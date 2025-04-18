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

-- name: GePageList :many
SELECT
    name
FROM
    pages
ORDER BY
    name DESC
LIMIT
    ?
OFFSET
    ?;
