CREATE TABLE words (
    id SERIAL PRIMARY KEY,
    word TEXT UNIQUE NOT NULL,
    data JSONB NOT NULL
);
