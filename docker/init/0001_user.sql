CREATE TABLE auth (
    username VARCHAR(255) PRIMARY KEY,
    password TEXT NOT NULL,
    rights INT
);