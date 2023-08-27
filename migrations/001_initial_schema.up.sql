
CREATE TYPE t_organization_status AS ENUM ('created', 'updated', 'unverified', 'verified');

CREATE TABLE IF NOT EXISTS organizations (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    processor_id VARCHAR2(255) NOT NULL,
    name TEXT NOT NULL,
    status t_organization_status NOT NULL DEFAULT 'created',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);