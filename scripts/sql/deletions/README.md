# Deletions Scripts

This directory contains SQL scripts for data cleanup, removal, and testing scenarios that require data deletion.

## Files Overview

### Test Data Cleanup
- **delete_test_order.sql** - Removes Jemmy's order for testing
  - Deletes today's order for employee_id: 1 (Jemmy)
  - Used to test footer edge cases when user has no orders
  - Includes verification query to confirm deletion

### Session Management
- **007_delete_today_order_session.sql** - Removes today's order session
  - Cleans up current day's order session for testing
  - Useful for resetting daily order states

### UI Testing
- **009_test_footer_edge_case.sql** - Footer edge case testing
  - Sets up specific conditions for testing footer behavior
  - Creates scenarios where users have no orders to test UI edge cases

## Usage

### Test Environment Cleanup
Remove specific test orders:
```bash
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f deletions/delete_test_order.sql
```

### Session Reset
Clear today's order session:
```bash
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f deletions/007_delete_today_order_session.sql
```

### UI Edge Case Testing
Set up edge case scenarios:
```bash
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -f deletions/009_test_footer_edge_case.sql
```

## Safety Guidelines

### ⚠️ Production Warning
**NEVER run deletion scripts in production without explicit approval and backup!**

### Pre-Deletion Checklist
1. **Backup First**: Always create a backup before deleting data
2. **Verify Target**: Double-check which records will be affected
3. **Test Environment**: Run in development first
4. **Document**: Record what was deleted and why

### Backup Before Deletion
```bash
# Full database backup
pg_dump -h localhost -p 5432 -U lunch_user lunch_delivery > backup_before_deletion_$(date +%Y%m%d_%H%M%S).sql

# Backup specific tables
pg_dump -h localhost -p 5432 -U lunch_user lunch_delivery -t individual_orders -t order_sessions > orders_backup.sql
```

### Verification Queries
Check what will be deleted before running scripts:
```bash
# Check Jemmy's orders
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -c "
SELECT * FROM individual_orders WHERE employee_id = 1
AND session_id = (SELECT id FROM order_sessions WHERE date = CURRENT_DATE AND company_id = 1);"

# Check today's sessions
PGPASSWORD=1234 psql -h localhost -p 5432 -U lunch_user -d lunch_delivery -c "
SELECT * FROM order_sessions WHERE date = CURRENT_DATE;"
```

## Notes
- These scripts are primarily for testing and development
- All deletion scripts include verification queries
- Use extreme caution - deletions are typically irreversible
- Consider using transactions for complex deletions
- Document any production deletions in change logs