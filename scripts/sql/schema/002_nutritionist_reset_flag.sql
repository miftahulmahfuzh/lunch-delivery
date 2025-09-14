-- migrations/003_nutritionist_reset_flag.sql

-- Add reset flag to daily_menus table
ALTER TABLE daily_menus ADD COLUMN nutritionist_reset BOOLEAN DEFAULT FALSE;

-- Add tracking table for users who used nutritionist selection
CREATE TABLE nutritionist_user_selections (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER REFERENCES employees(id),
    date DATE NOT NULL,
    order_id INTEGER REFERENCES individual_orders(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(employee_id, date)
);

CREATE INDEX idx_nutritionist_user_selections_date ON nutritionist_user_selections(date);
CREATE INDEX idx_nutritionist_user_selections_employee ON nutritionist_user_selections(employee_id);