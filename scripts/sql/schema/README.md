# Schema Scripts

This directory contains SQL scripts for database schema management - table creation, structure changes, and database setup.

## Files Overview

### Core Schema
- **001_initial.sql** - Initial database schema setup (creates all base tables)
  - Creates: `menu_items`, `companies`, `employees`, `daily_menus`, `order_sessions`, `individual_orders`
  - Includes indexes for performance optimization

### Feature Extensions
- **003_nutritionist_reset_flag.sql** - Adds nutritionist-related functionality
- **005_stock_empty_and_notifications.sql** - Stock management system
  - Creates: `stock_empty_items`, `user_stock_empty_notifications`
- **008_nutritionist_selections.sql** - Nutritionist selection features

## Usage

### Initial Setup
Run the initial schema first:
```bash
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f schema/001_initial.sql
```

### Feature Additions
Run feature-specific scripts in order:
```bash
# Add nutritionist features
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f schema/003_nutritionist_reset_flag.sql
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f schema/008_nutritionist_selections.sql

# Add stock management
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f schema/005_stock_empty_and_notifications.sql
```

## Notes
- Always run `001_initial.sql` first before any other schema scripts
- These scripts create tables and modify database structure
- Use with caution in production environments
- Consider backing up data before running schema changes