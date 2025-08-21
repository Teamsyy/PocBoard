-- Add visible and locked columns to elements table (only if they don't exist)
DO $$
BEGIN
    -- Add visible column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'elements' AND column_name = 'visible') THEN
        ALTER TABLE elements ADD COLUMN visible BOOLEAN DEFAULT true;
    END IF;
    
    -- Add locked column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'elements' AND column_name = 'locked') THEN
        ALTER TABLE elements ADD COLUMN locked BOOLEAN DEFAULT false;
    END IF;
END $$;

-- Create indexes for better query performance (only if they don't exist)
DO $$
BEGIN
    -- Create visible index if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_elements_visible') THEN
        CREATE INDEX idx_elements_visible ON elements(visible);
    END IF;
    
    -- Create locked index if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_elements_locked') THEN
        CREATE INDEX idx_elements_locked ON elements(locked);
    END IF;
END $$;
