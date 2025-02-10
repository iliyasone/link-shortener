CREATE TABLE IF NOT EXISTS url_mappings (
    id SERIAL PRIMARY KEY,
    short_url VARCHAR(10) NOT NULL UNIQUE,
    original_url TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create an index for fast lookups of original URLs. 
-- This is breaking normalization for speed up
CREATE INDEX idx_original_url ON url_mappings(original_url);
