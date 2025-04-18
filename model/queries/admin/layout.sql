-- name: CountLayout :one
SELECT
    count(*)
FROM
    layouts;

-- name: CountLayoutByName :one
SELECT
    count(*)
FROM
    layouts
WHERE
    name = ?;

-- name: CreateLayout :execlastid
INSERT INTO
    layouts (name, content)
VALUES
    (?, ?);

-- name: UpdateLayout :exec
UPDATE layouts
SET
    name = ?,
    content = ?
WHERE
    name = ?;

-- name: DeleteLayout :exec
DELETE FROM layouts
WHERE
    name = ?;

-- name: SelectLayout :many
SELECT
    *
FROM
    layouts
WHERE
    name = ?
LIMIT
    1;

-- name: SelectAllLayout :many
SELECT
    *
FROM
    layouts;

-- name: GetLayoutList :many
SELECT
    name
FROM
    layouts
ORDER BY
    name DESC
LIMIT
    ?
OFFSET
    ?;
