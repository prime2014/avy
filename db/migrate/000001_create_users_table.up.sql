CREATE TABLE users(
    id BIGSERIAL PRIMARY KEY,
    firstname VARCHAR(30) NOT NULL,
    lastname VARCHAR(30) NOT NULL,
    email VARCHAR(30) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT false,
    refresh_token VARCHAR(255),
    date_joined TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);