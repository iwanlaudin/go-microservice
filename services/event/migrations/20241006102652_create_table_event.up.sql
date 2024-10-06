CREATE TABLE events (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(255) NOT NULL,
    date                DATE NOT NULL,
    location            VARCHAR(255) NOT NULL,
    available_tickets   INT NOT NULL,
    created_by          VARCHAR(256),
    created_at          TIMESTAMP WITH TIME ZONE,
    updated_by          VARCHAR(256),
    updated_at          TIMESTAMP WITH TIME ZONE,
    is_deleted          BOOLEAN DEFAULT FALSE
);
