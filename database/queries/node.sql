-- name: GetRandomNode :many
SELECT
  *
FROM
  nodes
ORDER BY
  RANDOM()
LIMIT
  $1;

-- name: CreateNode :exec
INSERT INTO
  nodes (id, address)
VALUES
  ($1, $2) RETURNING *;

-- name: GetAll :many
SELECT
  *
FROM
  nodes;