-- Convert price representation from cents to rupiah
-- Current: User enters 15000, stored as 1500000 (cents), displayed as 15000
-- Target: User enters 4000, stored as 4000 (rupiah), displayed as 4000

-- Update existing prices by dividing by 100 to convert from cents to rupiah
UPDATE menu_items SET price = price / 100;
UPDATE individual_orders SET total_price = total_price / 100;

-- Update column comments
COMMENT ON COLUMN menu_items.price IS 'stored in rupiah';
COMMENT ON COLUMN individual_orders.total_price IS 'stored in rupiah';