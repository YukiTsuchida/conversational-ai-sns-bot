-- +goose Up
-- modify "conversations" table
ALTER TABLE "conversations" ADD COLUMN "abort_reason" character varying NULL;

-- +goose Down
-- reverse: modify "conversations" table
ALTER TABLE "conversations" DROP COLUMN "abort_reason";
