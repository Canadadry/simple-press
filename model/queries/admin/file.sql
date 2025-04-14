-- name: UploadFile :execlastid
INSERT INTO
    files (name, content)
VALUES
    (?, ?);

-- name: DownloadFile :many
SELECT
    content
FROM
    files
WHERE
    name = ?
LIMIT
    1;

-- name: DeleteFile :exec
DELETE FROM files
WHERE
    name = ?;

-- name: GetFileList :many
SELECT
    name
FROM
    files
ORDER BY
    id DESC
LIMIT
    ?
OFFSET
    ?;
