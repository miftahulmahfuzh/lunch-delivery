# Migration 007: Order History Test Data

## Purpose
This migration provides comprehensive test data for the Order History feature, specifically designed to test the date range filtering functionality.

## Test Data Overview
- **Total Orders**: 15 orders for Miftah (employee_id: 8) at Tuntun Sekuritas
- **Date Range**: December 2024 to September 2025
- **Payment Mix**: 10 PAID orders, 5 UNPAID orders  
- **Price Range**: Rp 80 - Rp 203

## Test Scenarios Covered

### Date Range Filters
- **"This Week"**: 5 orders (Sept 9-13, 2025)
- **"Last Week"**: Previous week orders from January 2025
- **"Last 2 Weeks"**: Combined recent orders
- **"This Month"**: Current month orders  
- **"Last Month"**: 4 orders (Dec 18-30, 2024)
- **"Custom Range"**: Any manually selected date range

### Realistic Data Patterns
- Various menu item combinations (rice, half-rice, proteins, vegetables)
- Mixed payment statuses for realistic testing
- Different order frequencies across dates
- Realistic Indonesian food prices in Rupiah

## How to Use

### Apply Migration
```bash
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f migrations/007_test_data_order_history.sql
```

### Verify Data
```bash
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -c "
SELECT COUNT(*) as total_test_orders 
FROM individual_orders 
WHERE employee_id = 8;
"
```

### Test the Feature
1. Login as Miftah: `miftahul.mahfuzh@tuntun.co.id`
2. Visit `/my-orders`
3. Try different date range shortcuts
4. Verify orders display correctly with menu item names and payment status

## Migration Safety
- Uses `ON CONFLICT DO NOTHING` to prevent duplicate data
- Can be run multiple times safely
- Only creates test data, doesn't modify existing schema
- Targets specific test user (employee_id: 8) to avoid affecting real data

## Clean Up (Optional)
To remove test data after testing:
```sql
DELETE FROM individual_orders WHERE employee_id = 8;
DELETE FROM order_sessions WHERE company_id = 1 AND date >= '2024-12-18';
```