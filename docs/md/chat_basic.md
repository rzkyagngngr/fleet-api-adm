# POS BE Chat - Technical Handover

## 1. Project Overview

`pos-be-chat` adalah backend service berbasis Go untuk orkestrasi komunikasi operasional vessel call melalui Telegram forum topic.

Scope inti:

- Menyimpan master vessel dan mapping ke Telegram supergroup (`telegram_chat_id`)
- Membuat vessel call (voyage) dan auto-create Telegram forum topic (`message_thread_id`)
- Mengelola participant per call (authorization level DB)
- Mengirim pesan bot ke topic vessel call
- Menerima webhook Telegram, melakukan validasi authorization, menyimpan message mirror ke PostgreSQL

Tech stack:

- Go `1.25.6`
- HTTP framework: `gin-gonic/gin`
- ORM: `gorm`
- DB: `PostgreSQL`
- API docs: `swaggo` (Swagger UI)


## 2. High-Level Architecture

Pola layer:

1. `handler`  
   HTTP boundary: parse request, validasi binding, format response.
2. `service`  
   Business orchestration: alur vessel call, topic operations, webhook rules.
3. `repository`  
   Persistence logic via GORM.
4. `telegram`  
   Telegram Bot API adapter (HTTP request ke `api.telegram.org`).

Entry point: `cmd/server/main.go`

- Load `.env`
- Init DB + auto migration
- Init repository/service/handler
- Register routes API, webhook, static UI, swagger


## 3. Data Model (Transactional Domain)

Model ada di `internal/model/models.go`.

### 3.1 `Vessel`

- `ID` (PK)
- `Name`
- `TelegramChatID` (supergroup ID)

### 3.2 `VesselCall`

- `ID` (PK)
- `VesselID` (FK logis ke `Vessel`)
- `VoyageCode`
- `ETA`
- `TelegramThreadID` (forum topic ID)

### 3.3 `CallParticipant`

- `ID` (PK)
- `CallID` (FK logis ke `VesselCall`)
- `TelegramUserID`
- `Role` (`agent` / `pbm` / `operator`)

### 3.4 `ChatMessage`

- `ID` (PK)
- `CallID`
- `TelegramMessageID`
- `SenderID`
- `Text`
- `Timestamp`
- `RawPayload` (`jsonb`)


## 4. Core Transactional Flows

## 4.1 Create Vessel Call -> Create Telegram Topic -> Persist Call

Flow ini adalah transaksi bisnis paling penting karena jadi penghubung DB dan thread Telegram.

```go
func (s *Service) CreateVesselCall(vesselID uint, voyageCode string, eta time.Time) (*model.VesselCall, error) {
	vessel, err := s.Repo.GetVesselByID(vesselID)
	if err != nil {
		return nil, err
	}

	// Create forum topic automatically
	topicName := fmt.Sprintf("ETA %s - %s", eta.Format("02Jan"), voyageCode)
	threadID, err := s.Bot.CreateForumTopic(vessel.TelegramChatID, topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to create forum topic: %w", err)
	}

	call := &model.VesselCall{
		VesselID:         vesselID,
		VoyageCode:       voyageCode,
		ETA:              eta,
		TelegramThreadID: threadID,
	}
	if err := s.Repo.CreateVesselCall(call); err != nil {
		return nil, err
	}
	return call, nil
}
```

Lokasi: `internal/service/service.go`

Catatan:

- Jika create topic sukses tapi insert DB gagal, saat ini belum ada compensating action (topic tidak dihapus).
- Perlu dipertimbangkan pattern idempotency atau saga jika production critical.


## 4.2 Send Message to Specific Vessel Topic

Kirim message dilakukan dengan lookup `VesselCall` -> `Vessel` untuk ambil `chat_id` + `thread_id`.

```go
func (s *Service) SendMessage(callID uint, text string) error {
	var call model.VesselCall
	if err := s.Repo.DB.First(&call, callID).Error; err != nil {
		return err
	}

	var vessel model.Vessel
	if err := s.Repo.DB.First(&vessel, call.VesselID).Error; err != nil {
		return err
	}

	return s.Bot.SendMessage(vessel.TelegramChatID, call.TelegramThreadID, text)
}
```

Lokasi: `internal/service/service.go`

Nilai bisnis:

- Menjamin pesan masuk ke thread voyage yang tepat.


## 4.3 Webhook Enforcement + Message Mirroring

Ini flow paling penting untuk governance komunikasi.

```go
func (s *Service) HandleWebhookUpdate(update *telegram.Update) {
	if update.Message == nil {
		log.Println("Webhook received but no message found inside update")
		return
	}
	msg := update.Message

	vc, err := s.Repo.GetVesselCallByThreadID(msg.Chat.ID, msg.MessageThreadID)
	if err != nil {
		return
	}

	if !s.Repo.IsParticipantAuthorized(vc.ID, msg.From.ID) {
		s.Bot.DeleteMessage(msg.Chat.ID, msg.MessageID)
	}

	raw, _ := json.Marshal(msg)
	chatMsg := &model.ChatMessage{
		CallID:            vc.ID,
		TelegramMessageID: msg.MessageID,
		SenderID:          msg.From.ID,
		Text:              msg.Text,
		Timestamp:         time.Unix(int64(msg.Date), 0),
		RawPayload:        string(raw),
	}
	s.Repo.SaveChatMessage(chatMsg)
}
```

Lokasi: `internal/service/service.go`

Nilai bisnis:

- Enforce participant whitelist per vessel call.
- Menyimpan jejak percakapan untuk audit/troubleshooting.


## 5. API Surface

Base path: `/api`

Endpoints utama:

- `POST /vessels`
- `GET /vessels`
- `POST /vessel-calls`
- `GET /vessel-calls`
- `POST /vessel-calls/:id/suspend`
- `POST /vessel-calls/:id/continue`
- `PUT /vessel-calls/:id/rename`
- `GET /vessel-calls/:id/participants`
- `DELETE /vessel-calls/:id/participants/:user_id`
- `POST /participants`
- `POST /invite` (alias requirement)
- `POST /send-message`
- `POST /webhook/telegram` (non `/api`)

Swagger UI:

- `GET /swagger/index.html`


## 6. Environment & Configuration

File: `.env`

Variabel runtime:

- `PORT`
- `DATABASE_URL`
- `TELEGRAM_BOT_TOKEN`

Security notes:

- `TELEGRAM_BOT_TOKEN` terdeteksi tersimpan plain-text di repo lokal.
- Wajib rotate token sebelum handover/deploy.
- Gunakan secret manager pada environment production.


## 7. Local Runbook

## 7.1 Start DB

```bash
docker compose up -d
```

## 7.2 Run App

```bash
go run ./cmd/server
```

## 7.3 Verify

- Health basic: `GET /api/vessels` harus return `200`
- Swagger: buka `/swagger/index.html`
- UI tester: buka `/`


## 8. Telegram Integration Notes

Implementasi client ada di `internal/telegram/client.go`.

Telegram methods yang dipakai:

- `createForumTopic`
- `sendMessage`
- `sendPhoto` / `sendDocument`
- `closeForumTopic`
- `reopenForumTopic`
- `editForumTopic`
- `deleteMessage`

Batasan yang sudah diakui di kode:

- Bot API tidak bisa create supergroup langsung.
- Fetch historical topic messages tidak tersedia via Bot API standar.


## 9. Known Gaps / Technical Debt

1. Invitation flow masih mock

- Mapping `phone_number -> telegram_user_id` belum real.
- Invite link return static dummy link.

2. Transaction consistency

- Belum ada rollback kompensasi saat call ke Telegram berhasil tapi DB gagal (dan sebaliknya).

3. Error handling webhook

- `SaveChatMessage` dipanggil tanpa cek error.
- Unauthorized message tetap disimpan (sesuai komentar kode saat ini).

4. Validation & auth API

- Belum ada auth middleware untuk endpoint API internal.
- Belum ada input-level validation lanjutan (format phone, role whitelist strict, dsb).

5. Observability

- Logging masih basic, belum structured logging + trace correlation.


## 10. Handover Checklist (Recommended)

1. Rotate `TELEGRAM_BOT_TOKEN` dan pindahkan ke secret manager.
2. Tambahkan migration strategy explicit (bukan hanya `AutoMigrate`) untuk production.
3. Finalisasi mapping participant dari sumber data user yang valid.
4. Tambahkan authentication/authorization layer untuk API.
5. Tambahkan test coverage minimum:
- service flow `CreateVesselCall`
- webhook authorization path
- repository query `GetVesselCallByThreadID`
6. Tambahkan retry/backoff policy untuk Telegram API calls.
7. Definisikan SOP ketika topic creation berhasil tapi persist gagal.


## 11. File Map Cepat

- `cmd/server/main.go` -> bootstrap app + route registry
- `internal/handler/handler.go` -> HTTP endpoints
- `internal/service/service.go` -> business logic transaksional
- `internal/repository/repository.go` -> query DB
- `internal/repository/db.go` -> DB init + automigrate
- `internal/telegram/client.go` -> adapter Telegram Bot API
- `internal/model/models.go` -> schema model domain
- `public/index.html` -> tester UI manual
- `docs/swagger.yaml` -> OpenAPI spec

