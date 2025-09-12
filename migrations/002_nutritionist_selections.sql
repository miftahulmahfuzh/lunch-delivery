-- migrations/002_nutritionist_selections.sql
CREATE TABLE nutritionist_selections (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL UNIQUE,
    menu_item_ids BIGINT[] NOT NULL,
    selected_indices INTEGER[] NOT NULL,
    reasoning TEXT,
    nutritional_summary JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_nutritionist_selections_date ON nutritionist_selections(date);