-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (
    token,
    created_at,
    updated_at,
    user_id,
    expires_at,
    revoked_at
) VALUES ($1, NOW(), NOW(), $2, $3, $4)
RETURNING *;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

-- name: GetUserFromRefreshToken :one
SELECT
    u.id,
    u.is_admin
FROM users AS u
INNER JOIN refresh_tokens AS rt ON u.id = rt.user_id
WHERE
    rt.token = $1
    AND rt.revoked_at IS NULL
    AND rt.expires_at > NOW();
