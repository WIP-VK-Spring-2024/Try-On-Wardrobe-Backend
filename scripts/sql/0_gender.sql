DO $$ BEGIN
    CREATE TYPE gender AS ENUM ('male', 'female', 'unisex', 'unknown');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
