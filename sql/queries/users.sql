-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password)
VALUES (gen_random_uuid(), now(), now(), $1, $2)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (gen_random_uuid(), now(), now(), $1, $2)
RETURNING *;

-- name: GetChirpsAsc :many
SELECT * FROM chirps
ORDER BY created_at;

-- name: GetChirpsDesc :many
SELECT * FROM chirps
ORDER BY created_at DESC;

-- name: GetChirpByID :one
SELECT * FROM chirps
WHERE id = $1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT u.* FROM users u
JOIN refresh_tokens rt ON u.id = rt.user_id
WHERE rt.token = $1 AND rt.expires_at > NOW() AND rt.revoked_at IS NULL;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens 
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

-- name: UpdateUser :one
UPDATE users
SET email = $1, password = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id=$1 AND user_id=$2;

-- name: UpdateUserIsChirpyRed :one
UPDATE users
SET is_chirpy_red = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: GetChirpsByAuthorID :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at;