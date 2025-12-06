-- name: DownloadFile :many
SELECT
    content
FROM
    file
WHERE
    name = ?
LIMIT
    1;
