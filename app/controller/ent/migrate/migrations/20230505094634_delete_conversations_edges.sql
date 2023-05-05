-- +goose Up
-- modify "conversations" table
ALTER TABLE "conversations" DROP CONSTRAINT "conversations_twitter_accounts_conversation";
-- modify "twitter_accounts" table
ALTER TABLE "twitter_accounts" ADD COLUMN "twitter_accounts_conversation" bigint NULL, ADD CONSTRAINT "twitter_accounts_conversations_conversation" FOREIGN KEY ("twitter_accounts_conversation") REFERENCES "conversations" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;

-- +goose Down
-- reverse: modify "twitter_accounts" table
ALTER TABLE "twitter_accounts" DROP CONSTRAINT "twitter_accounts_conversations_conversation", DROP COLUMN "twitter_accounts_conversation";
-- reverse: modify "conversations" table
ALTER TABLE "conversations" ADD CONSTRAINT "conversations_twitter_accounts_conversation" FOREIGN KEY ("twitter_accounts_conversation") REFERENCES "twitter_accounts" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
