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

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_pages_board_id ON pages(board_id);
CREATE INDEX IF NOT EXISTS idx_pages_board_order ON pages(board_id, order_idx);
CREATE INDEX IF NOT EXISTS idx_pages_date ON pages(date);