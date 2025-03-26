CREATE TABLE IF NOT EXISTS sequences (
  id SERIAL PRIMARY KEY,
  accession TEXT NOT NULL,          -- CM034354.1 (Unique GenBank accession number)
  species TEXT,                     -- Mytilus edulis (Organism name)
  chromosome TEXT,                  -- Chromosome name or scaffold ID (if available)
  description TEXT,                 -- Metadata (e.g., "whole genome shotgun sequence")
  sequence TEXT NOT NULL,           -- Raw DNA sequence
  start_position INT DEFAULT 1,     -- Always starts at position 1
  end_position INT NOT NULL,        -- Computed as LENGTH(sequence)
  created_at TIMESTAMP DEFAULT NOW()
);
