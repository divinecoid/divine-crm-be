-- Insert sample knowledge base (without embeddings for now)
INSERT INTO knowledge_bases (
    title,
    content,
    category,
    tags,
    source,
    active,
    created_at,
    updated_at
) VALUES 
(
    'Tentang Perusahaan',
    'Divine CRM adalah platform CRM berbasis AI yang membantu bisnis mengelola customer relationship dengan lebih efektif. Kami menyediakan integrasi WhatsApp, Instagram, dan Telegram dengan AI-powered responses.',
    'company_info',
    'company,about,crm',
    'internal',
    true,
    NOW(),
    NOW()
),
(
    'Fitur Utama',
    'Divine CRM memiliki fitur utama: 1) AI-powered chat responses dengan RAG (Retrieval Augmented Generation), 2) Multi-platform messaging yang mendukung WhatsApp, Instagram, dan Telegram, 3) Contact management untuk mengelola pelanggan, 4) Product catalog untuk menampilkan produk, 5) Broadcast messaging untuk kirim pesan massal, 6) Analytics dan reporting untuk analisa bisnis',
    'features',
    'features,capabilities,fitur',
    'internal',
    true,
    NOW(),
    NOW()
),
(
    'Harga dan Paket',
    'Kami menawarkan beberapa paket: Paket Basic seharga Rp 500.000 per bulan untuk 100 contacts, Paket Pro seharga Rp 1.500.000 per bulan untuk 500 contacts, dan Paket Enterprise dengan harga custom untuk unlimited contacts. Semua paket sudah termasuk AI responses, multi-platform, dan support 24/7.',
    'pricing',
    'pricing,subscription,payment,harga,biaya',
    'internal',
    true,
    NOW(),
    NOW()
),
(
    'Cara Pemesanan',
    'Untuk memesan Divine CRM sangat mudah: 1) Hubungi tim sales kami melalui WhatsApp atau email, 2) Pilih paket yang sesuai dengan kebutuhan bisnis Anda, 3) Lakukan pembayaran via transfer bank atau virtual account, 4) Tim kami akan setup akun Anda dalam waktu maksimal 1x24 jam, 5) Dapatkan training gratis untuk tim Anda',
    'sales',
    'order,purchase,how to buy,pemesanan,beli',
    'internal',
    true,
    NOW(),
    NOW()
),
(
    'Support dan Bantuan',
    'Tim support Divine CRM tersedia 24/7 untuk membantu Anda. Anda bisa menghubungi kami via WhatsApp di nomor +62812-3456-7890, email di support@divine-crm.com, atau melalui live chat di website. Response time rata-rata kami adalah 15 menit untuk urgent issues dan 1 jam untuk general inquiries.',
    'support',
    'help,support,contact,bantuan,cs',
    'internal',
    true,
    NOW(),
    NOW()
),
(
    'Demo dan Trial',
    'Kami menyediakan free trial selama 14 hari tanpa perlu kartu kredit. Anda bisa mencoba semua fitur Divine CRM secara gratis. Untuk demo live dengan tim kami, silakan booking jadwal melalui WhatsApp atau email. Demo biasanya berlangsung 30-45 menit.',
    'trial',
    'demo,trial,free,coba,gratis',
    'internal',
    true,
    NOW(),
    NOW()
),
(
    'Integrasi WhatsApp',
    'Divine CRM terintegrasi penuh dengan WhatsApp Business API. Fitur termasuk: auto-reply dengan AI, broadcast messaging, chat assignment ke agent, quick replies, media support (gambar, video, dokumen), dan WhatsApp template messages untuk notifikasi.',
    'whatsapp',
    'whatsapp,wa,integration,integrasi',
    'internal',
    true,
    NOW(),
    NOW()
),
(
    'Keamanan Data',
    'Keamanan data Anda adalah prioritas kami. Divine CRM menggunakan enkripsi end-to-end, server di Indonesia, backup harian otomatis, dan compliant dengan regulasi perlindungan data. Kami tidak pernah membagikan data pelanggan Anda ke pihak ketiga.',
    'security',
    'security,privacy,data,keamanan,privasi',
    'internal',
    true,
    NOW(),
    NOW()
);

-- Verify insertion
SELECT id, title, category, active FROM knowledge_bases;