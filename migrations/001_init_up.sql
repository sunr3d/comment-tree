CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    parent_id INTEGER REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    author VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_comments_created_at ON comments(created_at);
CREATE INDEX idx_comments_deleted_at ON comments(deleted_at);

-- Для полнотекстового поиска
CREATE INDEX idx_comments_content_gin ON comments USING gin(to_tsvector('russian', content));

GRANT ALL PRIVILEGES ON TABLE comments TO comment_tree_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO comment_tree_user;