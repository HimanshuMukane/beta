CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
DECLARE
    t text;
BEGIN
    FOR t IN 
        SELECT table_schema || '.' || table_name 
        FROM information_schema.columns 
        WHERE column_name = 'updated_at' 
          AND table_schema IN ('crm', 'identity', 'platform')
    LOOP
        EXECUTE format('DROP TRIGGER IF EXISTS trg_set_updated_at ON %s;', t);
        EXECUTE format('CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON %s FOR EACH ROW EXECUTE FUNCTION set_updated_at();', t);
    END LOOP;
END
$$;

CREATE INDEX IF NOT EXISTS idx_clients_name_trgm ON crm.clients USING GIN (first_name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_leads_contact_trgm ON crm.leads USING GIN (contact_name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_clients_search_fts ON crm.clients USING GIN (to_tsvector('simple', coalesce(first_name, '') || ' ' || coalesce(last_name, '') || ' ' || coalesce(email, '') || ' ' || coalesce(phone, '')));
CREATE INDEX IF NOT EXISTS idx_leads_search_fts ON crm.leads USING GIN (to_tsvector('simple', coalesce(title, '') || ' ' || coalesce(contact_name, '') || ' ' || coalesce(notes, '')));

CREATE INDEX IF NOT EXISTS idx_followups_due_queue ON crm.followups(due_date) WHERE status = 'pending';

