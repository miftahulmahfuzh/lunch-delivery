-- migrations/006_remove_global_stock_empty.sql

-- Drop the global stock_empty_items table since we're using user-specific tracking only
DROP TABLE IF EXISTS stock_empty_items;

-- The user_stock_empty_notifications table will be our only source of truth for stock status
-- This table tracks empty stock per individual order (user-specific)