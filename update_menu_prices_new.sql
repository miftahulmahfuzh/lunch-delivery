-- Update menu items with new price list

-- Update existing items that match (case-insensitive matching)
UPDATE menu_items SET price = 4000 WHERE LOWER(name) LIKE '%nasi%' AND LOWER(name) NOT LIKE '%1/2%' AND LOWER(name) NOT LIKE '%half%';
UPDATE menu_items SET price = 2000 WHERE LOWER(name) LIKE '%nasi%' AND (LOWER(name) LIKE '%1/2%' OR LOWER(name) LIKE '%half%');

-- Main Dishes (Lauk Protein)
UPDATE menu_items SET price = 20000 WHERE LOWER(name) LIKE '%oseng sapi setan%';
UPDATE menu_items SET price = 15000 WHERE LOWER(name) LIKE '%filet kungpao%' OR LOWER(name) LIKE '%fillet%kungpao%';
UPDATE menu_items SET price = 14000 WHERE LOWER(name) LIKE '%ikan kembung dabu dabu%';
UPDATE menu_items SET price = 12000 WHERE LOWER(name) LIKE '%ayam oseng pedes daun jeruk%' OR LOWER(name) LIKE '%fillet ayam oseng pedes daun jeruk%';
UPDATE menu_items SET price = 11000 WHERE LOWER(name) LIKE '%fuyunghai asam manis%';
UPDATE menu_items SET price = 8000 WHERE LOWER(name) LIKE '%baby cumi cabe ijo%';
UPDATE menu_items SET price = 8000 WHERE LOWER(name) LIKE '%rolade ayam%';
UPDATE menu_items SET price = 7000 WHERE LOWER(name) LIKE '%tempe lada hitam%';
UPDATE menu_items SET price = 5000 WHERE LOWER(name) LIKE '%dori cabe ijo%';
UPDATE menu_items SET price = 4000 WHERE LOWER(name) LIKE '%ayam suwir%';

-- Side Dishes (Sayur & Pelengkap)
UPDATE menu_items SET price = 6000 WHERE LOWER(name) LIKE '%perkedel%';
UPDATE menu_items SET price = 6000 WHERE LOWER(name) LIKE '%oncom leunca%';
UPDATE menu_items SET price = 5000 WHERE LOWER(name) LIKE '%bakwan sayur%';
UPDATE menu_items SET price = 5000 WHERE LOWER(name) LIKE '%telor ceplok%' OR LOWER(name) LIKE '%telur ceplok%';
UPDATE menu_items SET price = 5000 WHERE LOWER(name) LIKE '%orek basah%';
UPDATE menu_items SET price = 5000 WHERE LOWER(name) LIKE '%cah buncis tempe%';
UPDATE menu_items SET price = 4000 WHERE LOWER(name) LIKE '%cah jagung muda%';
UPDATE menu_items SET price = 4000 WHERE LOWER(name) LIKE '%cah kembang kol%';
UPDATE menu_items SET price = 4000 WHERE LOWER(name) LIKE '%cah labu%';
UPDATE menu_items SET price = 4000 WHERE LOWER(name) LIKE '%jamur crispy%';
UPDATE menu_items SET price = 4000 WHERE LOWER(name) LIKE '%kentang semur%';
UPDATE menu_items SET price = 3000 WHERE LOWER(name) LIKE '%kentang balado%';

-- Insert new items that don't exist (check if they don't exist first)
INSERT INTO menu_items (name, price, active, created_at) 
SELECT * FROM (
    SELECT 'Nasi (1 Porsi)' as name, 4000 as price, true as active, NOW() as created_at
    UNION ALL
    SELECT 'Nasi (1/2 Porsi)', 2000, true, NOW()
    UNION ALL
    SELECT 'Oseng Sapi Setan', 20000, true, NOW()
    UNION ALL
    SELECT 'Filet Kungpao', 15000, true, NOW()
    UNION ALL
    SELECT 'Ikan Kembung Dabu Dabu', 14000, true, NOW()
    UNION ALL
    SELECT 'Ayam Oseng Pedes Daun Jeruk', 12000, true, NOW()
    UNION ALL
    SELECT 'Fuyunghai Asam Manis', 11000, true, NOW()
    UNION ALL
    SELECT 'Baby Cumi Cabe Ijo', 8000, true, NOW()
    UNION ALL
    SELECT 'Rolade Ayam', 8000, true, NOW()
    UNION ALL
    SELECT 'Tempe Lada Hitam', 7000, true, NOW()
    UNION ALL
    SELECT 'Dori Cabe Ijo', 5000, true, NOW()
    UNION ALL
    SELECT 'Ayam Suwir', 4000, true, NOW()
    UNION ALL
    SELECT 'Perkedel', 6000, true, NOW()
    UNION ALL
    SELECT 'Oncom Leunca', 6000, true, NOW()
    UNION ALL
    SELECT 'Bakwan Sayur', 5000, true, NOW()
    UNION ALL
    SELECT 'Telor Ceplok', 5000, true, NOW()
    UNION ALL
    SELECT 'Orek Basah', 5000, true, NOW()
    UNION ALL
    SELECT 'Cah Buncis Tempe', 5000, true, NOW()
    UNION ALL
    SELECT 'Cah Jagung Muda', 4000, true, NOW()
    UNION ALL
    SELECT 'Cah Kembang Kol', 4000, true, NOW()
    UNION ALL
    SELECT 'Cah Labu', 4000, true, NOW()
    UNION ALL
    SELECT 'Jamur Crispy', 4000, true, NOW()
    UNION ALL
    SELECT 'Kentang Semur', 4000, true, NOW()
    UNION ALL
    SELECT 'Kentang Balado', 3000, true, NOW()
) AS new_items
WHERE NOT EXISTS (
    SELECT 1 FROM menu_items 
    WHERE LOWER(menu_items.name) = LOWER(new_items.name)
);