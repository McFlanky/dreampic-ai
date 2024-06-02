CREATE TABLE
    IF NOT EXISTS images (
        id SERIAL PRIMARY KEY,
        user_id UUID REFERENCES auth.users,
        status INT NOT NULL DEFAULT 1,
        prompt TEXT NOT NULL,
        deleted BOOLEAN NOT NULL DEFAULT 'false',
        image_location TEXT,
        deleted_at TIMESTAMP,
        created_at TIMESTAMP NOT NULL DEFAULT NOW ()
    );