-- name: DownloadFile :many
SELECT
    content
FROM
    files
WHERE
    name = ?
LIMIT
    1;
