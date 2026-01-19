CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE organization_members (
    organization_id UUID NOT NULL REFERENCES organizations(id),
    user_id UUID NOT NULL REFERENCES users(id),
    role VARCHAR(50) NOT NULL,
    joined_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (organization_id, user_id)
);

-- Backfill: Create organization for each existing user and link them
DO $$
DECLARE
    user_rec RECORD;
    new_org_id UUID;
    org_slug TEXT;
BEGIN
    FOR user_rec IN SELECT id, name FROM users LOOP
        -- Generate slug: slugify(name) + "-" + randomHex(4)
        org_slug := regexp_replace(lower(COALESCE(user_rec.name, 'My Workspace')), '[^a-z0-9]+', '-', 'g') || '-' || substring(md5(random()::text) from 1 for 4);

        INSERT INTO organizations (name, slug)
        VALUES (COALESCE(user_rec.name, 'My Workspace'), org_slug)
        RETURNING id INTO new_org_id;

        INSERT INTO organization_members (organization_id, user_id, role)
        VALUES (new_org_id, user_rec.id, 'owner');
    END LOOP;
END $$;
