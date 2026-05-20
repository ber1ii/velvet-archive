-- Create tables
CREATE TABLE series (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    cover_color TEXT NOT NULL, -- Hex code for 3D meshes
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE lore_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id UUID REFERENCES series(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    category TEXT NOT NULL, -- Character, Location, Faction, Artifact, Event
    content TEXT NOT NULL,
    metadata JSONB,          -- For system-specific attributes
    affinity TEXT[],         -- Tags like {Agi, Zio}
    search_vec TSVECTOR,     -- Handled by the trigger below
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE lore_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID REFERENCES lore_entries(id) ON DELETE CASCADE,
    target_id UUID REFERENCES lore_entries(id) ON DELETE CASCADE,
    label TEXT,              -- Relationship title (e.g., 'Mentored by')
    is_revealed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE admin_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Full-Text Search Trigger Logic
CREATE OR REPLACE FUNCTION lore_entries_tsvector_trigger() 
RETURNS TRIGGER AS $$
BEGIN
  NEW.search_vec := to_tsvector('english', coalesce(NEW.title, '') || ' ' || coalesce(NEW.content, ''));
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tsvectorupdate 
BEFORE INSERT OR UPDATE ON lore_entries 
FOR EACH ROW EXECUTE FUNCTION lore_entries_tsvector_trigger();