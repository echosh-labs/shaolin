-- 1. Disable foreign keys temporarily
PRAGMA foreign_keys = OFF;

-- 2. Drop student_term_registrations and waivers table to recreate
DROP TABLE IF EXISTS student_term_registrations;
DROP TABLE IF EXISTS waivers;

-- 3. Create waivers table
CREATE TABLE waivers (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Seed the official Martial Arts Academy Class Registration Waiver
INSERT INTO waivers (id, title, content, is_active) VALUES (
    'waiver-class-registration',
    'Waiver, Warranty and Consent',
    'Learning and practicing martial arts, kung fu, karate, kobudo, tai chi, and qigong is a rewarding experience but it comes with risks of injury. For example, if your training area is small, you are more likely to have an accident with a piece of furniture or with someone else! The drills and exercises we provide are safe but one can get injured if you push yourself too much or have an underlying health condition. Please be sure to read our Virtual Online Class Guide on our website to learn how to set up your training environment. Use common sense for your own safety and have fun getting stronger!

Waiver, Warranty and Consent
PLEASE READ BEFORE SIGNING!
IN CONSIDERATION of allowing me to participate in the program, related events and activities of Martial Arts Academy (also known as the Academy), in person at one of our school locations, in a Virtual Class, and through our online digital classes:

I WARRANT TO YOU THAT:

I am familiar with the risk of serious injury and death which any participant in this programme must assume, and
I believe that I am physically, emotionally and mentally able to participate in this programme, and that my equipment is mechanically fit for my use in this programme, and
I understand that all applicable rules for participation must be followed and that at all times the sole responsibility for personal safety remains with me, and
I will immediately remove myself from participation, if at any time I sense or observe any unusual hazard or unsafe condition or if I feel that I have experienced any deterioration in my physical, emotional or mental fitness for continued participation in the programme.
I UNDERSTAND AND AGREE, on behalf of myself, my heirs, assigns, personal representatives and next of kin, that my participation in this programme and execution of this document constitutes:

An unqualified ASSUMPTION OF ALL RISKS associated with participation in this programme by me even if arising from negligent, or gross negligence, including any compounding or aggravation of injuries caused by negligent rescue operations or procedures, of the programme organizer and any person associated therewith or participating therein, and
A FULL AND FINAL RELEASE AND WAIVER OF LIABILITY of the programme organizer and all persons and organizations associated with it and the programme including, without limiting the generality of the foregoing, its officers, directors, agents and/or employees, other participants, sponsors, advertisers, owners and/or lessors of the premises used to conduct the programme, sanctioning bodies, medical or rescue personnel (the RELEASEES), of and from with the respect to all injury, disability, death or loss or damage to person or property whether arising from the negligence, or negligent rescue of or by the foregoing or otherwise, and
An UNDERSTANDING NOT TO SUE the RELEASEES for any loss, injury, costs or damages of any form or type, howsoever caused or arising, and whether directly or indirectly from the participation in this programme by me, and
An AGREEMENT TO INDEMNIFY, and to SAVE and HOLD HARMLESS the RELEASEES, and each of them, from any litigation expense, legal fees, liability, damage award or cost, of any form or type whatsoever, they may incur due to any claim made against them or any one of them whether the claim is based on the negligence of the RELEASEES or otherwise.
I HAVE READ THIS DOCUMENT THOROUGHLY. I UNDERSTAND THAT THE RELEASEES ARE RELYING UPON MY WARRANTIES, ASSUMPTION, WAIVER AND RELEASE, UNDERTAKINGS AND AGREEMENTS WHEN ACCEPTING MY PARTICIPATING IN THIS PROGRAMME. I UNDERSTAND THAT BY SIGNING THIS DOCUMENT I GIVE UP SUBSTANTIAL LEGAL RIGHTS I WOULD OTHERWISE HAVE.

If you are under 18 years of age you must have your parent or guardian sign this waiver.',
    1
);

-- 5. Recreate student_term_registrations table with mandatory waiver signature fields
CREATE TABLE student_term_registrations (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    term_id TEXT NOT NULL,
    token_package_id TEXT REFERENCES token_packages(id) ON DELETE SET NULL, -- Nullable to allow "None" tokens
    membership_option_id TEXT REFERENCES membership_options(id) ON DELETE RESTRICT, -- References dynamic membership plans
    
    -- Uniform options
    uniform_package_id TEXT REFERENCES uniform_packages(id) ON DELETE SET NULL,
    uniform_ordered INTEGER NOT NULL DEFAULT 0 CHECK (uniform_ordered IN (0, 1)),
    uniform_size TEXT CHECK (uniform_size IS NULL OR uniform_size IN ('XS', 'S', 'M', 'L', 'XL', 'XXL')),
    shoe_size INTEGER CHECK (shoe_size IS NULL OR (shoe_size BETWEEN 30 AND 50)),
    
    -- Ping Pong club option
    ping_pong_package_id TEXT REFERENCES ping_pong_packages(id) ON DELETE SET NULL,
    ping_pong_club_joined INTEGER NOT NULL DEFAULT 0 CHECK (ping_pong_club_joined IN (0, 1)),
    
    -- Mandatory Waiver Signature fields
    waiver_id TEXT NOT NULL REFERENCES waivers(id) ON DELETE RESTRICT,
    waiver_signer_name TEXT NOT NULL,
    waiver_signed_date TEXT NOT NULL CHECK (waiver_signed_date LIKE '____-__-__'),
    waiver_signed_ip TEXT NOT NULL,
    
    -- Discount percentages applied
    family_discount_applied_percent INTEGER NOT NULL DEFAULT 0 CHECK (family_discount_applied_percent BETWEEN 0 AND 100),
    returning_discount_applied_percent INTEGER NOT NULL DEFAULT 0 CHECK (returning_discount_applied_percent BETWEEN 0 AND 100),
    
    -- Detailed pre-tax, tax, and non-taxable fees ledger
    class_tokens_fee_cents INTEGER NOT NULL DEFAULT 0 CHECK (class_tokens_fee_cents >= 0),
    ping_pong_fee_cents INTEGER NOT NULL DEFAULT 0 CHECK (ping_pong_fee_cents >= 0),
    uniform_fee_cents INTEGER NOT NULL DEFAULT 0 CHECK (uniform_fee_cents >= 0),
    tax_cents INTEGER NOT NULL DEFAULT 0 CHECK (tax_cents >= 0), -- 13% HST on classes and products
    membership_fee_cents INTEGER NOT NULL DEFAULT 0 CHECK (membership_fee_cents >= 0), -- Saved price paid for audit
    total_fee_cents INTEGER NOT NULL DEFAULT 0 CHECK (total_fee_cents >= 0), -- Grand total
    
    payment_status TEXT NOT NULL DEFAULT 'pending' CHECK (payment_status IN ('pending', 'paid', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(term_id) REFERENCES terms(id) ON DELETE CASCADE,
    UNIQUE (user_id, term_id),
    CHECK ((ping_pong_package_id IS NULL AND ping_pong_club_joined = 0) OR (ping_pong_package_id IS NOT NULL AND ping_pong_club_joined = 1)),
    CHECK ((uniform_package_id IS NULL AND uniform_ordered = 0) OR (uniform_package_id IS NOT NULL AND uniform_ordered = 1)),
    CHECK ((uniform_ordered = 0 AND uniform_size IS NULL AND shoe_size IS NULL) OR (uniform_ordered = 1 AND uniform_size IS NOT NULL AND shoe_size IS NOT NULL))
);

-- 6. Create registrations indexes
CREATE INDEX IF NOT EXISTS idx_student_term_registrations_user ON student_term_registrations(user_id);
CREATE INDEX IF NOT EXISTS idx_student_term_registrations_term ON student_term_registrations(term_id);

-- 7. Re-enable foreign keys
PRAGMA foreign_keys = ON;
