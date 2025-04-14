-- name: DownloadFile :one
SELECT
    content
FROM
    files
WHERE
    uuid = ?
LIMIT
    1;
