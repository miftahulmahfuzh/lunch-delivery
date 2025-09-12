-- migrations/005_stock_empty_and_notifications.sql

-- Table to track stock empty items per day
CREATE TABLE stock_empty_items (
    id SERIAL PRIMARY KEY,
    menu_item_id INTEGER REFERENCES menu_items(id),
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(menu_item_id, date)
);

-- Table to track individual stock empty notifications per user
CREATE TABLE user_stock_empty_notifications (
    id SERIAL PRIMARY KEY,
    individual_order_id INTEGER REFERENCES individual_orders(id),
    menu_item_id INTEGER REFERENCES menu_items(id),
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(individual_order_id, menu_item_id)
);

-- Table to store user notifications
CREATE TABLE user_notifications (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER REFERENCES employees(id),
    notification_type VARCHAR(50) NOT NULL, -- STOCK_EMPTY, PAID, SESSION_CLOSED, MENU_UPDATED
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    redirect_url VARCHAR(255),
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX idx_stock_empty_items_date ON stock_empty_items(date);
CREATE INDEX idx_stock_empty_items_menu_item ON stock_empty_items(menu_item_id);
CREATE INDEX idx_user_stock_empty_notifications_order ON user_stock_empty_notifications(individual_order_id);
CREATE INDEX idx_user_notifications_employee ON user_notifications(employee_id);
CREATE INDEX idx_user_notifications_unread ON user_notifications(employee_id, is_read) WHERE is_read = false;