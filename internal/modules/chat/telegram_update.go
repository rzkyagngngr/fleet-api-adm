package chat

type TelegramUpdate struct {
	UpdateID int64            `json:"update_id"`
	Message  *TelegramMessage `json:"message"`
}

type TelegramMessage struct {
	MessageID       int64              `json:"message_id"`
	MessageThreadID int64              `json:"message_thread_id"`
	Date            int64              `json:"date"`
	Text            string             `json:"text"`
	Caption         string             `json:"caption"`
	Photo           []TelegramPhoto    `json:"photo"`
	Document        *TelegramDocument  `json:"document"`
	Video           *TelegramVideo     `json:"video"`
	Audio           *TelegramAudio     `json:"audio"`
	Voice           *TelegramVoice     `json:"voice"`
	Animation       *TelegramAnimation `json:"animation"`
	Sticker         *TelegramSticker   `json:"sticker"`
	From            *TelegramUser      `json:"from"`
	Chat            *TelegramChat      `json:"chat"`
}

type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}


type TelegramChat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}


type TelegramPhoto struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size"`
}

type TelegramDocument struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileName     string `json:"file_name"`
	MimeType     string `json:"mime_type"`
	FileSize     int    `json:"file_size"`
}

type TelegramVideo struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type"`
	FileSize     int    `json:"file_size"`
}

type TelegramAudio struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	Performer    string `json:"performer"`
	Title        string `json:"title"`
	FileName     string `json:"file_name"`
	MimeType     string `json:"mime_type"`
	FileSize     int    `json:"file_size"`
}

type TelegramVoice struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type"`
	FileSize     int    `json:"file_size"`
}

type TelegramAnimation struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Duration     int    `json:"duration"`
	FileName     string `json:"file_name"`
	MimeType     string `json:"mime_type"`
	FileSize     int    `json:"file_size"`
}

type TelegramSticker struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Type         string `json:"type"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	IsAnimated   bool   `json:"is_animated"`
	IsVideo      bool   `json:"is_video"`
}
