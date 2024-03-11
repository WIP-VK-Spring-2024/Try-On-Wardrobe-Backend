DO $$ BEGIN
    CREATE TYPE season AS ENUM ('winter', 'spring', 'summer', 'autumn');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
