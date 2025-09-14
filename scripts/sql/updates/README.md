# Updates Scripts

This directory contains SQL scripts for data updates, modifications, and schema alterations to existing records.

## Files Overview

### Data Conversions
- **004_update_price_to_rupiah.sql** - Price format conversion
  - Converts price storage from cents to rupiah
  - Updates: `menu_items.price` and `individual_orders.total_price`
  - Divides existing prices by 100 for proper rupiah representation

### Schema Modifications
- **006_remove_global_stock_empty.sql** - Stock management cleanup
  - Removes global stock empty functionality
  - Cleans up related columns and constraints

## Usage

### Price Conversion
**⚠️ IMPORTANT**: This is a one-time conversion script. Only run if prices are stored in cents format:
```bash
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f updates/001_update_price_to_rupiah.sql
```

### Stock Management Updates
```bash
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f updates/002_remove_global_stock_empty.sql
```

## Before Running Updates

### Backup Recommendation
Always backup your database before running update scripts:
```bash
pg_dump -h localhost -p 5432 -U lunch_user lunch_delivery > backup_$(date +%Y%m%d_%H%M%S).sql
```

### Check Current State
Verify current data state before updates:
```bash
# Check current price format
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -c "SELECT name, price FROM menu_items LIMIT 5;"

# Check for global stock empty columns (before removal)
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -c "\d+ menu_items"
```

## Notes
- **Always backup before running updates**
- These scripts modify existing data - use with extreme caution in production
- Test in development environment first
- Some updates are irreversible (like data format conversions)
- Document any custom updates made to production data
