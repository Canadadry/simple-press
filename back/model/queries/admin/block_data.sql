-- name: CountBlockDataByID :one
SELECT
    count(*)
FROM
    block_data
WHERE
    id = ?;

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

-- name: SelectBlockDataByID :many
SELECT
    *
FROM
    block_data
WHERE
    id = ?
LIMIT 1;

-- name: DeleteBlockData :exec
DELETE FROM block_data
WHERE
    id = ?;

-- name: SelectBlockDataByArticle :many
SELECT
    bd.id as `id`,
    bd.position as `position`,
    bd.data as `data`,
    bd.article_id as `article_id`,
    bd.block_id as `block_id`,
    b.name as `name`
FROM
    block_data as bd
    LEFT JOIN block as b on b.id = bd.block_id
WHERE
    bd.article_id = ?
ORDER BY
    bd.position DESC;
