-- name: GetAllBlock :many
SELECT
    *
FROM
    blocks
ORDER BY
    height DESC;


-- name: CreateBlock :exec
INSERT INTO
    blocks (id, data)
VALUES
    ($1, $2) RETURNING *;