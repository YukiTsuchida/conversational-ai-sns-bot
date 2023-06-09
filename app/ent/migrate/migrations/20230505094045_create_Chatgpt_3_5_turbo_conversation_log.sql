-- +goose Up
-- create "chatgpt35turbo_conversation_logs" table
CREATE TABLE "chatgpt35turbo_conversation_logs" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "message" character varying NOT NULL, "purpose" character varying NOT NULL, "role" character varying NOT NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "chatgpt35turbo_conversation_log_conversation" bigint NULL, PRIMARY KEY ("id"), CONSTRAINT "chatgpt35turbo_conversation_logs_conversations_conversation" FOREIGN KEY ("chatgpt35turbo_conversation_log_conversation") REFERENCES "conversations" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);

-- +goose Down
-- reverse: create "chatgpt35turbo_conversation_logs" table
DROP TABLE "chatgpt35turbo_conversation_logs";
