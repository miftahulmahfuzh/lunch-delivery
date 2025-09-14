-- migration: 009_password_reset_tokens.sql
-- Add password reset tokens table for forgot password functionality

CREATE TABLE password_reset_tokens (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER REFERENCES employees(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for efficient token lookup
CREATE INDEX idx_password_reset_tokens_token ON password_reset_tokens(token);
CREATE INDEX idx_password_reset_tokens_employee ON password_reset_tokens(employee_id);
CREATE INDEX idx_password_reset_tokens_expires ON password_reset_tokens(expires_at);