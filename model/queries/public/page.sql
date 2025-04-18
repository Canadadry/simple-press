-- name: SelectPageByID :many
SELECT
    *
FROM
    pages
WHERE
    id = ?
LIMIT
    1;
