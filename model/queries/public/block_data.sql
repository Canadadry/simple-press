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
