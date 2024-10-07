CREATE TABLE tickets (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id    UUID NOT NULL,
    user_id     UUID NOT NULL,
    quantity    INT NOT NULL,
    status      VARCHAR(50) NOT NULL,
    created_by  VARCHAR(256),
    created_at  TIMESTAMP WITH TIME ZONE,
    updated_by  VARCHAR(256),
    updated_at  TIMESTAMP WITH TIME ZONE,
    is_deleted  BOOLEAN DEFAULT FALSE
);
