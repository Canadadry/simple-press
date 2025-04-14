-- name: UploadFile :execlastid
INSERT INTO
    files (name, content, uuid)
VALUES
    (?, ?, ?);

-- name: DownloadFile :one
SELECT
    content
FROM
    files
WHERE
    uuid = ?
LIMIT
    1;

-- name: DeleteFile :exec
DELETE FROM files
WHERE
    uuid = ?;

-- name: GetFileList :many
SELECT
    uuid,
    name
FROM
    files
ORDER BY
    id DESC
LIMIT
    ?
OFFSET
    ?;
