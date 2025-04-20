-- name: CountBlock :one
SELECT
    count(*)
FROM
    block;

-- name: CountBlockByName :one
SELECT
    count(*)
FROM
    block
WHERE
    name = ?;

-- name: CreateBlock :execlastid
INSERT INTO
    block (name, content, definition)
VALUES
    (?, ?, ?);

-- name: UpdateBlock :exec
UPDATE block
SET
    name = ?,
    content = ?,
    definition = ?
WHERE
    name = ?;

-- name: DeleteBlock :exec
DELETE FROM block
WHERE
    name = ?;

-- name: SelectBlock :many
SELECT
    *
FROM
    block
WHERE
    name = ?
LIMIT
    1;

-- name: SelectAllBlock :many
SELECT
    *
FROM
    block;

-- name: GetBlockList :many
SELECT
    name
FROM
    block
ORDER BY
    name DESC
LIMIT
    ?
OFFSET
    ?;
