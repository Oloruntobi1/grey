-- name: CreateTransaction :one
INSERT INTO transactions(
from_user_id,
to_user_id,
amount
) VALUES (
    $1, $2, $3
) RETURNING *;