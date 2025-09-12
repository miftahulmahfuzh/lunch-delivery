-- Migration: Seed menu items
-- This script adds all the standard menu items for the lunch delivery system

INSERT INTO menu_items (name, price) VALUES
-- Staples (Nasi)
('Nasi (1 Porsi)', 4000),
('Nasi (1/2 Porsi)', 2000),

-- Main Dishes (Lauk Protein)
('Oseng Sapi Setan', 20000),
('Filet Kungpao', 15000),
('Ikan Kembung Dabu Dabu', 14000),
('Ayam Oseng Pedes Daun Jeruk', 12000),
('Fuyunghai Asam Manis', 11000),
('Baby Cumi Cabe Ijo', 8000),
('Rolade Ayam', 8000),
('Tempe Lada Hitam', 7000),
('Dori Cabe Ijo', 5000),
('Ayam Suwir', 4000),

-- Side Dishes (Sayur & Pelengkap)
('Perkedel', 6000),
('Oncom Leunca', 6000),
('Bakwan Sayur', 5000),
('Telor Ceplok', 5000),
('Orek Basah', 5000),
('Cah Buncis Tempe', 5000),
('Cah Jagung Muda', 4000),
('Cah Kembang Kol', 4000),
('Cah Labu', 4000),
('Jamur Crispy', 4000),
('Kentang Semur', 4000),
('Kentang Balado', 3000),

-- Legacy items (keeping for backward compatibility)
('Cah pare rebon', 4200),
('Terong cabe ijo', 4300),
('Cah oyong telur', 5000),
('Cah toge', 4500),
('Kikil balado', 5000),
('Udang gede crispy cabe garlic', 6000),
('Udang crispy cabe garam', 5000),
('Tahu wijennss', 4500),
('Tongkol sarden', 5000),
('Ikan cue sarden', 5000),
('Ceker oseng pedes', 4800),
('Ceker cabe ijo', 3500),
('Ayam goreng kandar merah', 4500),
('Bihun goreng', 4700),
('Mie goreng', 4800),
('Nasi', 4000),
('Nasi 1/2', 2000),
('Gorengan - bakwan sayur 🥰', 5000),
('Puding lumut ijo 🥰', 4200),
('Es mambo 🥰', 4000),
('Es ketan itam', 4200),
('Es sirup merah', 4000),
('Es coklat', 4100),
('Es Thai tea', 4500)

ON CONFLICT (name) DO NOTHING;