-- =============================
-- USERS
-- =============================
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,

    first_name VARCHAR(100),
    last_name VARCHAR(100),
    gender VARCHAR(20),
    birth_date VARCHAR(20),
    bio TEXT,
    profile_picture VARCHAR(255),

    rating INT NOT NULL DEFAULT 0,
    level INT NOT NULL DEFAULT 1,
    league VARCHAR(50) NOT NULL DEFAULT 'Green Seed',

    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX ON users (rating DESC);


-- =============================
-- EMAIL VERIFICATION
-- =============================
CREATE TABLE IF NOT EXISTS email_verifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL
);


-- =============================
-- PASSWORD RESETS
-- =============================
CREATE TABLE IF NOT EXISTS password_resets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN NOT NULL DEFAULT FALSE
);


-- =============================
-- NEWS
-- =============================
CREATE TABLE IF NOT EXISTS news (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    link TEXT UNIQUE NOT NULL,
    published_at TIMESTAMP NOT NULL,
    source VARCHAR(255),
    description TEXT
);


-- =============================
-- ECO QUESTIONS (справочник)
-- =============================
CREATE TABLE IF NOT EXISTS eco_questions (
    id BIGSERIAL PRIMARY KEY,
    category VARCHAR(50) NOT NULL,
    question TEXT NOT NULL,
    max_value INT NOT NULL DEFAULT 5
);


-- =============================
-- ECO ANSWERS (ответы)
-- =============================
CREATE TABLE IF NOT EXISTS eco_answers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question_id BIGINT NOT NULL REFERENCES eco_questions(id) ON DELETE CASCADE,
    value INT NOT NULL CHECK (value >= 0 AND value <= 5),
    created_at TIMESTAMP DEFAULT NOW()
);


-- =============================
-- ECO RESULTS (история расчётов)
-- =============================
CREATE TABLE IF NOT EXISTS eco_results (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total_score INT NOT NULL,
    category VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


-- =============================
-- ECO ACTIONS (не обязательно, но полезно)
-- =============================
CREATE TABLE IF NOT EXISTS eco_actions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    points INT NOT NULL DEFAULT 5
);


-- =============================
-- USER ACTIONS (история действий)
-- =============================
CREATE TABLE IF NOT EXISTS user_actions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action_id BIGINT NOT NULL REFERENCES eco_actions(id),
    points INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


-- =============================
-- NOTIFICATIONS
-- =============================
CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    read BOOLEAN DEFAULT FALSE
);
