-- name: CountArticles :one
SELECT
    count(*)
FROM
    articles;

-- name: CreateArticle :execlastid
INSERT INTO
    articles (title, date, author, content, slug, draft)
VALUES
    (?, ?, ?, ?, ?, ?);

-- name: UpdateArticle :exec
UPDATE articles
SET
    title = ?,
    date = ?,
    author = ?,
    content = ?,
    slug = ?,
    draft = ?
WHERE
    slug = ?;

-- name: DeleteArticle :exec
DELETE FROM articles
WHERE
    slug = ?;

-- name: SelectArticleBySlug :many
SELECT
    *
FROM
    articles
WHERE
    slug = ?
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
ORDER BY
    date DESC
LIMIT
    ?
OFFSET
    ?;
