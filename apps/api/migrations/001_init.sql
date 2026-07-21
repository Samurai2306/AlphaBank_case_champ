-- Placeholder for full Postgres mode (USE_MEMORY_STORE=0).
-- Demo MVP uses in-memory seeded store by default.

CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS profiles (
  user_id TEXT PRIMARY KEY,
  display_name TEXT NOT NULL,
  business_sphere TEXT,
  city TEXT,
  monthly_revenue_estimate NUMERIC,
  tax_regime TEXT,
  cjm_level INT DEFAULT 1,
  persona_key TEXT
);
