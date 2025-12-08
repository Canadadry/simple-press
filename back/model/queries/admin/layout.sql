-- name: CountLayout :one
SELECT
    count(*)
FROM
    layout;

-- name: CountLayoutByName :one
SELECT
    count(*)
FROM
    layout
WHERE
    name = ?;

-- name: CountLayoutByID :one
SELECT
    count(*)
FROM
    layout
WHERE
    id = ?;

-- name: CreateLayout :execlastid
INSERT INTO
    layout (name, content)
VALUES
    (?, ?);

-- name: UpdateLayout :exec
UPDATE layout
SET
    name = ?,
    content = ?
WHERE
    name = ?;

-- name: DeleteLayout :exec
DELETE FROM layout
WHERE
    name = ?;

-- name: SelectLayout :many
SELECT
    *
FROM
    layout
WHERE
    name = ?
LIMIT
    1;

-- name: SelectLayoutByID :many
SELECT
    *
FROM
    layout
WHERE
    id = ?
LIMIT
    1;

-- name: GetLayoutList :many
SELECT
    id,
    name,
    substr(content,0,50) as `content`
FROM
    layout
ORDER BY
    name DESC
LIMIT
    ?
OFFSET
    ?;


-- name: GetAllLayout :many
SELECT
    id,
    name
FROM
    layout
ORDER BY
    id DESC
;
