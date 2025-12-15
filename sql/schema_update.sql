-- ==================== DIVINE CRM - SCHEMA UPDATE ====================
-- Tambah kolom yang kurang di database
-- Execute: docker exec -it divine-postgres psql -U postgres -d divine_crm -f /tmp/schema_update.sql

-- ==================== AI_AGENTS - Tambah kolom ====================
ALTER TABLE ai_agents ADD COLUMN IF NOT EXISTS platform VARCHAR(50);
ALTER TABLE ai_agents ADD COLUMN IF NOT EXISTS intro_message TEXT;

-- Update existing data
UPDATE ai_agents SET platform = 'whatsapp' WHERE name IN ('Diva', 'Clara');
UPDATE ai_agents SET platform = 'telegram' WHERE name = 'Kana';
UPDATE ai_agents SET platform = 'instagram' WHERE name = 'Gema';

UPDATE ai_agents SET intro_message = 'Halo! 👋 Perkenalkan, saya Diva dari Divine CRM. Ada yang bisa Diva bantu?' WHERE name = 'Diva';
UPDATE ai_agents SET intro_message = 'Halo! 👋 Saya Clara dari Divine CRM. Ada yang bisa Clara bantu?' WHERE name = 'Clara';
UPDATE ai_agents SET intro_message = 'Hai! 🎉 Aku Kana dari Divine CRM! Senang ketemu kamu!' WHERE name = 'Kana';
UPDATE ai_agents SET intro_message = 'Haii! ✨ Aku Gema dari Divine CRM! Thanks udah DM!' WHERE name = 'Gema';

-- ==================== CONNECTED_PLATFORMS - Tambah kolom ====================
ALTER TABLE connected_platforms ADD COLUMN IF NOT EXISTS token TEXT;
ALTER TABLE connected_platforms ADD COLUMN IF NOT EXISTS client_id VARCHAR(255);
ALTER TABLE connected_platforms ADD COLUMN IF NOT EXISTS client_secret VARCHAR(255);
ALTER TABLE connected_platforms ADD COLUMN IF NOT EXISTS webhook_url TEXT;
ALTER TABLE connected_platforms ADD COLUMN IF NOT EXISTS phone_number_id VARCHAR(100);
ALTER TABLE connected_platforms ADD COLUMN IF NOT EXISTS active BOOLEAN DEFAULT true;

-- Drop NOT NULL constraint on contact_id if it exists (karena platform bisa tanpa contact)
ALTER TABLE connected_platforms ALTER COLUMN contact_id DROP NOT NULL;

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_connected_platforms_active ON connected_platforms(active);

-- ==================== VERIFY ====================
SELECT '✅ Schema update completed!' as status;

SELECT '' as spacer;
SELECT '📋 AI_AGENTS columns:' as info;
SELECT column_name, data_type FROM information_schema.columns 
WHERE table_name = 'ai_agents' ORDER BY ordinal_position;

SELECT '' as spacer;
SELECT '📋 CONNECTED_PLATFORMS columns:' as info;
SELECT column_name, data_type FROM information_schema.columns 
WHERE table_name = 'connected_platforms' ORDER BY ordinal_position;