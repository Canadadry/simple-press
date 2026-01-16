-- name: CountFile :one
SELECT
    count(*)
FROM
    file;

-- name: CountFileByName :one
SELECT
    count(*)
FROM
    file
WHERE
    name = ?;

-- name: UploadFile :execlastid
INSERT INTO
    file (name, content)
VALUES
    (?, ?);

-- name: DownloadFile :many
SELECT
    content
FROM
    file
WHERE
    name = ?
LIMIT
    1;

-- name: DeleteFile :exec
DELETE FROM file
WHERE
    name = ?;

-- name: GetFileList :many
SELECT
    name
FROM
    file
ORDER BY
    id DESC
LIMIT
    ?
OFFSET
    ?;

-- name: SelectFoldersInFolderFile :many
SELECT DISTINCT
  substr(
    substr(name, length(:path) + 1),
    1,
    instr(substr(name, length(:path) + 1), '/') - 1
  ) AS folder
FROM file
WHERE name LIKE :path || '%'
  AND instr(substr(name, length(:path) + 1), '/') > 0;

-- name: SelectFilesInFolderFile :many
SELECT
    substr(name, length(:path) + 1) AS filename
FROM file
WHERE name LIKE :path || '%'
AND instr(substr(name, length(:path) + 1), '/') = 0;
