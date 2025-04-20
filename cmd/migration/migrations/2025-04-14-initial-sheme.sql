CREATE TABLE article (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    date DATE NOT NULL,
    author TEXT NOT NULL,
    content TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    draft INTEGER NOT NULL,
    layout_id INTEGER NOT NULL,
    FOREIGN KEY (layout_id) REFERENCES layout (id)
);

CREATE TABLE template (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    content TEXT NOT NULL
);

CREATE TABLE layout (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    content TEXT NOT NULL
);

CREATE TABLE file (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    content BLOB NOT NULL
);
