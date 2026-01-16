-- name: CountArticleLikeTitle :one
SELECT
    count(*)
FROM
    article
WHERE
    title LIKE ?;


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
    substr(content, 0, 50) AS `content`,
    slug,
    draft
FROM
    article
WHERE
    title LIKE ?
ORDER BY
    date DESC
LIMIT
    ?
OFFSET
    ?;


-- name: SelectFoldersInFolderArticle :many
SELECT DISTINCT
    substr(
    substr(slug, length(:path) + 1),
    1,
    instr(substr(slug, length(:path) + 1), '/') - 1
    ) AS folder
FROM article
WHERE slug LIKE :path || '%'
    AND instr(substr(slug, length(:path) + 1), '/') > 0;

-- name: SelectArticlesInFolderArticle :many
SELECT
    substr(slug, length(:path) + 1) AS filename
FROM article
WHERE slug LIKE :path || '%'
AND instr(substr(slug, length(:path) + 1), '/') = 0;
