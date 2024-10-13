CREATE TABLE IF NOT EXISTS "urls" (
	original_url TEXT PRIMARY KEY NOT NULL,
	short_url TEXT NOT NULL,
	created_at TEXT NOT NULL
);

-- Create index on the short_url
CREATE index IF NOT EXISTS idx_short ON urls (short_url);

