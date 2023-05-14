-- +goose Up
-- modify "chatgpt35turbo_conversation_logs" table
ALTER TABLE "chatgpt35turbo_conversation_logs" ALTER COLUMN "purpose" DROP NOT NULL;

-- +goose Down
-- reverse: modify "chatgpt35turbo_conversation_logs" table
ALTER TABLE "chatgpt35turbo_conversation_logs" ALTER COLUMN "purpose" SET NOT NULL;
