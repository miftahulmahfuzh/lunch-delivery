-- Migration: Seed menu items
-- This script adds all the standard menu items for the lunch delivery system

INSERT INTO menu_items (name, price) VALUES
-- Vegetables
('Cah buncis tempe', 400000),
('Cah pare rebon', 420000),
('Cah labu', 400000),
('Cah kembang kol', 410000),
('Cah jagung muda', 420000),
('Terong cabe ijo', 430000),
('Cah oyong telur', 500000),
('Cah toge', 450000),

-- Main dishes - Meat/Protein
('Kikil balado', 500000),
('Oncom leuncah', 450000),
('Orek basah', 450000),
('Kentang balado', 450000),
('Kentang semur', 480000),
('Baby cumi cabe ijo', 550000),
('Dori cabe ijo', 580000),
('Fillet ayam oseng pedes daun jeruk', 550000),
('Fillet ayam kungpao', 580000),
('Udang gede crispy cabe garlic', 600000),
('Udang crispy cabe garam', 500000),
('Oseng sapi setan', 580000),
('Ikan kembung dabu dabu', 520000),
('Telur ceplok cabe ijo', 450000),
('Fuyunghai asam manis', 480000),
('Jamur crispyyyyy', 460000),
('Rolade ayam', 520000),
('Perkedel', 420000),
('Tempe lada hitammss', 450000),
('Tahu wijennss', 450000),
('Tongkol sarden', 500000),
('Ikan cue sarden', 500000),
('Ceker oseng pedes', 480000),
('Ceker cabe ijo', 350000),
('Ayam suwir', 500000),
('Ayam goreng kandar merah', 450000),

-- Noodles/Rice
('Bihun goreng', 470000),
('Mie goreng', 480000),
('Nasi', 500000),
('Nasi 1/2', 300000),

-- Snacks
('Gorengan - bakwan sayur 🥰', 400000),

-- Desserts
('Puding lumut ijo 🥰', 420000),

-- Drinks
('Es mambo 🥰', 400000),
('Es ketan itam', 420000),
('Es sirup merah', 400000),
('Es coklat', 410000),
('Es Thai tea', 450000)

ON CONFLICT (name) DO NOTHING;