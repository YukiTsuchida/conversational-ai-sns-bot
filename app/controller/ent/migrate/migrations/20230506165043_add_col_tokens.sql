-- +goose Up
-- modify "conversations" table
ALTER TABLE "conversations" DROP COLUMN "twitter_accounts_conversation";
-- modify "twitter_accounts" table
ALTER TABLE "twitter_accounts" DROP COLUMN "bearer_token", ADD COLUMN "access_token" character varying NOT NULL, ADD COLUMN "refresh_token" character varying NOT NULL;

-- +goose Down
-- reverse: modify "twitter_accounts" table
ALTER TABLE "twitter_accounts" DROP COLUMN "refresh_token", DROP COLUMN "access_token", ADD COLUMN "bearer_token" character varying NOT NULL;
-- reverse: modify "conversations" table
ALTER TABLE "conversations" ADD COLUMN "twitter_accounts_conversation" bigint NULL;
