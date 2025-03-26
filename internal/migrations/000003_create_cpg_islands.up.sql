CREATE TABLE IF NOT EXISTS cpg_islands (
  id SERIAL PRIMARY KEY,
  sequence_id INT REFERENCES sequences(id) ON DELETE CASCADE, -- Links to sequences table
  start_position INT NOT NULL, 
  end_position INT NOT NULL, 
  gc_content DOUBLE PRECISION NOT NULL, -- GC% in this region
  obs_exp_cpg_ratio DOUBLE PRECISION NOT NULL, -- Observed/Expected CpG ratio
  created_at TIMESTAMP DEFAULT NOW()
);

-- CREATE TABLE IF NOT EXISTS promoters (
--   id SERIAL PRIMARY KEY,
--   sequence_id INT REFERENCES sequences(id) ON DELETE CASCADE, 
--   start_position INT NOT NULL, 
--   end_position INT NOT NULL, 
--   tata_box BOOLEAN DEFAULT FALSE,  -- True if TATA Box is found
--   gc_box BOOLEAN DEFAULT FALSE,    -- True if GC Box is found
--   caat_box BOOLEAN DEFAULT FALSE,  -- True if CAAT Box is found
--   created_at TIMESTAMP DEFAULT NOW()
-- );

-- CREATE TABLE IF NOT EXISTS gene_annotations (
--   id SERIAL PRIMARY KEY,
--   sequence_id INT REFERENCES sequences(id) ON DELETE CASCADE,
--   gene_name TEXT NOT NULL, 
--   start_position INT NOT NULL, 
--   end_position INT NOT NULL, 
--   strand CHAR(1) CHECK (strand IN ('+', '-')), -- Forward or Reverse strand
--   feature_type TEXT CHECK (feature_type IN ('exon', 'intron', 'CDS', 'UTR', 'gene')), -- Gene feature type
--   created_at TIMESTAMP DEFAULT NOW()
-- );

-- CREATE TABLE IF NOT EXISTS mutations (
--   id SERIAL PRIMARY KEY,
--   sequence_id INT REFERENCES sequences(id) ON DELETE CASCADE,
--   position INT NOT NULL, 
--   mutation_type TEXT CHECK (mutation_type IN ('SNP', 'insertion', 'deletion')), -- Variant type
--   reference_base TEXT NOT NULL, 
--   mutated_base TEXT NOT NULL, 
--   created_at TIMESTAMP DEFAULT NOW()
-- );

-- CREATE TABLE IF NOT EXISTS sequence_metadata (
--   id SERIAL PRIMARY KEY,
--   sequence_id INT REFERENCES sequences(id) ON DELETE CASCADE,
--   sequencing_technology TEXT, -- e.g., Illumina, PacBio, Oxford Nanopore
--   coverage DOUBLE PRECISION, -- Depth of sequencing coverage
--   assembly_method TEXT, -- e.g., SPAdes, Canu, Flye
--   reference_genome TEXT, -- If the sequence was aligned to a reference
--   created_at TIMESTAMP DEFAULT NOW()
-- );
