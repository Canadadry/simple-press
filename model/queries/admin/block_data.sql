-- name: CreateBlockData :execlastid
INSERT INTO
    block_data (position, data, block_id, article_id)
VALUES
    (?, ?, ?, ?);

-- name: UpdateBlockData :exec
UPDATE block_data
SET
    position = ?,
    data = ?
WHERE
    id = ?;

-- name: DeleteBlockData :exec
DELETE FROM block_data
WHERE
    id = ?;

-- name: SelectBlockDataByArticle :many
SELECT
    *
FROM
    block_data
WHERE
    article_id = ?
ORDER BY
    position DESC;
