-- name: SelectBlockDataByArticle :many
SELECT
    *
FROM
    block_data
WHERE
    article_id = ?
ORDER BY
    position DESC;
