-- Migration: Seed menu items
-- This script adds all the standard menu items for the lunch delivery system

INSERT INTO menu_items (name, price) VALUES
-- Vegetables
('Cah buncis tempe', 4000),
('Cah pare rebon', 4200),
('Cah labu', 4000),
('Cah kembang kol', 4100),
('Cah jagung muda', 4200),
('Terong cabe ijo', 4300),
('Cah oyong telur', 5000),
('Cah toge', 4500),

-- Main dishes - Meat/Protein
('Kikil balado', 5000),
('Oncom leuncah', 4500),
('Orek basah', 4500),
('Kentang balado', 4500),
('Kentang semur', 4800),
('Baby cumi cabe ijo', 5500),
('Dori cabe ijo', 5800),
('Fillet ayam oseng pedes daun jeruk', 5500),
('Fillet ayam kungpao', 5800),
('Udang gede crispy cabe garlic', 6000),
('Udang crispy cabe garam', 5000),
('Oseng sapi setan', 5800),
('Ikan kembung dabu dabu', 5200),
('Telur ceplok cabe ijo', 4500),
('Fuyunghai asam manis', 4800),
('Jamur crispyyyyy', 4600),
('Rolade ayam', 5200),
('Perkedel', 4200),
('Tempe lada hitammss', 4500),
('Tahu wijennss', 4500),
('Tongkol sarden', 5000),
('Ikan cue sarden', 5000),
('Ceker oseng pedes', 4800),
('Ceker cabe ijo', 3500),
('Ayam suwir', 5000),
('Ayam goreng kandar merah', 4500),

-- Noodles/Rice
('Bihun goreng', 4700),
('Mie goreng', 4800),
('Nasi', 5000),
('Nasi 1/2', 3000),

-- Snacks
('Gorengan - bakwan sayur 🥰', 4000),

-- Desserts
('Puding lumut ijo 🥰', 4200),

-- Drinks
('Es mambo 🥰', 4000),
('Es ketan itam', 4200),
('Es sirup merah', 4000),
('Es coklat', 4100),
('Es Thai tea', 4500)

ON CONFLICT (name) DO NOTHING;