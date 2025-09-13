-- Migration: Seed menu items
-- This script adds all the menu items currently in the database

INSERT INTO menu_items (name, price) VALUES
-- Main Dishes (Premium)
('Oseng Sapi Setan', 20000),
('Fillet Kungpao', 15000),
('Fillet Ayam Kungpao', 15000),
('Ikan Kembung Dabu Dabu', 14000),
('Ayam Oseng Pedes Daun Jeruk', 12000),
('Fillet Ayam Oseng Pedes Daun Jeruk', 12000),
('Fuyunghai Asam Manis', 11000),

-- Main Dishes (Standard)
('Baby Cumi Cabe Ijo', 8000),
('Rolade Ayam', 8000),
('Ayam Goreng Kandar Merah', 8000),
('Tempe Lada Hitam', 7000),
('Tempe lada hitammss', 7000),
('Ayam Suwir', 6000),
('Oncom Leunca', 6000),
('Oncom leuncah', 6000),
('Perkedel', 6000),
('Udang Gede Crispy Cabe Garlic', 6000),
('Test Item', 6000),

-- Side Dishes & Vegetables
('Bakwan Sayur', 5000),
('Cah Buncis Tempe', 5000),
('Telor Ceplok', 5000),
('Orek Basah', 5000),
('Dori Cabe Ijo', 5000),
('Cah Oyong Telur', 5000),
('Cah Pare Rebon', 5000),
('Ceker Cabe Ijo', 5000),
('Ceker Oseng Pedes', 5000),
('Ikan Cue Sarden', 5000),
('Kikil balado', 5000),
('Kikil Balado', 5000),
('Mie Goreng', 5000),
('Bihun Goreng', 5000),
('Tahu Wijen', 5000),
('Telur Ceplok Cabe Ijo', 5000),
('Terong Cabe Ijo', 5000),
('Tongkol Sarden', 5000),
('Udang crispy cabe garam', 5000),
('Cah Jagung Muda', 5000),

-- Vegetables (Basic)
('Cah Toge', 4500),
('Cah Jagung Muda', 4000),
('Cah kembang kol', 4000),
('Cah Kembang Kol', 4000),
('Cah labu', 4000),
('Cah Labu', 4000),
('Jamur Crispy', 4000),
('Jamur crispyyyyy', 4000),
('Kentang Semur', 4000),
('Kentang Balado', 4000),

-- Staples
('Nasi (1 Porsi)', 4000),
('Nasi', 4000),
('Nasi (1/2 Porsi)', 2000),
('Nasi 1/2', 2000),

-- Beverages & Desserts
('Es Coklat', 8000),
('Es Ketan Itam', 8000),
('Es Mambo', 8000),
('Es Sirup Merah', 8000),
('Es Thai Tea', 8000),
('Puding Lumut Ijo', 8000)

ON CONFLICT (name) DO NOTHING;