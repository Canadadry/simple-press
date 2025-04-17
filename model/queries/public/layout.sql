-- name: SelectLayout :many
SELECT
    *
FROM
    layouts
WHERE
    name = ?
LIMIT
    1;

-- name: SelectBaseLayout :many
SELECT
    *
FROM
    layouts
WHERE
    name like "_layout/%";
