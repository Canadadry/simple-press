-- name: CountPage :one
SELECT
    count(*)
FROM
    pages;

-- name: CountPageByName :one
SELECT
    count(*)
FROM
    pages
WHERE
    name = ?;

-- name: CountPageByID :one
SELECT
    count(*)
FROM
    pages
WHERE
    id = ?;

-- name: CreatePage :execlastid
INSERT INTO
    pages (name, content)
VALUES
    (?, ?);

-- name: UpdatePage :exec
UPDATE pages
SET
    name = ?,
    content = ?
WHERE
    name = ?;

-- name: DeletePage :exec
DELETE FROM pages
WHERE
    name = ?;

-- name: SelectPage :many
SELECT
    *
FROM
    pages
WHERE
    name = ?
LIMIT
    1;

-- name: SelectPageByID :many
SELECT
    *
FROM
    pages
WHERE
    id = ?
LIMIT
    1;

-- name: GetPageList :many
SELECT
    *
FROM
    pages
ORDER BY
    name DESC
LIMIT
    ?
OFFSET
    ?;

-- name: GetAllPages :many
SELECT
    *
FROM
    pages
ORDER BY
    id DESC;
