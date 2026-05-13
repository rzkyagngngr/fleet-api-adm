package chat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type TelegramClient interface {
	CreateForumTopic(chatID int64, topicName string) (int64, error)
	SendMessage(chatID int64, threadID int64, text string) (int64, error)
	CloseForumTopic(chatID int64, threadID int64) error
	ReopenForumTopic(chatID int64, threadID int64) error
	EditForumTopic(chatID int64, threadID int64, topicName string) error
	DeleteMessage(chatID int64, messageID int64) error
	GetFileURL(fileID string) (string, error)
	DownloadFile(fileURL string) ([]byte, error)
}

type telegramClient struct {
	token   string
	baseURL string
	client  *http.Client
}

var errMissingTelegramToken = errors.New("telegram bot token is not configured")

func NewTelegramClientFromEnv() TelegramClient {
	return NewTelegramClient(os.Getenv("TELEGRAM_BOT_TOKEN"))
}

func NewTelegramClient(token string) TelegramClient {
	return &telegramClient{
		token:   strings.TrimSpace(token),
		baseURL: "https://api.telegram.org",
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *telegramClient) CreateForumTopic(chatID int64, topicName string) (int64, error) {
	var result struct {
		MessageThreadID int64 `json:"message_thread_id"`
	}

	if err := c.call("createForumTopic", map[string]interface{}{
		"chat_id":    chatID,
		"name":       topicName,
		"icon_color": 7322096,
	}, &result); err != nil {
		return 0, err
	}
	if result.MessageThreadID == 0 {
		return 0, errors.New("telegram createForumTopic did not return message_thread_id")
	}

	return result.MessageThreadID, nil
}

func (c *telegramClient) SendMessage(chatID int64, threadID int64, text string) (int64, error) {
	var result struct {
		MessageID int64 `json:"message_id"`
	}
	if err := c.call("sendMessage", map[string]interface{}{
		"chat_id":           chatID,
		"message_thread_id": threadID,
		"text":              text,
	}, &result); err != nil {
		return 0, err
	}
	return result.MessageID, nil
}

func (c *telegramClient) CloseForumTopic(chatID int64, threadID int64) error {
	return c.call("closeForumTopic", map[string]interface{}{
		"chat_id":           chatID,
		"message_thread_id": threadID,
	}, nil)
}

func (c *telegramClient) ReopenForumTopic(chatID int64, threadID int64) error {
	return c.call("reopenForumTopic", map[string]interface{}{
		"chat_id":           chatID,
		"message_thread_id": threadID,
	}, nil)
}

func (c *telegramClient) EditForumTopic(chatID int64, threadID int64, topicName string) error {
	return c.call("editForumTopic", map[string]interface{}{
		"chat_id":           chatID,
		"message_thread_id": threadID,
		"name":              topicName,
	}, nil)
}

func (c *telegramClient) DeleteMessage(chatID int64, messageID int64) error {
	return c.call("deleteMessage", map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
	}, nil)
}

func (c *telegramClient) GetFileURL(fileID string) (string, error) {
	var result struct {
		FilePath string `json:"file_path"`
	}
	if err := c.call("getFile", map[string]interface{}{
		"file_id": fileID,
	}, &result); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/file/bot%s/%s", c.baseURL, c.token, result.FilePath), nil
}

func (c *telegramClient) DownloadFile(fileURL string) ([]byte, error) {
	resp, err := c.client.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file, status: %d", resp.StatusCode)
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file body: %w", err)
	}

	return buf.Bytes(), nil
}

type telegramAPIResponse struct {
	OK          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	Description string          `json:"description"`
}

func (c *telegramClient) call(method string, payload interface{}, output interface{}) error {
	if c.token == "" {
		return errMissingTelegramToken
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal telegram payload: %w", err)
	}

	endpoint := fmt.Sprintf("%s/bot%s/%s", c.baseURL, c.token, method)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build telegram request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("telegram %s request failed: %w", method, err)
	}
	defer resp.Body.Close()

	var telegramResp telegramAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&telegramResp); err != nil {
		return fmt.Errorf("decode telegram %s response: %w", method, err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("telegram %s returned %d: %s", method, resp.StatusCode, telegramResp.Description)
	}
	if !telegramResp.OK {
		return fmt.Errorf("telegram %s failed: %s", method, telegramResp.Description)
	}

	if output != nil && len(telegramResp.Result) > 0 {
		if err := json.Unmarshal(telegramResp.Result, output); err != nil {
			return fmt.Errorf("decode telegram %s result: %w", method, err)
		}
	}

	return nil
}
