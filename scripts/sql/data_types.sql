DO $$ BEGIN
    CREATE TYPE gender AS ENUM ('male', 'female', 'unisex', 'unknown');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE season AS ENUM ('winter', 'spring', 'summer', 'autumn');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
