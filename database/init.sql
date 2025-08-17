-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create boards table
CREATE TABLE IF NOT EXISTS boards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    skin TEXT DEFAULT 'default',
    edit_token UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    public_token UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create pages table
CREATE TABLE IF NOT EXISTS pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    board_id UUID NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    order_idx INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

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
CREATE INDEX IF NOT EXISTS idx_pages_board_order ON pages(board_id, order_idx);
CREATE INDEX IF NOT EXISTS idx_pages_date ON pages(date);
CREATE INDEX IF NOT EXISTS idx_elements_page_z ON elements(page_id, z);