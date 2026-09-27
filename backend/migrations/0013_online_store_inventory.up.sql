-- 1. Create store_categories table
CREATE TABLE IF NOT EXISTS store_categories (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Create store_products table
CREATE TABLE IF NOT EXISTS store_products (
    id TEXT PRIMARY KEY,
    category_id TEXT NOT NULL,
    name TEXT NOT NULL UNIQUE,
    price_cents INTEGER NOT NULL CHECK (price_cents >= 0),
    stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    description TEXT,
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(category_id) REFERENCES store_categories(id) ON DELETE RESTRICT
);

-- 3. Create store_orders table (for tracking purchases of general store inventory)
CREATE TABLE IF NOT EXISTS store_orders (
    id TEXT PRIMARY KEY,
    user_id TEXT, -- Nullable to support guest checkout
    pre_tax_cents INTEGER NOT NULL DEFAULT 0 CHECK (pre_tax_cents >= 0),
    tax_cents INTEGER NOT NULL DEFAULT 0 CHECK (tax_cents >= 0), -- 13% HST
    total_cents INTEGER NOT NULL DEFAULT 0 CHECK (total_cents >= 0),
    payment_status TEXT NOT NULL DEFAULT 'pending' CHECK (payment_status IN ('pending', 'paid', 'refunded', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 4. Create store_order_items table
CREATE TABLE IF NOT EXISTS store_order_items (
    id TEXT PRIMARY KEY,
    order_id TEXT NOT NULL,
    product_id TEXT NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price_paid_cents INTEGER NOT NULL CHECK (price_paid_cents >= 0),
    FOREIGN KEY(order_id) REFERENCES store_orders(id) ON DELETE CASCADE,
    FOREIGN KEY(product_id) REFERENCES store_products(id) ON DELETE RESTRICT
);

-- 5. Create indexes to speed up joins
CREATE INDEX IF NOT EXISTS idx_store_products_category ON store_products(category_id);
CREATE INDEX IF NOT EXISTS idx_store_orders_user ON store_orders(user_id);
CREATE INDEX IF NOT EXISTS idx_store_order_items_order ON store_order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_store_order_items_product ON store_order_items(product_id);

-- 6. Seed store categories
INSERT INTO store_categories (id, name, display_order) VALUES
    ('cat-uniform', 'Uniform', 1),
    ('cat-clothing', 'Clothing', 2),
    ('cat-equipment', 'Equipment', 3),
    ('cat-accessories', 'Accessories', 4),
    ('cat-miscellaneous', 'Miscellaneous', 5);

-- 7. Seed products catalog
INSERT INTO store_products (id, category_id, name, price_cents, stock_quantity, description) VALUES
    -- Uniform
    ('prod-hightop-shoes', 'cat-uniform', 'Feiyue Authentic Shaolin Hightop Training Shoes for Men and Women', 3500, 50, 'Authentic Shaolin high-top training shoes with durable canvas and flexible rubber soles.'),
    ('prod-bamboo-tshirt', 'cat-uniform', 'Shaolin Yellow T-Shirt Bamboo Student Uniform', 3000, 100, 'Soft, breathable bamboo yellow student uniform t-shirt with embroidered Shaolin logo.'),
    ('prod-drifit-tshirt', 'cat-uniform', 'Shaolin T-Shirt Uniform Dri-fit Yellow', 2500, 75, 'High-performance dri-fit student uniform t-shirt ideal for intensive training sessions.'),
    ('prod-lowcut-shoes', 'cat-uniform', 'Feiyue Authentic Shaolin Training Shoes for Men and Women', 2500, 120, 'Authentic low-cut training shoes, lightweight and perfect for kung fu stance work.'),
    ('prod-std-uniform-set', 'cat-uniform', 'Standard Uniform Set for Training and Classes', 5000, 60, 'Complete uniform set including a cotton yellow t-shirt, training pants, and low-cut shoes.'),
    ('prod-dlx-uniform-set', 'cat-uniform', 'Deluxe Uniform Set for Training and Classes', 8000, 30, 'Premium uniform set featuring a bamboo yellow shirt, premium pants, and high-top shoes.'),
    ('prod-bamboo-pants', 'cat-uniform', 'Shaolin Pants Uniform Bamboo', 4000, 80, 'Ultra-comfortable kung fu training pants made from premium organic bamboo fibers.'),
    ('prod-term-locker', 'cat-uniform', 'Term Locker Rental', 2500, 15, 'Locker rental for one term. Safely store your uniform, shoes, and weapons at the school.'),
    ('prod-long-sleeve-shirt', 'cat-uniform', 'Shaolin Long Sleeve Yellow Student Shirt', 3500, 45, 'Long sleeve version of the Shaolin yellow student uniform, excellent for cooler terms.'),
    ('prod-robe-uniform', 'cat-uniform', 'Shaolin Robe Uniforms', 9000, 20, 'Official Shaolin lay disciple robe for special ceremonies and traditional forms training.'),
    ('prod-grey-robe-set', 'cat-uniform', 'Shaolin Grey Robe and Pants Set', 11000, 15, 'Traditional grey monk-style robe and trousers set imported directly from Henan.'),

    -- Clothing
    ('prod-bamboo-col', 'cat-clothing', 'Shaolin Bamboo T-Shirt Assorted Colours', 2500, 50, 'Casual everyday bamboo t-shirt available in slate, blue, grey, and black.'),
    ('prod-hoodie-unisex', 'cat-clothing', 'Official Shaolin Team Canada Unisex Hoodie', 6000, 30, 'Warm fleece unisex hoodie featuring the official Shaolin Team Canada insignia.'),
    ('prod-warrior-tshirt', 'cat-clothing', 'Shaolin Warrior Monk Bamboo T-Shirt Unisex', 3000, 60, 'Unisex bamboo shirt with a graphic design of a Shaolin warrior monk in stance.'),
    ('prod-med-tshirt', 'cat-clothing', 'Shaolin Meditating Monk Bamboo T-Shirt Unisex', 3000, 60, 'Breathable bamboo shirt depicting a monk in deep meditation, perfect for Tai Chi practitioners.'),
    ('prod-med-tank-men', 'cat-clothing', 'Shaolin Meditating Monk Bamboo Tank Top for Men', 2500, 40, 'Mens lightweight athletic tank top with meditating monk graphics.'),
    ('prod-warrior-tank-men', 'cat-clothing', 'Shaolin Warrior Monk Bamboo Tank Top for Men', 2500, 40, 'Mens athletic tank top featuring dynamic warrior monk print.'),
    ('prod-warrior-tank-women', 'cat-clothing', 'Shaolin Warrior Monk Bamboo Tank Top for Women', 2500, 40, 'Womens tailored bamboo tank top with martial arts graphics.'),
    ('prod-med-tank-women', 'cat-clothing', 'Shaolin Meditation Monk Bamboo Tank Top for Women', 2500, 40, 'Womens tailored tank top highlighting the zen meditation graphic.'),
    ('prod-hoodie-women', 'cat-clothing', 'Official Shaolin Team Canada Womens Hoodie', 6000, 25, 'Tailored women''s fit team hoodie with embroidered national emblems.'),
    ('prod-monk-socks', 'cat-clothing', 'Monk Socks', 1500, 100, 'Traditional orange wrap-style socks worn by Shaolin practitioners during practice.'),
    ('prod-monk-vest', 'cat-clothing', 'Martial Arts Vest', 4500, 35, 'Reversible training vest displaying traditional martial calligraphy and trim.'),
    ('prod-perf-belt', 'cat-clothing', 'Performance Belt', 2000, 150, 'Traditional satin sash belt used for performance, available in red, yellow, and black.'),
    ('prod-anniv-tshirt', 'cat-clothing', 'Academy Limited Edition T-Shirt', 3000, 40, 'Collector''s edition t-shirt celebrating traditional martial arts.'),
    ('prod-dragon-tshirt', 'cat-clothing', 'Dragon and Phoenix Bamboo-Cotton T-Shirt in White or Black', 3500, 55, 'Premium blended graphic shirt featuring the iconic dragon and phoenix motif.'),
    ('prod-jacket-unisex', 'cat-clothing', 'Academy Hoodie Jacket Unisex', 7000, 20, 'Premium full-zip jacket hoodie with heavy lining for outdoor training.'),
    ('prod-perf-costume', 'cat-clothing', 'Performance Costume', 9000, 15, 'Ornate silk-satin uniform set designed for stage performances and demonstrations.'),

    -- Equipment
    ('prod-punch-pad', 'cat-equipment', 'Heavy Duty Wall Mounted Punching Training Strike Pad', 5000, 10, 'Wall mountable high-density foam strike pad for hand and fist hardening.'),
    ('prod-cudgel-staff', 'cat-equipment', 'Waxwood Bo Staff Cudgel', 2500, 80, 'Flexible white wax wood staff, the fundamental weapon of traditional martial arts.'),
    ('prod-incense-med', 'cat-equipment', 'Incense for Meditation', 1200, 150, 'Natural sandalwood incense sticks to create a serene environment for meditation.'),
    ('prod-pillow-zafu', 'cat-equipment', 'Meditation Pillow Zafu with Buckwheat Filling', 4500, 25, 'Traditional round meditation cushion filled with organic buckwheat hulls.'),
    ('prod-sanda-gear', 'cat-equipment', 'Sanda Tournament Gear Set with Martial Arts Duffle Bag', 18000, 8, 'Full sparring kit containing gloves, headgear, shin guards, chest protector, and bag.'),
    ('prod-broadsword', 'cat-equipment', 'Traditional Broadsword', 9000, 12, 'Semi-flexible steel broadsword (Dao) for form practice and weapon instruction.'),
    ('prod-blocker-pad', 'cat-equipment', 'Blocker Strike Pad', 3500, 18, 'Handheld focus target pad for kicking and punching drills.'),
    ('prod-incense-holder', 'cat-equipment', 'Incense Holder', 1500, 40, 'Handcrafted wooden incense burner tray with brass accents.'),
    ('prod-sword-bag', 'cat-equipment', 'Sword Protective Bag', 2000, 50, 'Canvas carrying bag with adjustable strap to protect your broadsword.'),
    ('prod-staff-bag', 'cat-equipment', 'Staff Protective Bag', 2000, 50, 'Long cylindrical zip bag designed to carry and store your white wax wood staff.'),
    ('prod-massage-balls', 'cat-equipment', 'Massage Ball Set', 1500, 60, 'Double massage mobility lacrosse balls for trigger point release and muscle recovery.'),

    -- Accessories
    ('prod-shaker-bottle', 'cat-accessories', 'Academy Protein Shaker Water Bottle', 1500, 90, 'Leak-proof shaker bottle with wire whisk ball, branded with the school virtues.'),
    ('prod-virtue-mug', 'cat-accessories', 'Martial Virtues Mug Assorted', 1200, 110, 'Ceramic mug printed with one of the core martial virtues (Respect, Focus, Perseverance).'),
    ('prod-sling-bag', 'cat-accessories', 'Martial Arts Sling Bag', 2500, 40, 'Compact shoulder sling pack featuring traditional graphic prints.'),
    ('prod-duffle-bag', 'cat-accessories', 'Martial Arts Convertible Duffle Bag', 5000, 20, 'Spacious gym duffle bag that converts into a backpack with hidden straps.'),
    ('prod-magic-mug', 'cat-accessories', 'Martial Virtues Magic Mug', 1800, 35, 'Heat-activated color changing mug that reveals martial arts calligraphy.'),
    ('prod-qigong-poster', 'cat-accessories', 'Ba Duan Jin Qigong Poster', 1000, 120, 'Instructional wall chart detailing the 8 pieces of brocade Qigong movements.'),
    ('prod-mouse-pad', 'cat-accessories', 'Dojo Mouse Pad', 1200, 80, 'Non-slip desk mat featuring the Dojo layout plan and training schedule.'),
    ('prod-virtue-bookmark', 'cat-accessories', 'Martial Virtues Bookmark', 500, 300, 'Laminated paper bookmark featuring the 12 martial virtues.'),
    ('prod-shield-umbrella', 'cat-accessories', 'Martial Arts Shield Umbrella', 3500, 25, 'Heavy-duty windproof umbrella with an ergonomic handle grip.'),
    ('prod-thymox-spray', 'cat-accessories', 'Thymox Ext Disinfectant Anti-Bacterial and Anti-Viral Spray', 1500, 40, 'Non-toxic, botanical anti-viral spray to sanitize equipment and footwear.'),
    ('prod-cap-hat', 'cat-accessories', 'Academy Cap Hat', 2000, 60, 'Snapback style baseball cap with the Academy emblem embroidered in gold.'),
    ('prod-thermos-bottle', 'cat-accessories', 'Bamboo Thermos Tea Water Bottle', 3000, 50, 'Insulated stainless steel water bottle wrapped in natural bamboo with tea infuser basket.'),
    ('prod-beanie-hat', 'cat-accessories', 'Beanie Winter Hat', 2000, 70, 'Warm acrylic knit beanie hat with Academy logo tag, perfect for winter commutes.'),
    ('prod-virtue-towel', 'cat-accessories', 'Martial Virtues Sweat Towel', 1500, 85, 'Microfiber quick-dry sports towel featuring key martial virtues.'),
    ('prod-weapon-fig', 'cat-accessories', 'Martial Arts Weapon Figurines', 2500, 30, 'Set of miniature figurines holding different classical martial weapons.'),
    ('prod-hand-fan', 'cat-accessories', 'Martial Arts Hand Fan', 1200, 100, 'Bamboo folding fan used for cooling off or martial fan form practice.'),
    ('prod-stances-fig', 'cat-accessories', 'Martial Arts Stances Figurines', 2500, 30, 'Set of mini figurines showing correct postures for Horse, Bow, and Flat stances.'),
    ('prod-metallic-bookmark', 'cat-accessories', 'Martial Virtues Metallic Bookmark', 800, 200, 'Premium laser-engraved metal bookmark with red tassel.'),
    ('prod-bamboo-pen', 'cat-accessories', 'Academy Bamboo Pen', 500, 500, 'Eco-friendly retractable ballpoint pen crafted from sustainable bamboo wood.'),
    ('prod-metallic-keychain', 'cat-accessories', 'Academy Metallic Keychain', 800, 150, 'Polished zinc-alloy keyring stamped with the Academy emblem.'),
    ('prod-sweatband', 'cat-accessories', 'Academy Sweatband', 1000, 120, 'Elastic cotton wrist and forehead sweatbands pack in black or red.'),

    -- Miscellaneous
    ('prod-instructor-shirt', 'cat-miscellaneous', 'Instructor Shirt 50/50 Bamboo-Cotton', 3500, 20, 'Official grey instructor t-shirt, restricted to certified teaching staff.'),
    ('prod-gift-cert', 'cat-miscellaneous', 'Gift Certificate', 5000, 999, 'Digital or physical gift voucher worth $50, redeemable for classes or gear.');
