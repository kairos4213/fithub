-- name: CreateAuthProvider :one
INSERT INTO auth_providers (
    user_id,
    provider,
    provider_user_id
) VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAuthProvider :one
SELECT * FROM auth_providers
WHERE provider = $1 AND provider_user_id = $2;

-- name: GetAuthProvidersByUserID :many
SELECT * FROM auth_providers
WHERE user_id = $1;

-- name: DeleteAuthProvider :exec
DELETE FROM auth_providers
WHERE id = $1;
