CREATE TABLE "userToken"
(
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    refresh_token   VARCHAR(256) NOT NULL,
    expiry_at       TIMESTAMP WITH TIME ZONE,
    used_at         TIMESTAMP WITH TIME ZONE,
    is_used         BOOLEAN DEFAULT FALSE,
    created_by      VARCHAR(256),
    created_at      TIMESTAMP WITH TIME ZONE,
    updated_by      VARCHAR(256),
    updated_at      TIMESTAMP WITH TIME ZONE,
    is_deleted      BOOLEAN DEFAULT FALSE
)