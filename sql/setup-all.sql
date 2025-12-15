-- ==================== DIVINE CRM - COMPLETE DATABASE SEED ====================
-- Sesuai wireframe design
-- Execute: docker exec -it divine-postgres psql -U postgres -d divine_crm -f /tmp/seed.sql

-- ==================== 1. CHAT LABELS (Sesuai Wireframe Page 3) ====================
-- Columns: Label, Description, Label Color
INSERT INTO chat_labels (code, name, color, description, created_at, updated_at) VALUES
('TANYA', 'Nanya-nanya doang', '#ef4444', 'Customer yang cuma nanya tapi belum beli', NOW(), NOW()),
('VIP_BACOT', 'Customer bacot tapi kaya', '#8b5cf6', 'Customer yang banyak tanya tapi spending tinggi', NOW(), NOW()),
('KESAYANGAN', 'Customer kesayangan', '#ec4899', 'Customer loyal dan mudah dihandle', NOW(), NOW()),
('CHECKOUT', 'Hampir Checkout', '#10b981', 'Customer yang sudah siap bayar', NOW(), NOW()),
('HOT', 'Hot Lead', '#f97316', 'Prospek dengan kemungkinan closing tinggi', NOW(), NOW()),
('WARM', 'Warm Lead', '#eab308', 'Prospek yang menunjukkan ketertarikan', NOW(), NOW()),
('COLD', 'Cold Lead', '#3b82f6', 'Prospek baru belum menunjukkan interest', NOW(), NOW()),
('KOMPLAIN', 'Komplain', '#dc2626', 'Customer yang sedang komplain', NOW(), NOW()),
('FOLLOWUP', 'Perlu Follow Up', '#14b8a6', 'Perlu di-follow up', NOW(), NOW()),
('PENDING', 'Pending Payment', '#f59e0b', 'Menunggu pembayaran', NOW(), NOW())
ON CONFLICT (code) DO UPDATE SET 
    name = EXCLUDED.name,
    color = EXCLUDED.color,
    description = EXCLUDED.description,
    updated_at = NOW();

-- ==================== 2. AI CONFIGURATIONS (Sesuai Wireframe Page 3) ====================
-- Columns: AI Engine, Token
-- OpenAI, Deepseek, Grok xAI, Gemini
INSERT INTO ai_configurations (code, name, provider, model, api_endpoint, temperature, max_tokens, is_active, is_default, created_at, updated_at) VALUES
('OPENAI', 'OpenAI', 'openai', 'gpt-3.5-turbo', 'https://api.openai.com/v1/chat/completions', 0.7, 1000, true, true, NOW(), NOW()),
('DEEPSEEK', 'Deepseek', 'deepseek', 'deepseek-chat', 'https://api.deepseek.com/v1/chat/completions', 0.7, 1000, true, false, NOW(), NOW()),
('GROK', 'Grok xAI', 'grok', 'grok-beta', 'https://api.x.ai/v1/chat/completions', 0.7, 1000, true, false, NOW(), NOW()),
('GEMINI', 'Gemini', 'gemini', 'gemini-pro', 'https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent', 0.7, 1000, true, false, NOW(), NOW())
ON CONFLICT (code) DO UPDATE SET 
    name = EXCLUDED.name,
    updated_at = NOW();

-- ==================== 3. CONNECTED PLATFORMS (Sesuai Wireframe Page 3-5) ====================
-- Columns: Platform, ID, Access Token, Client ID, Client Secret
INSERT INTO connected_platforms (code, name, platform_type, is_active, config, created_at, updated_at) VALUES
('WHATSAPP', 'WhatsApp', 'whatsapp', true, 
 '{"phone_id": "+6287777888125", "access_token": "EAAQaYK3ZAp6sBP...", "client_id": "1231241241", "client_secret": "srgsdfgdsfgdfg"}', 
 NOW(), NOW()),
('TELEGRAM', 'Telegram', 'telegram', true, 
 '{"bot_id": "@testtelegram", "access_token": "8310666863:AAHd4nW...", "client_id": "2345", "client_secret": "srgsdfgdsfgdfg"}', 
 NOW(), NOW()),
('INSTAGRAM', 'Instagram', 'instagram', true, 
 '{"username": "@testinstagram", "access_token": "EAAQaYK3ZAp6sBPx...", "client_id": "43623", "client_secret": "srgsdfgdsfgdfg"}', 
 NOW(), NOW())
ON CONFLICT (code) DO UPDATE SET 
    config = EXCLUDED.config,
    updated_at = NOW();

-- ==================== 4. AI AGENTS (Sesuai Wireframe Page 5 & 7) ====================
-- Columns: AI Agent, AI Engine, Basic Prompt
-- Diva → OpenAI, Clara → Deepseek, Kana → Grok xAI, Gema → Gemini

INSERT INTO ai_agents (code, name, platform, ai_config_id, system_prompt, welcome_message, is_active, is_default, created_at, updated_at) VALUES

-- DIVA - OpenAI (WhatsApp Sales)
('DIVA', 'Diva', 'whatsapp', 
 (SELECT id FROM ai_configurations WHERE code = 'OPENAI'),
 'Kamu adalah Diva, customer service dari Divine CRM.

🔴 ATURAN WAJIB PERKENALAN:
Setiap kali customer menyapa (halo, hi, hai, hello, hey, selamat pagi/siang/sore/malam, assalamualaikum, permisi, p, hallo, dll), WAJIB perkenalkan diri dengan format:
"Halo! 👋 Perkenalkan, saya Diva dari Divine CRM. Ada yang bisa Diva bantu?"

KEPRIBADIAN:
- Ramah, profesional, dan helpful
- Bahasa Indonesia sopan tapi friendly
- Gunakan emoji secukupnya
- SELALU sebut nama "Diva" saat berinteraksi

TUGAS:
1. Perkenalkan diri saat disapa
2. Jawab pertanyaan produk & layanan
3. Bantu proses pembelian
4. Arahkan ke human agent jika perlu

Jika tidak bisa jawab: "Diva akan hubungkan dengan tim kami ya!"',
 'Halo! 👋 Perkenalkan, saya Diva dari Divine CRM!

Diva siap membantu:
✨ Info produk & layanan
💰 Konsultasi paket
🎁 Info promo

Ada yang bisa Diva bantu? 😊',
 true, true, NOW(), NOW()),

-- CLARA - Deepseek (WhatsApp Support)
('CLARA', 'Clara', 'whatsapp',
 (SELECT id FROM ai_configurations WHERE code = 'DEEPSEEK'),
 'Kamu adalah Clara, customer service dari Divine CRM.

🔴 ATURAN WAJIB PERKENALAN:
Setiap kali customer menyapa (halo, hi, hai, hello, hey, selamat pagi/siang/sore/malam, assalamualaikum, permisi, p, hallo, dll), WAJIB perkenalkan diri dengan format:
"Halo! 👋 Saya Clara dari Divine CRM. Ada yang bisa Clara bantu?"

KEPRIBADIAN:
- Sabar, empati, dan solutif
- Tenang menghadapi komplain
- Detail memahami masalah
- SELALU sebut nama "Clara"

TUGAS:
1. Perkenalkan diri saat disapa
2. Bantu selesaikan masalah teknis
3. Jawab pertanyaan penggunaan
4. Eskalasi jika masalah kompleks

Jika tidak bisa jawab: "Clara akan hubungkan dengan tim specialist ya!"',
 'Halo! 👋 Saya Clara dari Divine CRM.

Clara siap bantu kendala yang kamu alami. Ceritakan masalahnya ya! 💪

Ada yang bisa Clara bantu?',
 true, false, NOW(), NOW()),

-- KANA - Grok xAI (Telegram)
('KANA', 'Kana', 'telegram',
 (SELECT id FROM ai_configurations WHERE code = 'GROK'),
 'Kamu adalah Kana, customer service dari Divine CRM di Telegram.

🔴 ATURAN WAJIB PERKENALAN:
Setiap kali user menyapa (halo, hi, hai, hello, /start, selamat pagi/siang/sore/malam, p, dll), WAJIB perkenalkan diri dengan format:
"Hai! 🎉 Aku Kana dari Divine CRM! Senang ketemu kamu!"

KEPRIBADIAN:
- Ceria, energik, dan friendly
- Bahasa santai (pakai "aku" dan "kamu")
- Suka emoji yang fun 🎉✨😊
- SELALU sebut nama "Kana"

TUGAS:
1. Perkenalkan diri dengan ceria
2. Jawab pertanyaan umum
3. Info produk dan promo
4. Bantu navigasi bot

Gunakan bahasa seperti ngobrol sama teman!',
 'Hai hai! 🎉

Aku Kana dari Divine CRM! Senang ketemu kamu! ✨

Kana bisa bantu:
📦 Info produk
💡 Tips penggunaan
🎁 Update promo

Langsung chat Kana ya! 😊',
 true, true, NOW(), NOW()),

-- GEMA - Gemini (Instagram)
('GEMA', 'Gema', 'instagram',
 (SELECT id FROM ai_configurations WHERE code = 'GEMINI'),
 'Kamu adalah Gema, customer service dari Divine CRM di Instagram.

🔴 ATURAN WAJIB PERKENALAN:
Setiap kali user DM dengan sapaan (halo, hi, hai, hello, heyy, p, dll), WAJIB perkenalkan diri dengan format:
"Haii! ✨ Aku Gema dari Divine CRM! Thanks udah DM!"

KEPRIBADIAN:
- Trendy, gaul, dan relatable
- Bahasa anak muda
- Emoji aesthetic ✨💕🔥
- SELALU sebut nama "Gema"

TUGAS:
1. Perkenalkan diri yang catchy
2. Respons DM cepat dan friendly
3. Info promo dan diskon
4. Share link pembelian

Gunakan bahasa Instagram-style!',
 'Haii! ✨

Aku Gema dari Divine CRM! Thanks udah DM! 🙌

Gema bisa bantu:
🛍️ Produk & Layanan
🔥 Promo & Diskon
💬 Tanya-tanya

Chat Gema ya! 💕',
 true, true, NOW(), NOW())

ON CONFLICT (code) DO UPDATE SET 
    name = EXCLUDED.name,
    system_prompt = EXCLUDED.system_prompt,
    welcome_message = EXCLUDED.welcome_message,
    updated_at = NOW();

-- ==================== 5. PRODUCTS (Sesuai Wireframe Page 1) ====================
-- Columns: Code, Name, Price, Stock, Uploaded By
INSERT INTO products (code, name, description, price, category, stock, is_active, created_at, updated_at) VALUES
-- Sample dari wireframe
('F003', 'Produk 1', 'Deskripsi produk 1', 28000, 'General', 853, true, NOW(), NOW()),
('V005', 'Produk 2', 'Deskripsi produk 2', 58000, 'General', 347, true, NOW(), NOW()),
('O002', 'Produk 3', 'Deskripsi produk 3', 352000, 'General', 495, true, NOW(), NOW()),
-- CRM Packages
('CRM01', 'Divine CRM Starter', 'Paket CRM untuk bisnis kecil. 1 user, 1000 kontak, WhatsApp, Basic reporting.', 299000, 'CRM Package', 999, true, NOW(), NOW()),
('CRM02', 'Divine CRM Business', 'Paket CRM bisnis. 5 users, 10000 kontak, Multi-channel, AI Assistant.', 799000, 'CRM Package', 999, true, NOW(), NOW()),
('CRM03', 'Divine CRM Enterprise', 'Paket enterprise. Unlimited users & kontak, All integrations.', 1999000, 'CRM Package', 999, true, NOW(), NOW()),
-- Add-ons
('ADD01', 'WhatsApp Business API', 'Integrasi WhatsApp Business API', 150000, 'Add-on', 999, true, NOW(), NOW()),
('ADD02', 'Telegram Bot', 'Integrasi Telegram Bot', 100000, 'Add-on', 999, true, NOW(), NOW()),
('ADD03', 'Instagram DM', 'Integrasi Instagram Direct Message', 100000, 'Add-on', 999, true, NOW(), NOW()),
('ADD04', 'AI Assistant Pro', 'AI dengan 100.000 token/bulan', 299000, 'Add-on', 999, true, NOW(), NOW()),
('ADD05', 'Broadcast Pro', 'Broadcast unlimited + scheduling', 199000, 'Add-on', 999, true, NOW(), NOW())
ON CONFLICT (code) DO UPDATE SET 
    name = EXCLUDED.name,
    price = EXCLUDED.price,
    stock = EXCLUDED.stock,
    updated_at = NOW();

-- ==================== 6. BROADCAST TEMPLATES ====================
INSERT INTO broadcast_templates (code, name, category, content, variables, platform, is_active, created_at, updated_at) VALUES
('WELCOME', 'Welcome Customer', 'Welcome', 
 'Halo {{name}}! 👋\n\nSelamat datang di Divine CRM!\n\nSalam,\nTim Divine CRM ✨',
 '["name"]', 'all', true, NOW(), NOW()),
('PROMO', 'Promo Discount', 'Promotion',
 '🎉 PROMO! 🎉\n\nHai {{name}},\nDiskon {{discount}}%!\nBerlaku hingga {{end_date}}\n\nKode: {{promo_code}}',
 '["name", "discount", "end_date", "promo_code"]', 'all', true, NOW(), NOW()),
('FLASH', 'Flash Sale', 'Promotion',
 '⚡ FLASH SALE ⚡\n\n{{name}}, diskon {{discount}}% HARI INI!\nBerakhir {{end_time}} WIB',
 '["name", "discount", "end_time"]', 'all', true, NOW(), NOW()),
('FOLLOWUP', 'Follow Up', 'Follow Up',
 'Hai {{name}},\n\nTerima kasih sudah menghubungi kami tentang {{product}}.\nAda pertanyaan lain? 😊',
 '["name", "product"]', 'all', true, NOW(), NOW()),
('ORDER', 'Order Confirmation', 'Notification',
 '✅ PESANAN DIKONFIRMASI\n\nHai {{name}},\n📦 Order: {{order_id}}\n💰 Total: Rp {{total}}\n\nTerima kasih! 🙏',
 '["name", "order_id", "total"]', 'all', true, NOW(), NOW()),
('PAYMENT', 'Payment Reminder', 'Notification',
 '⏰ REMINDER\n\nHai {{name}},\nPesanan {{order_id}} menunggu pembayaran.\n💰 Total: Rp {{total}}\n⏳ Batas: {{deadline}}',
 '["name", "order_id", "total", "deadline"]', 'all', true, NOW(), NOW()),
('BIRTHDAY', 'Birthday', 'Seasonal',
 '🎂 HAPPY BIRTHDAY {{name}}! 🎉\n\nDari tim Divine CRM.\nHadiah: Diskon {{discount}}%\nKode: {{promo_code}}',
 '["name", "discount", "promo_code"]', 'all', true, NOW(), NOW())
ON CONFLICT (code) DO UPDATE SET 
    content = EXCLUDED.content,
    updated_at = NOW();

-- ==================== 7. QUICK REPLIES ====================
INSERT INTO quick_replies (code, title, content, category, is_active, created_at, updated_at) VALUES
('HALO', 'Salam Pembuka', 'Halo! 👋 Terima kasih sudah menghubungi Divine CRM. Ada yang bisa kami bantu?', 'Greeting', true, NOW(), NOW()),
('PAGI', 'Selamat Pagi', 'Selamat pagi! ☀️ Ada yang bisa kami bantu?', 'Greeting', true, NOW(), NOW()),
('SIANG', 'Selamat Siang', 'Selamat siang! 🌤️ Ada yang bisa kami bantu?', 'Greeting', true, NOW(), NOW()),
('MALAM', 'Selamat Malam', 'Selamat malam! 🌙 Ada yang bisa dibantu?', 'Greeting', true, NOW(), NOW()),
('HARGA', 'Daftar Harga', 'Daftar harga:\n💼 Starter: Rp 299.000/bln\n🏢 Business: Rp 799.000/bln\n🏛️ Enterprise: Rp 1.999.000/bln', 'Product', true, NOW(), NOW()),
('PROMO', 'Info Promo', 'Saat ini ada promo spesial! 🎉 Hubungi tim sales untuk info.', 'Product', true, NOW(), NOW()),
('TUNGGU', 'Mohon Tunggu', 'Mohon tunggu sebentar ya! ⏳', 'Support', true, NOW(), NOW()),
('ESKALASI', 'Eskalasi', 'Saya hubungkan dengan tim specialist. Mohon tunggu. 👨‍💼', 'Support', true, NOW(), NOW()),
('SELESAI', 'Masalah Selesai', 'Senang bisa membantu! ✅ Ada lagi yang bisa dibantu? 🙏', 'Closing', true, NOW(), NOW()),
('TUTUP', 'Penutup', 'Terima kasih! 🙏 Semoga hari Anda menyenangkan! ✨', 'Closing', true, NOW(), NOW()),
('JAM', 'Jam Operasional', 'Jam operasional:\n🕐 Senin-Jumat: 08.00-17.00\n🕐 Sabtu: 08.00-12.00\n🚫 Minggu: Tutup', 'Others', true, NOW(), NOW())
ON CONFLICT (code) DO UPDATE SET 
    content = EXCLUDED.content,
    updated_at = NOW();

-- ==================== VERIFY DATA ====================
SELECT '═══════════════════════════════════════════════════════' as line;
SELECT '✅ DIVINE CRM DATABASE SEED COMPLETED!' as status;
SELECT '═══════════════════════════════════════════════════════' as line;

SELECT '' as spacer;
SELECT '📊 DATA SUMMARY:' as header;
SELECT '───────────────────────────────────────────────────────' as line;

SELECT 'Chat Labels' as "Table", COUNT(*)::text as "Count" FROM chat_labels
UNION ALL SELECT 'AI Configurations', COUNT(*)::text FROM ai_configurations
UNION ALL SELECT 'Connected Platforms', COUNT(*)::text FROM connected_platforms
UNION ALL SELECT 'AI Agents', COUNT(*)::text FROM ai_agents
UNION ALL SELECT 'Products', COUNT(*)::text FROM products
UNION ALL SELECT 'Broadcast Templates', COUNT(*)::text FROM broadcast_templates
UNION ALL SELECT 'Quick Replies', COUNT(*)::text FROM quick_replies;

SELECT '' as spacer;
SELECT '🤖 AI AGENTS MAPPING:' as header;
SELECT '───────────────────────────────────────────────────────' as line;

SELECT 
    a.name as "Agent", 
    c.name as "AI Engine",
    a.platform as "Platform",
    CASE WHEN a.is_default THEN '✓' ELSE '' END as "Default"
FROM ai_agents a
JOIN ai_configurations c ON a.ai_config_id = c.id
ORDER BY 
    CASE a.name 
        WHEN 'Diva' THEN 1 
        WHEN 'Clara' THEN 2 
        WHEN 'Kana' THEN 3 
        WHEN 'Gema' THEN 4 
    END;

SELECT '' as spacer;
SELECT '🏷️ CHAT LABELS:' as header;
SELECT '───────────────────────────────────────────────────────' as line;
SELECT name as "Label", color as "Color" FROM chat_labels LIMIT 5;