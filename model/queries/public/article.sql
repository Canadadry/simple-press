-- name: SelectArticleBySlug :many
SELECT
    *
FROM
    articles
WHERE
    slug = ?
    and draft = false
LIMIT
    1;

-- name: GetArticlesList :many
SELECT
    title,
    date,
    author,
    slug,
    draft
FROM
    articles
WHERE
    draft = false
ORDER BY
    date DESC
LIMIT
    ?
OFFSET
    ?;
