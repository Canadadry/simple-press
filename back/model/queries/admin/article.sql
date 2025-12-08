-- name: CountArticle :one
SELECT
    count(*)
FROM
    article;

-- name: CountArticleBySlug :one
SELECT
    count(*)
FROM
    article
WHERE
    slug = ?;

-- name: CreateArticle :execlastid
INSERT INTO
    article (
        title,
        date,
        author,
        content,
        slug,
        draft,
        layout_id
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateArticle :exec
UPDATE article
SET
    title = ?,
    date = ?,
    author = ?,
    content = ?,
    slug = ?,
    draft = ?,
    layout_id = ?
WHERE
    slug = ?;

-- name: DeleteArticle :exec
DELETE FROM article
WHERE
    slug = ?;

-- name: SelectArticleBySlug :many
SELECT
    *
FROM
    article
WHERE
    slug = ?
LIMIT
    1;

-- name: GetArticleList :many
SELECT
    title,
    date,
    author,
    content,
    slug,
    draft
FROM
    article
ORDER BY
    date DESC
LIMIT
    ?
OFFSET
    ?;
