to create tables, run:
psql -h localhost -U lunch_user -d lunch_delivery -f migrations/001_initial.sql
or
sudo -i -u postgres psql -d lunch_delivery -f /path/to/your/project/migrations/001_initial.sql
