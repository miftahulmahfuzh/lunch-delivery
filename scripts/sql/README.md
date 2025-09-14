# SQL Scripts Directory

This directory contains organized SQL scripts for the Tuntun Lunch Delivery system, grouped by functionality for better maintainability and clarity.

## Directory Structure

```
scripts/sql/
├── schema/          # Database structure & table creation
├── seeds/           # Data insertion & initial setup
├── updates/         # Data modifications & schema changes
├── deletions/       # Data cleanup & testing scenarios
└── README.md        # This overview file
```

## Quick Start

### 1. Fresh Database Setup
For a completely new database installation:

```bash
# 1. Create schema (tables, indexes)
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f schema/001_initial.sql

# 2. Add feature tables
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f schema/003_nutritionist_reset_flag.sql
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f schema/005_stock_empty_and_notifications.sql
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f schema/008_nutritionist_selections.sql

# 3. Seed initial data
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f seeds/002_menu_items_seed.sql
```

### 2. Development Environment
Add test data for development and testing:

```bash
# Add test orders for Order History feature testing
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f seeds/007_test_data_order_history.sql
```

### 3. Maintenance Operations
Use update and deletion scripts as needed:

```bash
# Example: Convert price format (one-time operation)
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f updates/004_update_price_to_rupiah.sql

# Example: Clean test data
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f deletions/delete_test_order.sql
```

## Directory Details

### 📋 schema/
**Purpose**: Database structure management
- Table creation scripts
- Index definitions
- Schema modifications
- **Start here** for new installations

### 🌱 seeds/
**Purpose**: Data population
- Menu items catalog
- Test data for features
- Reference data
- Safe to run multiple times

### 🔄 updates/
**Purpose**: Data and schema modifications
- Price conversions
- Schema alterations
- **⚠️ Use with caution** - modifies existing data

### 🗑️ deletions/
**Purpose**: Data cleanup and testing
- Test data removal
- Session cleanup
- Edge case setup
- **⚠️ Production risk** - always backup first

## Best Practices

### Development Workflow
1. **Schema** → **Seeds** → **Test** → **Update/Delete** as needed
2. Always run in development environment first
3. Use version control for all script changes
4. Document any custom modifications

### Production Safety
- **Always backup** before running updates or deletions
- Test scripts in staging environment first
- Use transactions for complex operations
- Document all production changes

### Database Connection
All scripts use the standard connection:
```bash
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f <script_path>
```

## Common Commands

### Verification
```bash
# Check table structure
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -c "\dt"

# Check menu items count
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -c "SELECT COUNT(*) FROM menu_items;"

# Check test user orders
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -c "SELECT COUNT(*) FROM individual_orders WHERE employee_id = 8;"
```

### Backup
```bash
# Full backup
pg_dump -h localhost -p 5432 -U lunch_user lunch_delivery > backup_$(date +%Y%m%d_%H%M%S).sql

# Restore
psql -h localhost -p 5432 -U lunch_user -d lunch_delivery < backup_file.sql
```

---

For specific details about each directory, see the README.md file in each subdirectory.