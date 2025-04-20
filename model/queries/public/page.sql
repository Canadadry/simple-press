-- name: SelectPageByID :many
SELECT
    *
FROM
    page
WHERE
    id = ?
LIMIT
    1;
