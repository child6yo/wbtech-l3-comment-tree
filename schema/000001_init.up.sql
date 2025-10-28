CREATE TABLE comments (
    id          serial       NOT NULL PRIMARY KEY,
    content     text         NOT NULL,
    answer_at   integer      REFERENCES comments(id) ON DELETE CASCADE
);