ALTER TABLE crm.tenant_counters ENABLE ROW LEVEL SECURITY;
ALTER TABLE crm.tenant_counters FORCE ROW LEVEL SECURITY;

ALTER TABLE crm.lead_workflows ENABLE ROW LEVEL SECURITY;
ALTER TABLE crm.lead_workflows FORCE ROW LEVEL SECURITY;

ALTER TABLE crm.lead_stages ENABLE ROW LEVEL SECURITY;
ALTER TABLE crm.lead_stages FORCE ROW LEVEL SECURITY;

ALTER TABLE crm.leads ENABLE ROW LEVEL SECURITY;
ALTER TABLE crm.leads FORCE ROW LEVEL SECURITY;

ALTER TABLE crm.clients ENABLE ROW LEVEL SECURITY;
ALTER TABLE crm.clients FORCE ROW LEVEL SECURITY;

ALTER TABLE crm.projects ENABLE ROW LEVEL SECURITY;
ALTER TABLE crm.projects FORCE ROW LEVEL SECURITY;

ALTER TABLE crm.project_segments ENABLE ROW LEVEL SECURITY;
ALTER TABLE crm.project_segments FORCE ROW LEVEL SECURITY;

ALTER TABLE crm.sessions ENABLE ROW LEVEL SECURITY;
ALTER TABLE crm.sessions FORCE ROW LEVEL SECURITY;

ALTER TABLE crm.followups ENABLE ROW LEVEL SECURITY;
ALTER TABLE crm.followups FORCE ROW LEVEL SECURITY;

ALTER TABLE crm.packages ENABLE ROW LEVEL SECURITY;
ALTER TABLE crm.packages FORCE ROW LEVEL SECURITY;

ALTER TABLE crm.services ENABLE ROW LEVEL SECURITY;
ALTER TABLE crm.services FORCE ROW LEVEL SECURITY;

ALTER TABLE crm.quotations ENABLE ROW LEVEL SECURITY;
ALTER TABLE crm.quotations FORCE ROW LEVEL SECURITY;

CREATE POLICY tenant_counters_isolation ON crm.tenant_counters
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID);

CREATE POLICY lead_workflows_isolation ON crm.lead_workflows
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID);

CREATE POLICY lead_stages_isolation ON crm.lead_stages
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID);

CREATE POLICY leads_isolation ON crm.leads
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID);

CREATE POLICY clients_isolation ON crm.clients
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID);

CREATE POLICY projects_isolation ON crm.projects
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID);

CREATE POLICY project_segments_isolation ON crm.project_segments
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID);

CREATE POLICY sessions_isolation ON crm.sessions
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID);

CREATE POLICY followups_isolation ON crm.followups
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID);

CREATE POLICY packages_isolation ON crm.packages
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID);

CREATE POLICY services_isolation ON crm.services
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID);

CREATE POLICY quotations_isolation ON crm.quotations
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::UUID);

GRANT ALL ON ALL TABLES IN SCHEMA crm TO crm_app;
GRANT ALL ON ALL TABLES IN SCHEMA identity TO crm_app;
GRANT ALL ON ALL TABLES IN SCHEMA audit TO crm_app;

GRANT ALL ON ALL TABLES IN SCHEMA platform TO platform_app;
GRANT SELECT ON identity.users TO platform_app;
GRANT SELECT ON crm.tenants TO platform_app;
GRANT ALL ON ALL TABLES IN SCHEMA audit TO platform_app;

