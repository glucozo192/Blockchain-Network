DROP TABLE IF EXISTS nodes;

CREATE TABLE IF NOT EXISTS nodes (
  id text PRIMARY KEY,
  address text NOT NULL,
  created_at timestamptz,
  updated_at timestamptz NULL,
  deleted_at timestamptz
);

DROP TABLE IF EXISTS blocks;
CREATE TABLE IF NOT EXISTS blocks (
  id text PRIMARY KEY,
  data INTEGER[],
  height SERIAL,
  created_at timestamptz,
  updated_at timestamptz NULL,
  deleted_at timestamptz
);

DROP TABLE IF EXISTS markers;

CREATE TABLE IF NOT EXISTS markers (
  block_id text PRIMARY KEY,
  created_at timestamptz,
  updated_at timestamptz NULL,
  deleted_at timestamptz
);

