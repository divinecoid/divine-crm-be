-- =====================================================
-- DIVINE CRM - COMPLETE SETUP SCRIPT
-- =====================================================

-- Clean slate (OPTIONAL - comment out to preserve data)
-- TRUNCATE TABLE contacts, products, chat_labels, chat_messages, 
--          ai_configurations, connected_platforms, ai_agents, 
--          human_agents, broadcast_templates, quick_replies RESTART IDENTITY CASCADE;

-- =====================================================
-- 1. CONTACTS (Leads & Contacts)
-- =====================================================

INSERT INTO contacts (code, channel, channel_id, name, contact_status, temperature, first_contact, last_contact, last_agent, last_agent_type, created_at, updated_at)
VALUES
    ('M001', 'Telegram', '+628389956789', 'John Doe', 'Leads', 'Cold', '2025-10-23 09:28:27', '2025-10-25 08:28:27', 'Diva', 'AI', NOW(), NOW()),
    ('E003', 'WhatsApp', '+628389956789', 'Marcel', 'Contact', 'Warm', '2025-10-23 11:29:39', '2025-10-25 09:29:39', 'Niken', 'AI', NOW(), NOW()),
    ('W008', 'Instagram', '@wafxx', 'Wafi', 'Contact', 'Hot', '2025-10-23 15:14:14', '2025-10-25 14:14:14', 'Martono', 'Human', NOW(), NOW())
ON CONFLICT (code) DO UPDATE SET 
    name = EXCLUDED.name,
    temperature = EXCLUDED.temperature,
    updated_at = NOW();

-- =====================================================
-- 2. PRODUCTS
-- =====================================================

INSERT INTO products (code, name, price, stock, description, uploaded_by, created_at, updated_at)
VALUES
    ('F003', 'Produk 1', 28000, 853, 'Produk berkualitas tinggi untuk kebutuhan sehari-hari', 'Markonah', NOW(), NOW()),
    ('V005', 'Produk 2', 58000, 347, 'Produk premium dengan fitur lengkap dan modern', 'Markonak', NOW(), NOW()),
    ('O002', 'Produk 3', 352000, 495, 'Produk eksklusif dengan teknologi terkini', 'Martono', NOW(), NOW())
ON CONFLICT (code) DO UPDATE SET 
    name = EXCLUDED.name,
    price = EXCLUDED.price,
    stock = EXCLUDED.stock,
    updated_at = NOW();

-- =====================================================
-- 3. CHAT LABELS
-- =====================================================

INSERT INTO chat_labels (label, description, color, created_at, updated_at)
VALUES
    ('Nanya-nanya doang', 'Customer yang hanya bertanya tanpa niat beli', 'Red', NOW(), NOW()),
    ('Customer bacot tapi kaya', 'Customer yang banyak complain tapi sering beli', 'Purple', NOW(), NOW()),
    ('Customer kesayangan', 'Customer loyal dengan transaksi tinggi', 'Pink', NOW(), NOW()),
    ('Hampir Checkout', 'Customer yang sudah di tahap akhir pembelian', 'Green', NOW(), NOW())
ON CONFLICT (label) DO UPDATE SET 
    description = EXCLUDED.description,
    color = EXCLUDED.color,
    updated_at = NOW();

-- =====================================================
-- 4. AI CONFIGURATIONS
-- =====================================================

INSERT INTO ai_configurations (ai_engine, token, endpoint, model, active, created_at, updated_at)
VALUES
    ('openai', 'sk-proj-YOUR-OPENAI-KEY', 'https://api.openai.com/v1/chat/completions', 'gpt-3.5-turbo', true, NOW(), NOW()),
    ('deepseek', 'YOUR-DEEPSEEK-KEY', 'https://api.deepseek.com/v1/chat/completions', 'deepseek-chat', false, NOW(), NOW()),
    ('grok', 'YOUR-GROK-KEY', 'https://api.x.ai/v1/chat/completions', 'grok-beta', false, NOW(), NOW()),
    ('gemini', 'YOUR-GEMINI-KEY', 'https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent', 'gemini-pro', false, NOW(), NOW())
ON CONFLICT (ai_engine) DO UPDATE SET 
    token = EXCLUDED.token,
    endpoint = EXCLUDED.endpoint,
    model = EXCLUDED.model,
    active = EXCLUDED.active,
    updated_at = NOW();

-- =====================================================
-- 5. CONNECTED PLATFORMS
-- =====================================================

INSERT INTO connected_platforms (platform, platform_id, token, phone_number_id, client_id, client_secret, active, created_at, updated_at)
VALUES
    ('WhatsApp', '+6287777888125', 'EAAQaYK3ZAp6sBP8Eo5Jezr84jOHR2B8bRi0577BeEMPc3cKnQ4YXmNMZC3rjrsUL0fV6e4ZAtGAtNPZB9TtTtbp1GEl2bmYkZB8MZAjqwMMAcRJaqbedhfWXX0nurd9UjxLUwu8L2ixUzUHtn6RfjFNdViEf4APLGQ7xvJ4eiTW6ZCHesQkrRzVZBZBfNfubFqqTVuAZDZD', '816432081555850', '1231241241', 'whatsapp-secret', true, NOW(), NOW()),
    ('Telegram', '@testtelegram', '8596779579:AAHnoHUc_Dc-uMUFPUCj97MSpXecS9MznwE', '', '2345234134', 'telegram-secret', true, NOW(), NOW()),
    ('Instagram', '@testinstagram', 'IGAALqhu4OuYdBZAFE0TWtWVkFpVFN6Vm9TVnRaRjl2UzFackJaRllHWE5wTEhZAdjNLamxkVTI5VFRQcExJWWZAGUFFkODk0REEwMDJwcm1QaGVGdDZARV0M1T1U0ZAjBiWm5NOTB2RE82RnIyVFZALRnFGbHctcVZA6VGJucW5EMS1LRQZDZD', '820815193880967', '43623453454', 'instagram-secret', true, NOW(), NOW())
ON CONFLICT (platform) DO UPDATE SET 
    token = EXCLUDED.token,
    platform_id = EXCLUDED.platform_id,
    phone_number_id = EXCLUDED.phone_number_id,
    active = EXCLUDED.active,
    updated_at = NOW();

-- =====================================================
-- 6. AI AGENTS
-- =====================================================

INSERT INTO ai_agents (name, ai_engine, basic_prompt, active, created_at, updated_at)
VALUES
    ('Diva', 'openai', 'Kamu adalah customer service bernama Diva dari Divine CRM. Selalu jawab dalam Bahasa Indonesia dengan ramah dan profesional. Bantu customer dengan pertanyaan produk dan layanan.', true, NOW(), NOW()),
    ('Clara', 'deepseek', 'Kamu adalah customer service bernama Clara dari Divine CRM. Selalu jawab dalam Bahasa Indonesia dengan ramah dan profesional.', false, NOW(), NOW()),
    ('Kana', 'grok', 'Kamu adalah customer service bernama Kana dari Divine CRM. Selalu jawab dalam Bahasa Indonesia dengan ramah dan sedikit humor.', false, NOW(), NOW()),
    ('Gema', 'gemini', 'Kamu adalah customer service bernama Gema dari Divine CRM. Selalu jawab dalam Bahasa Indonesia dengan ramah dan detail.', false, NOW(), NOW())
ON CONFLICT (name) DO UPDATE SET 
    basic_prompt = EXCLUDED.basic_prompt,
    ai_engine = EXCLUDED.ai_engine,
    active = EXCLUDED.active,
    updated_at = NOW();

-- =====================================================
-- 7. HUMAN AGENTS
-- =====================================================

-- Password: password123 (hashed with bcrypt)
INSERT INTO human_agents (username, password, email, full_name, role, active, created_at, updated_at)
VALUES
    ('feli1210', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'feli@divine.com', 'Feli Tan', 'Agent', true, NOW(), NOW()),
    ('nico24', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'nico@divine.com', 'Nikolas Chen', 'Agent', true, NOW(), NOW()),
    ('tomi.siahaan', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'tomi@divine.com', 'Tomi Siahaan', 'Supervisor', true, NOW(), NOW()),
    ('agus123', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'agus@divine.com', 'Agus Santoso', 'Admin', true, NOW(), NOW())
ON CONFLICT (username) DO UPDATE SET 
    email = EXCLUDED.email,
    full_name = EXCLUDED.full_name,
    updated_at = NOW();

-- =====================================================
-- 8. QUICK REPLIES
-- =====================================================

INSERT INTO quick_replies (trigger, response, active, created_at, updated_at)
VALUES
    ('halo', 'Halo! Selamat datang di Divine CRM. Ada yang bisa saya bantu?', true, NOW(), NOW()),
    ('terima kasih', 'Sama-sama! Senang bisa membantu. Ada lagi yang bisa dibantu?', true, NOW(), NOW()),
    ('jam buka', 'Kami melayani Senin-Jumat pukul 09:00-17:00 WIB. Sabtu 09:00-15:00 WIB.', true, NOW(), NOW())
ON CONFLICT (trigger) DO UPDATE SET 
    response = EXCLUDED.response,
    updated_at = NOW();

-- =====================================================
-- 9. BROADCAST TEMPLATES
-- =====================================================

INSERT INTO broadcast_templates (name, content, channel, created_by, active, created_at, updated_at)
VALUES
    ('Promo Weekend', 'Halo {name}! Ada promo spesial weekend untuk Anda. Diskon 20% untuk semua produk. Buruan order sebelum terlambat!', 'All', 'Admin', true, NOW(), NOW()),
    ('Follow Up Order', 'Halo {name}, terima kasih sudah order di Divine CRM. Pesanan Anda sedang diproses. Estimasi tiba dalam 2-3 hari kerja.', 'WhatsApp', 'Admin', true, NOW(), NOW())
ON CONFLICT (name) DO UPDATE SET 
    content = EXCLUDED.content,
    updated_at = NOW();

-- =====================================================
-- VERIFICATION QUERIES
-- =====================================================

-- Show comprehensive status
SELECT '==================== SETUP VERIFICATION ====================' as info;

SELECT 'Contacts' as table_name, COUNT(*) as total FROM contacts
UNION ALL SELECT 'Products', COUNT(*) FROM products
UNION ALL SELECT 'Chat Labels', COUNT(*) FROM chat_labels
UNION ALL SELECT 'AI Configurations', COUNT(*) FROM ai_configurations
UNION ALL SELECT 'Connected Platforms', COUNT(*) FROM connected_platforms
UNION ALL SELECT 'AI Agents', COUNT(*) FROM ai_agents
UNION ALL SELECT 'Human Agents', COUNT(*) FROM human_agents
UNION ALL SELECT 'Quick Replies', COUNT(*) FROM quick_replies
UNION ALL SELECT 'Broadcast Templates', COUNT(*) FROM broadcast_templates;

-- Show active configurations
SELECT '==================== ACTIVE CONFIGURATIONS ====================' as info;

SELECT 'AI Config' as type, ai_engine as name, active, model as detail FROM ai_configurations WHERE active = true
UNION ALL
SELECT 'AI Agent', name, active, ai_engine FROM ai_agents WHERE active = true
UNION ALL
SELECT 'Platform', platform, active, platform_id FROM connected_platforms WHERE active = true;

-- ✅ Setup Complete!
SELECT '==================== ✅ SETUP COMPLETE! ====================' as info;