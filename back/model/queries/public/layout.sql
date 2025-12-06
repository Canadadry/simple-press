-- name: SelectLayoutByID :many
SELECT
    *
FROM
    layout
WHERE
    id = ?
LIMIT
    1;
