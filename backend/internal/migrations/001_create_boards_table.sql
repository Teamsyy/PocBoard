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

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_boards_edit_token ON boards(edit_token);
CREATE INDEX IF NOT EXISTS idx_boards_public_token ON boards(public_token);