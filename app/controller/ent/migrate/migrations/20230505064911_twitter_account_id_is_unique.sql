-- +goose Up
-- create index "twitter_accounts_twitter_account_id_key" to table: "twitter_accounts"
CREATE UNIQUE INDEX "twitter_accounts_twitter_account_id_key" ON "twitter_accounts" ("twitter_account_id");

-- +goose Down
-- reverse: create index "twitter_accounts_twitter_account_id_key" to table: "twitter_accounts"
DROP INDEX "twitter_accounts_twitter_account_id_key";
