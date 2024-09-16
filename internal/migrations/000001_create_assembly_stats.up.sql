CREATE TABLE IF NOT EXISTS assembly_stats (
  id                    SERIAL PRIMARY KEY,
  assembly_id           INTEGER NOT NULL,
  contigL50             INTEGER NOT NULL,
  contigN50             INTEGER NOT NULL,
  total_sequence_length BIGINT NOT NULL,
  total_ungapped_length BIGINT NOT NULL,
  gc_count              BIGINT NOT NULL,
  gc_percent            DOUBLE PRECISION NOT NULL,
  genome_coverage       VARCHAR(50),
  number_of_contigs     INTEGER NOT NULL,
  created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
