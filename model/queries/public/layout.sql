-- name: SelectLayout :many
SELECT
    *
FROM
    layouts
WHERE
    name = ?
LIMIT
    1;
