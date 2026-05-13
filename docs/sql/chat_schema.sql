-- Chat schema bootstrap for omniport main database.
-- Source of vessel-call data is plan.post_vessel_schedules (no chat.vessels / chat.vessel_calls dependency).

CREATE SCHEMA IF NOT EXISTS chat;

ALTER TABLE plan.post_vessel_schedules
    ADD COLUMN IF NOT EXISTS telegram_topic_id BIGINT,
    ADD COLUMN IF NOT EXISTS telegram_topic_name VARCHAR(200);

ALTER TABLE plan.post_vessel_schedules
    DROP CONSTRAINT IF EXISTS chk_post_vessel_schedules_telegram_topic_status,
    DROP COLUMN IF EXISTS telegram_topic_status;

CREATE INDEX IF NOT EXISTS idx_post_vessel_schedules_telegram_topic_id
    ON plan.post_vessel_schedules (telegram_topic_id);

CREATE TABLE IF NOT EXISTS chat.schedule_participants (
    id BIGSERIAL PRIMARY KEY,
    schedule_id BIGINT NOT NULL REFERENCES plan.post_vessel_schedules(id) ON DELETE CASCADE,
    internal_user_id BIGINT,
    telegram_user_id BIGINT NOT NULL,
    role VARCHAR(20) NOT NULL,
    status SMALLINT NOT NULL DEFAULT 1,
    creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    creation_by VARCHAR(100),
    last_updated_date TIMESTAMP WITHOUT TIME ZONE,
    last_updated_by VARCHAR(100),
    CONSTRAINT chk_chat_schedule_participant_role CHECK (role IN ('agent', 'pbm', 'operator'))
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_chat_schedule_participants_schedule_telegram_user
    ON chat.schedule_participants (schedule_id, telegram_user_id);

CREATE TABLE IF NOT EXISTS chat.schedule_messages (
    id BIGSERIAL PRIMARY KEY,
    schedule_id BIGINT NOT NULL REFERENCES plan.post_vessel_schedules(id) ON DELETE CASCADE,
    telegram_message_id BIGINT NOT NULL,
    telegram_chat_id BIGINT NOT NULL,
    telegram_topic_id BIGINT NOT NULL,
    sender_id BIGINT NOT NULL,
    text TEXT,
    attatchment JSONB,
    message_timestamp TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    raw_payload JSONB,
    is_authorized BOOLEAN NOT NULL DEFAULT TRUE,
    creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE chat.schedule_messages
    ADD COLUMN IF NOT EXISTS attatchment JSONB,
    ADD COLUMN IF NOT EXISTS telegram_chat_id BIGINT,
    ADD COLUMN IF NOT EXISTS telegram_topic_id BIGINT;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'chat'
          AND table_name = 'schedule_messages'
          AND column_name = 'telegram_channel_id'
    ) THEN
        UPDATE chat.schedule_messages
        SET telegram_topic_id = telegram_channel_id
        WHERE telegram_topic_id IS NULL
          AND telegram_channel_id IS NOT NULL;
    END IF;
END $$;

ALTER TABLE chat.schedule_messages
    DROP COLUMN IF EXISTS telegram_channel_id;

CREATE UNIQUE INDEX IF NOT EXISTS uq_chat_schedule_messages_schedule_msg
    ON chat.schedule_messages (schedule_id, telegram_message_id);

CREATE INDEX IF NOT EXISTS idx_chat_schedule_messages_chat_topic
    ON chat.schedule_messages (telegram_chat_id, telegram_topic_id);

GRANT USAGE ON SCHEMA chat TO omniport;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA chat TO omniport;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA chat TO omniport;
GRANT SELECT, UPDATE ON TABLE plan.post_vessel_schedules TO omniport;
