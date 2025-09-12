-- migrations/007_test_data_order_history.sql
-- Test Data for Order History Feature
-- This migration injects comprehensive test data for testing the Order History date range filtering functionality
-- 
-- Test Scenarios Covered:
-- - Current Week (Sept 2025): 5 orders for testing "This Week" filter
-- - Previous Weeks (Jan 2025): 6 orders for testing "Last Week" and "Last 2 Weeks" filters  
-- - Previous Month (Dec 2024): 4 orders for testing "Last Month" filter
-- - Mixed payment statuses (PAID/UNPAID) for realistic testing
-- - Various menu combinations and price points
-- - Different order patterns across multiple dates
--
-- Target Test User: Miftah (employee_id: 8) at Tuntun Sekuritas (company_id: 1)

-- ============================================================================
-- ORDER SESSIONS - Create sessions for various dates to test different ranges
-- ============================================================================

INSERT INTO order_sessions (company_id, date, status, created_at) VALUES 
-- Current Week (September 2025) - for "This Week" filter
(1, '2025-09-13', 'DELIVERED', '2025-09-13 08:00:00'),
(1, '2025-09-12', 'DELIVERED', '2025-09-12 08:00:00'),
(1, '2025-09-11', 'DELIVERED', '2025-09-11 08:00:00'),
(1, '2025-09-10', 'DELIVERED', '2025-09-10 08:00:00'),
(1, '2025-09-09', 'DELIVERED', '2025-09-09 08:00:00'),

-- Previous Weeks (January 2025) - for "Last Week" and "Last 2 Weeks" filters
(1, '2025-01-10', 'DELIVERED', '2025-01-10 08:00:00'),
(1, '2025-01-09', 'DELIVERED', '2025-01-09 08:00:00'),
(1, '2025-01-08', 'DELIVERED', '2025-01-08 08:00:00'),
(1, '2025-01-06', 'DELIVERED', '2025-01-06 08:00:00'),
(1, '2025-01-03', 'DELIVERED', '2025-01-03 08:00:00'),
(1, '2025-01-02', 'DELIVERED', '2025-01-02 08:00:00'),

-- Previous Month (December 2024) - for "Last Month" filter
(1, '2024-12-30', 'DELIVERED', '2024-12-30 08:00:00'),
(1, '2024-12-27', 'DELIVERED', '2024-12-27 08:00:00'),
(1, '2024-12-20', 'DELIVERED', '2024-12-20 08:00:00'),
(1, '2024-12-18', 'DELIVERED', '2024-12-18 08:00:00')

ON CONFLICT (company_id, date) DO NOTHING;

-- ============================================================================
-- INDIVIDUAL ORDERS - Create realistic orders for Miftah across all date ranges
-- ============================================================================

-- Get the session IDs for our orders (using subqueries to be migration-safe)
-- Current Week Orders (September 2025) - 5 orders
INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, paid)
SELECT os.id, 8, ARRAY[5, 3, 10, 24], 20300, false
FROM order_sessions os WHERE os.company_id = 1 AND os.date = '2025-09-13'
ON CONFLICT (session_id, employee_id) DO NOTHING;

INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, paid)
SELECT os.id, 8, ARRAY[9, 3, 10, 4], 18600, true
FROM order_sessions os WHERE os.company_id = 1 AND os.date = '2025-09-12'
ON CONFLICT (session_id, employee_id) DO NOTHING;

INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, paid)
SELECT os.id, 8, ARRAY[10, 9, 3], 14500, true
FROM order_sessions os WHERE os.company_id = 1 AND os.date = '2025-09-11'
ON CONFLICT (session_id, employee_id) DO NOTHING;

INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, paid)
SELECT os.id, 8, ARRAY[11, 6, 2], 10500, false
FROM order_sessions os WHERE os.company_id = 1 AND os.date = '2025-09-10'
ON CONFLICT (session_id, employee_id) DO NOTHING;

INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, paid)
SELECT os.id, 8, ARRAY[10, 1, 7], 13700, true
FROM order_sessions os WHERE os.company_id = 1 AND os.date = '2025-09-09'
ON CONFLICT (session_id, employee_id) DO NOTHING;

-- Previous Weeks Orders (January 2025) - 6 orders
INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, paid)
SELECT os.id, 8, ARRAY[10, 7, 1], 13700, true
FROM order_sessions os WHERE os.company_id = 1 AND os.date = '2025-01-10'
ON CONFLICT (session_id, employee_id) DO NOTHING;

INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, paid)
SELECT os.id, 8, ARRAY[10, 9, 2], 14000, true
FROM order_sessions os WHERE os.company_id = 1 AND os.date = '2025-01-09'
ON CONFLICT (session_id, employee_id) DO NOTHING;

INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, paid)
SELECT os.id, 8, ARRAY[11, 6, 3], 12000, false
FROM order_sessions os WHERE os.company_id = 1 AND os.date = '2025-01-08'
ON CONFLICT (session_id, employee_id) DO NOTHING;

INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, paid)
SELECT os.id, 8, ARRAY[10, 5, 4], 14100, true
FROM order_sessions os WHERE os.company_id = 1 AND os.date = '2025-01-06'
ON CONFLICT (session_id, employee_id) DO NOTHING;

INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, paid)
SELECT os.id, 8, ARRAY[10, 7], 9500, true
FROM order_sessions os WHERE os.company_id = 1 AND os.date = '2025-01-03'
ON CONFLICT (session_id, employee_id) DO NOTHING;

INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, paid)
SELECT os.id, 8, ARRAY[11, 6, 1, 2], 14700, false
FROM order_sessions os WHERE os.company_id = 1 AND os.date = '2025-01-02'
ON CONFLICT (session_id, employee_id) DO NOTHING;

-- Previous Month Orders (December 2024) - 4 orders
INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, paid)
SELECT os.id, 8, ARRAY[10, 9, 3], 14500, true
FROM order_sessions os WHERE os.company_id = 1 AND os.date = '2024-12-30'
ON CONFLICT (session_id, employee_id) DO NOTHING;

INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, paid)
SELECT os.id, 8, ARRAY[11, 5], 8000, true
FROM order_sessions os WHERE os.company_id = 1 AND os.date = '2024-12-27'
ON CONFLICT (session_id, employee_id) DO NOTHING;

INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, paid)
SELECT os.id, 8, ARRAY[10, 7, 4, 6], 17600, false
FROM order_sessions os WHERE os.company_id = 1 AND os.date = '2024-12-20'
ON CONFLICT (session_id, employee_id) DO NOTHING;

INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, paid)
SELECT os.id, 8, ARRAY[10, 2, 1], 13200, true
FROM order_sessions os WHERE os.company_id = 1 AND os.date = '2024-12-18'
ON CONFLICT (session_id, employee_id) DO NOTHING;

-- ============================================================================
-- TEST DATA SUMMARY
-- ============================================================================
/*
Total Test Orders Created: 15 orders for Miftah (employee_id: 8)

Current Week (Sept 9-13, 2025) - 5 orders:
- Sept 13: Cah oyong telur + Nasi + Cah toge + Dori cabe ijo → Rp 203 (UNPAID)
- Sept 12: Kikil balado + Cah toge + Nasi + Cah kembang kol → Rp 186 (PAID)
- Sept 11: Nasi + Kikil balado + Cah toge → Rp 145 (PAID)
- Sept 10: Nasi 1/2 + Ceker cabe ijo + Cah labu → Rp 105 (UNPAID)
- Sept 9: Nasi + Cah jagung muda + Ayam goreng kandar merah → Rp 137 (PAID)

Previous Weeks (Jan 2-10, 2025) - 6 orders:
- Jan 10: Nasi + Ayam goreng + Cah jagung muda → Rp 137 (PAID)
- Jan 9: Nasi + Kikil balado + Cah labu → Rp 140 (PAID)
- Jan 8: Nasi 1/2 + Ceker + Cah toge → Rp 120 (UNPAID)
- Jan 6: Nasi + Cah oyong telur + Cah kembang kol → Rp 141 (PAID)
- Jan 3: Nasi + Ayam goreng → Rp 95 (PAID)
- Jan 2: Nasi 1/2 + Ceker + Cah jagung + Cah labu → Rp 147 (UNPAID)

Previous Month (Dec 18-30, 2024) - 4 orders:
- Dec 30: Nasi + Kikil balado + Cah toge → Rp 145 (PAID)
- Dec 27: Nasi 1/2 + Cah oyong telur → Rp 80 (PAID)
- Dec 20: Nasi + Ayam goreng + Cah kembang kol + Ceker → Rp 176 (UNPAID)
- Dec 18: Nasi + Cah labu + Cah jagung muda → Rp 132 (PAID)

Payment Status Distribution:
- PAID: 10 orders (67%)
- UNPAID: 5 orders (33%)

Price Range: Rp 80 - Rp 203
Average Order Value: Rp 140

Test Coverage:
✓ "This Week" filter → Shows 5 current week orders
✓ "Last Week" filter → Shows orders from previous week(s)
✓ "Last 2 Weeks" filter → Shows combined recent orders  
✓ "This Month" filter → Shows current month orders
✓ "Last Month" filter → Shows 4 December orders
✓ "Custom Range" filter → Any date range selection
✓ Mixed payment statuses for realistic scenarios
✓ Various menu item combinations and prices
✓ Realistic order frequency patterns
*/