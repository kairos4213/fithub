-- +goose Up
ALTER TABLE exercises ADD COLUMN video_url TEXT;

-- +goose Down
ALTER TABLE exercises DROP COLUMN video_url;
