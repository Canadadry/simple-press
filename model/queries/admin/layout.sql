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

-- name: GeLayoutList :many
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
