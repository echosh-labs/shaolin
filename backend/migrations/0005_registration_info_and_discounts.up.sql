-- 1. Create registration_faqs table for new student onboarding content
CREATE TABLE IF NOT EXISTS registration_faqs (
    id TEXT PRIMARY KEY,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    category TEXT NOT NULL DEFAULT 'general' CHECK (category IN ('general', 'discounts', 'uniforms', 'terms', 'ping_pong')),
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Create families / households table to manage family relationships
CREATE TABLE IF NOT EXISTS families (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL, -- e.g. 'Smith Household'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Create family_members table mapping users to families
CREATE TABLE IF NOT EXISTS family_members (
    family_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (family_id, user_id),
    FOREIGN KEY(family_id) REFERENCES families(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 4. Create user_memberships table to track annual charitable donations/receipts
CREATE TABLE IF NOT EXISTS user_memberships (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    calendar_year INTEGER NOT NULL,
    amount_cents INTEGER NOT NULL CHECK (amount_cents >= 0),
    payment_date TEXT NOT NULL CHECK (payment_date LIKE '____-__-__'),
    receipt_issued INTEGER NOT NULL DEFAULT 0 CHECK (receipt_issued IN (0, 1)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (user_id, calendar_year)
);

-- 5. Create discount_rules configuration table
CREATE TABLE IF NOT EXISTS discount_rules (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    discount_type TEXT NOT NULL CHECK (discount_type IN ('family', 'returning')),
    trigger_value INTEGER NOT NULL, -- family size or completed terms threshold
    discount_percentage INTEGER NOT NULL CHECK (discount_percentage BETWEEN 0 AND 100),
    description TEXT
);

-- 6. Seed Registration FAQs
INSERT INTO registration_faqs (id, question, answer, category, display_order) VALUES
    ('faq-1', 'Where can I find school information?', 'Click here for information about our school, classes, rules', 'general', 1),
    ('faq-2', 'How Do I Add Family Members?', 'To add family members, click on ''Account > Family Manager''. You will find the option to add and manage your family members there.', 'general', 2),
    ('faq-3', 'How Do School Terms Work?', 'There are 3 terms each year. Each term is 4 months and there are 14 weeks of actual classes each term. There is also usually a 1 week break around the middle of the term, as well as a 2 week break at the end of each term.', 'terms', 3),
    ('faq-4', 'How Does Ping Pong Work?', 'Join our fun and challenging Ping Pong Club! Already included with Unlimited Shaolin Tokens. You can use regular Shaolin Tokens to play ping pong as well. More information can be found here.', 'ping_pong', 4),
    ('faq-5', 'How do I choose # of tokens?', 'Ideally you should plan to purchase the number of tokens based on your expected hours of attendance each week. For example, if you plan to attend 3 classes each week, you will want 14 x 3 = 42 tokens.', 'terms', 5),
    ('faq-6', 'Can I carry forward tokens?', 'Tokens must be used within the term and cannot be carried forward. They can be transferred to another family member for a small exchange fee (Account > Family Manager > Carry Forward).', 'terms', 6),
    ('faq-7', 'Do I need to purchase a uniform?', 'A uniform is required. The uniform set comes with t-shirt, pants, and Feiyue shoes. You can upgrade to high-top Feiyue shoes if you prefer more ankle support. Contact us if you require an exception.', 'uniforms', 7),
    ('faq-8', 'Is a uniform needed for trial class?', 'No uniform is required for trial classes. Please wear something comfortable such as jogging pants and a t-shirt, with clean, indoor, non-marking shoes. Bare feet is not allowed for hygiene purposes.', 'uniforms', 8),
    ('faq-9', 'Why is membership required?', 'The membership is a charitable donation to Shaolin Luohan Temple. You can obtain a receipt for tax reduction purposes. Part of the funds are used for benevolent projects and social activities.', 'general', 9),
    ('faq-10', 'Are there children or senior discounts?', 'Yes, children and seniors get a discounted rate and you will see cheaper rate options for tokens. This is enabled after you have submitted proof of age with a birth certificate, driver''s license, or passport.', 'discounts', 10),
    ('faq-11', 'Are there family discounts?', 'Yes! With 2 family members there is a 3% discount for EACH member. 3 members get 4% each, and 4+ members get 5% each. Must live in the same household and purchase at least 14 tokens.', 'discounts', 11),
    ('faq-12', 'Are there Returning Student discounts?', 'Yes! Students who register BEFORE the term starts will receive a discount based on # Completed Terms (i.e. at least 14 tokens). 1-3 Terms = 5%, 4-6 Terms = 10%, 7-9 Terms = 15%, 10+ Terms = 20%', 'discounts', 12);

-- 7. Seed Discount Rules Catalog
INSERT INTO discount_rules (id, name, discount_type, trigger_value, discount_percentage, description) VALUES
    ('disc-fam-2', 'Family Discount (2 members)', 'family', 2, 3, '3% discount for each member with 2 household members purchasing >= 14 tokens'),
    ('disc-fam-3', 'Family Discount (3 members)', 'family', 3, 4, '4% discount for each member with 3 household members purchasing >= 14 tokens'),
    ('disc-fam-4', 'Family Discount (4+ members)', 'family', 4, 5, '5% discount for each member with 4+ household members purchasing >= 14 tokens'),
    ('disc-ret-1', 'Returning Student (1-3 terms)', 'returning', 1, 5, '5% discount for 1-3 completed terms'),
    ('disc-ret-4', 'Returning Student (4-6 terms)', 'returning', 4, 10, '10% discount for 4-6 completed terms'),
    ('disc-ret-7', 'Returning Student (7-9 terms)', 'returning', 7, 15, '15% discount for 7-9 completed terms'),
    ('disc-ret-10', 'Returning Student (10+ terms)', 'returning', 10, 20, '20% discount for 10+ completed terms');

-- 8. Create indexes to speed up discount calculation joins
CREATE INDEX IF NOT EXISTS idx_family_members_user_id ON family_members(user_id);
CREATE INDEX IF NOT EXISTS idx_user_memberships_user_id ON user_memberships(user_id);
CREATE INDEX IF NOT EXISTS idx_discount_rules_type ON discount_rules(discount_type);
