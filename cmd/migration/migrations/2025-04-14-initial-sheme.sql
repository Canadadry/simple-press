CREATE TABLE articles (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    date DATE NOT NULL,
    author TEXT NOT NULL,
    content TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    draft INTEGER NOT NULL
);

CREATE TABLE layouts (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    content TEXT NOT NULL
);

CREATE TABLE files (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    content BLOB NOT NULL
);
