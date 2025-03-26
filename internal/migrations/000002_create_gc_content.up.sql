CREATE TABLE IF NOT EXISTS gc_content (
  id SERIAL PRIMARY KEY,
  sequence_id INT REFERENCES sequences(id) ON DELETE CASCADE, -- Links to sequences table
  start_position INT NOT NULL,     -- Start position of the analyzed region
  end_position INT NOT NULL,       -- End position of the analyzed region
  gc_percentage DOUBLE PRECISION NOT NULL, -- GC% in this region
  region_type TEXT CHECK (region_type IN ('rich', 'poor')), -- GC classification
  created_at TIMESTAMP DEFAULT NOW()
);
