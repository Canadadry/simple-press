-- name: GetGlobalDefinition :one
SELECT
    content
FROM
    global
WHERE
    name = 'definition'
LIMIT
    1;


-- name: GetGlobalData :one
SELECT
    content
FROM
    global
WHERE
    name = 'data'
LIMIT
    1;


-- name: UpdateGlobalDefinition :exec
UPDATE
    global
SET
    content = ?
WHERE
    name = 'definition';


-- name: UpdateGlobalData :exec
UPDATE
    global
SET
    content = ?
WHERE
    name = 'data';
