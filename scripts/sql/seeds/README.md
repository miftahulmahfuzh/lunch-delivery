# Seeds Scripts

This directory contains SQL scripts for data insertion - initial data, test data, and sample records.

## Files Overview

### Production Data
- **002_menu_items_seed.sql** - Complete menu items catalog
  - Includes 60+ Indonesian food items with realistic pricing
  - Categories: Premium dishes, standard dishes, sides, vegetables, staples, beverages
  - Prices in Rupiah (Rp 2,000 - Rp 20,000)

### Test Data
- **007_test_data_order_history.sql** - Comprehensive test data for Order History feature
  - 15 orders for test user Miftah (employee_id: 8)
  - Date range: December 2024 to September 2025
  - Mixed payment statuses (PAID/UNPAID)
  - Tests various date range filters ("This Week", "Last Month", etc.)

## Usage

### Production Setup
After schema is created, seed the menu items:
```bash
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f seeds/002_menu_items_seed.sql
```

### Test Environment Setup
Add test data for development and testing:
```bash
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f seeds/007_test_data_order_history.sql
```

### Verification
Verify seed data was inserted:
```bash
# Check menu items
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -c "SELECT COUNT(*) FROM menu_items;"

# Check test orders
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -c "SELECT COUNT(*) FROM individual_orders WHERE employee_id = 8;"
```

## Notes
- Run schema scripts first before seeding data
- Uses `ON CONFLICT DO NOTHING` to prevent duplicate insertions
- Test data is safe to run multiple times
- Test user login: `miftahul.mahfuzh@tuntun.co.id` for testing Order History feature