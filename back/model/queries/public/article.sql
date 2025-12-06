-- name: SelectArticleBySlug :many
SELECT
    *
FROM
    article
WHERE
    slug = ?
    and draft = false
LIMIT
    1;

-- name: GetArticleList :many
SELECT
    title,
    date,
    author,
    slug,
    draft
FROM
    article
WHERE
    draft = false
ORDER BY
    date DESC
LIMIT
    ?
OFFSET
    ?;
