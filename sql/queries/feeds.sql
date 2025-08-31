-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: ListFeeds :many
SELECT * FROM feeds;

-- name: GetFeedAndCreator :many
SELECT feeds.name, feeds.url, users.name AS user_name
FROM feeds
INNER JOIN users
ON feeds.user_id = users.id
ORDER BY feeds.created_at ASC;

-- name: GetFeedFromURL :one
SELECT * FROM feeds WHERE feeds.url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1
RETURNING *;
-- UPDATE feeds 
-- SET updated_at = $1, last_fetched_at = $1 
-- WHERE id = $2
-- RETURNING *;


-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;

-- SELECT * FROM feeds
-- WHERE last_fetched_at < $1 OR last_fetched_at IS NULL
-- ORDER BY last_fetched_at ASC NULLS FIRST, id ASC
-- LIMIT 1;
