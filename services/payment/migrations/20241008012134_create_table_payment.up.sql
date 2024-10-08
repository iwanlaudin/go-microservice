CREATE TABLE "payments" (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id   INT NOT NULL,
    user_id     INT NOT NULL,
    amount      DECIMAL(10, 2) NOT NULL,
    status      VARCHAR(50) NOT NULL,
    date        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by  VARCHAR(256),
    created_at  TIMESTAMP WITH TIME ZONE,
    updated_by  VARCHAR(256),
    updated_at  TIMESTAMP WITH TIME ZONE,
    is_deleted  BOOLEAN DEFAULT FALSE
);
