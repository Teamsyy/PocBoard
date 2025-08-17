-- Create elements table
CREATE TABLE IF NOT EXISTS elements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_id UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    kind TEXT NOT NULL CHECK (kind IN ('text','image','sticker','shape')),
    x FLOAT NOT NULL,
    y FLOAT NOT NULL,
    w FLOAT NOT NULL,
    h FLOAT NOT NULL,
    rotation FLOAT DEFAULT 0,
    z INTEGER DEFAULT 0,
    payload JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_elements_page_id ON elements(page_id);
CREATE INDEX IF NOT EXISTS idx_elements_page_z ON elements(page_id, z);
CREATE INDEX IF NOT EXISTS idx_elements_kind ON elements(kind);