-- Delete today's order session for a specific company (for testing purposes)
-- This script helps test the menu validation feature by removing any existing order session for today
-- Usage: Replace COMPANY_ID with the actual company ID you want to test with

-- First, let's see what sessions exist for today
SELECT
    os.id,
    os.company_id,
    c.name as company_name,
    os.date,
    os.status,
    os.created_at
FROM order_sessions os
JOIN companies c ON os.company_id = c.id
WHERE os.date = CURRENT_DATE;

-- Delete order session for company ID 1 (Tuntun company) on today's date
-- Change the company_id value (1) to match the company you want to test with
DELETE FROM order_sessions
WHERE company_id = 1
AND date = CURRENT_DATE;

-- Verify the deletion
SELECT
    os.id,
    os.company_id,
    c.name as company_name,
    os.date,
    os.status,
    os.created_at
FROM order_sessions os
JOIN companies c ON os.company_id = c.id
WHERE os.date = CURRENT_DATE;

-- If you want to delete for ALL companies today (use with caution!)
-- DELETE FROM order_sessions WHERE date = CURRENT_DATE;