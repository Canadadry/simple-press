-- name: CountPage :one
SELECT
    count(*)
FROM
    page;

-- name: CountPageByName :one
SELECT
    count(*)
FROM
    page
WHERE
    name = ?;

-- name: CountPageByID :one
SELECT
    count(*)
FROM
    page
WHERE
    id = ?;

-- name: CreatePage :execlastid
INSERT INTO
    page (name, content)
VALUES
    (?, ?);

-- name: UpdatePage :exec
UPDATE page
SET
    name = ?,
    content = ?
WHERE
    name = ?;

-- name: DeletePage :exec
DELETE FROM page
WHERE
    name = ?;

-- name: SelectPage :many
SELECT
    *
FROM
    page
WHERE
    name = ?
LIMIT
    1;

-- name: SelectPageByID :many
SELECT
    *
FROM
    page
WHERE
    id = ?
LIMIT
    1;

-- name: GetPageList :many
SELECT
    *
FROM
    page
ORDER BY
    name DESC
LIMIT
    ?
OFFSET
    ?;

-- name: GetAllPage :many
SELECT
    *
FROM
    page
ORDER BY
    id DESC;
