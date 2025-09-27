-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (gen_random_uuid(), now(), now(), $1)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (gen_random_uuid(), now(), now(), $1, $2)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at;

-- name: GetChirpByID :one
SELECT * FROM chirps
WHERE id = $1;