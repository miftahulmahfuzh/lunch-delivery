-- Add status field to individual_orders table for preparation tracking
ALTER TABLE individual_orders
ADD COLUMN status VARCHAR(50) DEFAULT 'PENDING';

-- Update existing orders to PENDING status
UPDATE individual_orders SET status = 'PENDING' WHERE status IS NULL;

-- Create index for status queries
CREATE INDEX idx_individual_orders_status ON individual_orders(status);