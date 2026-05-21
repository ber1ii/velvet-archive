-- name: ListSeries :many
SELECT * FROM series 
ORDER BY title ASC;

-- name: GetSeries :one
SELECT * FROM series 
WHERE id = $1;

-- name: GetLoreEntriesBySeries :many
SELECT * FROM lore_entries 
WHERE series_id = $1 
ORDER BY title ASC;

-- name: GetLoreEntry :one
SELECT * FROM lore_entries 
WHERE id = $1;

-- name: GetRevealedLinksForEntry :many
SELECT 
    ll.id as link_id,
    ll.label,
    ll.is_revealed,
    le.id as target_entry_id,
    le.title as target_title,
    le.category as target_category
FROM lore_links ll
JOIN lore_entries le ON ll.target_id = le.id
WHERE ll.source_id = $1 AND ll.is_revealed = TRUE;

-- name: SearchLoreEntries :many
SELECT id, series_id, title, category, content, created_at
FROM lore_entries
WHERE search_vec @@ plainto_tsquery('english', $1)
ORDER BY ts_rank(search_vec, plainto_tsquery('english', $1)) DESC;

-- name: CreateSeries :one
INSERT INTO series (title, author, cover_color, description)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateLoreEntry :one
INSERT INTO lore_entries (series_id, title, category, content, metadata, affinity)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CreateLoreLink :one
INSERT INTO lore_links (source_id, target_id, label, is_revealed)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: RevealLoreLink :one
UPDATE lore_links
SET is_revealed = TRUE
WHERE id = $1
RETURNING *;

-- name: UpdateLoreEntry :one
UPDATE lore_entries
SET title = $2, category = $3, content = $4, metadata = $5, affinity = $6
WHERE id = $1
RETURNING *;

-- name: DeleteLoreEntry :exec
DELETE FROM lore_entries 
WHERE id = $1;

-- name: GetAdminByEmail :one
SELECT * FROM admin_users 
WHERE email = $1;

-- name: CreateAdminUser :one
INSERT INTO admin_users (email, password_hash)
VALUES ($1, $2)
RETURNING *;