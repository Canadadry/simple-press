-- name: CountTemplate :one
SELECT
    count(*)
FROM
    template;

-- name: CountTemplateByName :one
SELECT
    count(*)
FROM
    template
WHERE
    name = ?;

-- name: CreateTemplate :execlastid
INSERT INTO
    template (name, content)
VALUES
    (?, ?);

-- name: UpdateTemplate :exec
UPDATE template
SET
    name = ?,
    content = ?
WHERE
    name = ?;

-- name: DeleteTemplate :exec
DELETE FROM template
WHERE
    name = ?;

-- name: SelectTemplate :many
SELECT
    *
FROM
    template
WHERE
    name = ?
LIMIT
    1;

-- name: SelectAllTemplate :many
SELECT
    *
FROM
    template;

-- name: GetTemplateList :many
SELECT
    name,
    substr(content,0,50) as `content`
FROM
    template
ORDER BY
    name DESC
LIMIT
    ?
OFFSET
    ?;
