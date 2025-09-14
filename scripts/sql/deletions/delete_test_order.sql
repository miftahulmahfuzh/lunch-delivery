-- delete_test_order.sql
-- Quick script to delete Jemmy's order for today to test the footer edge case

-- Delete Jemmy's order for today
DELETE FROM individual_orders
WHERE employee_id = 1  -- Jemmy
AND session_id = (
    SELECT id FROM order_sessions
    WHERE date = CURRENT_DATE
    AND company_id = 1
);

-- Verify the deletion
SELECT
    CASE
        WHEN COUNT(*) = 0 THEN '✅ SUCCESS: Jemmy has no order for today - ready to test footer!'
        ELSE '❌ Order still exists - deletion failed'
    END as test_status,
    COUNT(*) as remaining_orders,
    CURRENT_DATE as test_date,
    'Login as jemmy@techcorp.com and click "Edit Orders" in footer' as next_step
FROM individual_orders io
JOIN order_sessions os ON io.session_id = os.id
WHERE io.employee_id = 1 AND os.date = CURRENT_DATE;