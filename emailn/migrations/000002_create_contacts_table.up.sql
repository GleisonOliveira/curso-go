CREATE TABLE IF NOT EXISTS contacts (
    id UUID DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    campaign_id UUID NOT NULL,
    PRIMARY KEY (id, campaign_id),
    CONSTRAINT fk_contacts_campaign FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
);
