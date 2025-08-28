1. to create tables, run:
psql -h localhost -U lunch_user -d lunch_delivery -f migrations/001_initial.sql
or
sudo -i -u postgres psql -d lunch_delivery -f /path/to/your/project/migrations/001_initial.sql

2. to populate menu_items, run:
INSERT INTO menu_items (name, price) VALUES
('Cah jagung muda', 1500),
('Cah labu', 1500),
('Cah toge', 1200),
('Cah kembang kol', 1500),
('Cah oyong telur', 1800),
('Ceker cabe ijo', 2000),
('Ayam goreng kandar merah', 2500),
('Udang crispy cabe garam', 3000);

3. to populate companies, run:
INSERT INTO companies (name, address, contact) VALUES
('Tuntun Sekuritas', 'ASG Tower Lt. 11', '082317384627'),
('Indomaret', 'ASG Tower Lt. 9', '08231748353'),
('Lo Accounting Ltd', 'PIK Avenue Lt. 2', '08132649297');
