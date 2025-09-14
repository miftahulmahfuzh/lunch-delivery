-- migrations/009_test_footer_edge_case.sql
-- Test script to simulate the footer "Edit Orders" edge case
-- This script helps test the scenario where a user has no order for today
-- and clicks "Edit Orders" in the footer

-- First, let's ensure we have a test order for employee ID 1 (Jemmy) for today
-- We'll create it if it doesn't exist, then provide deletion script

-- Step 1: Get today's order session ID
WITH today_session AS (
    SELECT id as session_id, company_id
    FROM order_sessions
    WHERE date = CURRENT_DATE
    AND company_id = 1
    LIMIT 1
)
-- Step 2: Create a test order if it doesn't exist
INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, paid)
SELECT
    ts.session_id,
    1 as employee_id, -- Jemmy (jemmy@techcorp.com)
    ARRAY[1, 2, 3] as menu_item_ids, -- Select first 3 menu items
    75000 as total_price, -- 750 rupiah total
    false as paid
FROM today_session ts
WHERE NOT EXISTS (
    SELECT 1 FROM individual_orders io
    JOIN order_sessions os ON io.session_id = os.id
    WHERE io.employee_id = 1
    AND os.date = CURRENT_DATE
);

-- Confirm the order was created/exists
SELECT
    io.id as order_id,
    e.name as employee_name,
    e.email,
    io.menu_item_ids,
    io.total_price,
    io.paid,
    os.date as order_date,
    os.status as session_status
FROM individual_orders io
JOIN employees e ON io.employee_id = e.id
JOIN order_sessions os ON io.session_id = os.id
WHERE e.id = 1 AND os.date = CURRENT_DATE;

-- ==================================================================
-- TESTING INSTRUCTIONS:
-- ==================================================================
--
-- 1. Run this script to create a test order for Jemmy
-- 2. Login as Jemmy (jemmy@techcorp.com) and verify order exists in /my-orders
-- 3. Then run the DELETE script below to remove the order
-- 4. Stay logged in as Jemmy and click "Edit Orders" in the footer
-- 5. The footer should redirect to /order which should redirect to proper order form
--
-- ==================================================================
-- DELETE SCRIPT (run this after Step 2 above):
-- ==================================================================
/*

DELETE FROM individual_orders
WHERE employee_id = 1
AND session_id = (
    SELECT id FROM order_sessions
    WHERE date = CURRENT_DATE
    AND company_id = 1
);

-- Verify deletion
SELECT
    COUNT(*) as remaining_orders,
    CURRENT_DATE as test_date
FROM individual_orders io
JOIN order_sessions os ON io.session_id = os.id
WHERE io.employee_id = 1 AND os.date = CURRENT_DATE;

*/