-- ========================================
-- Divine CRM Multi-Platform AI Migration
-- ========================================
-- Updated to match actual Divine CRM schema

-- ========================================
-- STEP 1: Ensure pgvector Extension
-- ========================================

CREATE EXTENSION IF NOT EXISTS vector;

-- ========================================
-- STEP 2: Create AI Agents Table (if not exists)
-- ========================================

CREATE TABLE IF NOT EXISTS ai_agents (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    ai_engine VARCHAR(50) NOT NULL,
    basic_prompt TEXT,
    system_prompt TEXT,
    temperature FLOAT DEFAULT 0.7,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ========================================
-- STEP 3: Create AI Configurations Table (if not exists)
-- ========================================

CREATE TABLE IF NOT EXISTS ai_configurations (
    id BIGSERIAL PRIMARY KEY,
    ai_engine VARCHAR(50) NOT NULL UNIQUE,
    token TEXT NOT NULL,
    endpoint VARCHAR(500),
    model VARCHAR(100),
    max_tokens INTEGER DEFAULT 2000,
    temperature FLOAT DEFAULT 0.7,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ========================================
-- STEP 4: Add AI Agent column to Connected Platforms
-- ========================================

ALTER TABLE connected_platforms 
ADD COLUMN IF NOT EXISTS ai_agent_id bigint REFERENCES ai_agents(id) ON DELETE SET NULL;

-- ========================================
-- STEP 5: Create AI Agents (Diva, Clara, Kana, Gema)
-- ========================================

-- Diva for WhatsApp + OpenAI
INSERT INTO ai_agents (name, ai_engine, basic_prompt, active, created_at, updated_at)
VALUES (
    'Diva',
    'openai',
    'Nama Anda adalah Diva, customer service AI untuk WhatsApp Divine CRM.
Anda ramah, profesional, dan selalu siap membantu customer.
Jawab pertanyaan dengan detail dan gunakan emoji untuk lebih friendly.
Jika customer ingin melanjutkan ke human agent, tawarkan dengan baik.',
    true,
    NOW(),
    NOW()
) ON CONFLICT (name) DO NOTHING;

-- Clara for Telegram + DeepSeek
INSERT INTO ai_agents (name, ai_engine, basic_prompt, active, created_at, updated_at)
VALUES (
    'Clara',
    'deepseek',
    'Nama Anda adalah Clara, customer service AI untuk Telegram Divine CRM.
Anda expert dalam memberikan informasi produk dan penawaran spesial.
Selalu jaga komunikasi tetap ringkas karena medium Telegram.
Gunakan formatting yang rapi dan mudah dibaca.',
    true,
    NOW(),
    NOW()
) ON CONFLICT (name) DO NOTHING;

-- Kana for Instagram + Grok
INSERT INTO ai_agents (name, ai_engine, basic_prompt, active, created_at, updated_at)
VALUES (
    'Kana',
    'grok',
    'Nama Anda adalah Kana, customer service AI untuk Instagram Direct Message.
Anda stylish, trendy, dan sangat memahami social media culture.
Jawab dengan bahasa yang casual tapi tetap profesional.
Sering gunakan emoji dan emoticon sesuai konteks percakapan.',
    true,
    NOW(),
    NOW()
) ON CONFLICT (name) DO NOTHING;

-- Gema for Gemini (Multi-Platform)
INSERT INTO ai_agents (name, ai_engine, basic_prompt, active, created_at, updated_at)
VALUES (
    'Gema',
    'gemini',
    'Nama Anda adalah Gema, customer service AI universal Divine CRM.
Anda versatile dan bisa beradaptasi dengan berbagai platform dan customer.
Selalu berikan respons yang personal dan relevan dengan kebutuhan customer.
Prioritaskan kepuasan customer dan retention.',
    true,
    NOW(),
    NOW()
) ON CONFLICT (name) DO NOTHING;

-- ========================================
-- STEP 6: Setup Connected Platforms with AI Agents
-- ========================================

-- Update WhatsApp platforms to use Diva (OpenAI)
UPDATE connected_platforms 
SET ai_agent_id = (SELECT id FROM ai_agents WHERE name = 'Diva')
WHERE platform = 'WhatsApp' AND ai_agent_id IS NULL;

-- Update Telegram platforms to use Clara (DeepSeek)
UPDATE connected_platforms 
SET ai_agent_id = (SELECT id FROM ai_agents WHERE name = 'Clara')
WHERE platform = 'Telegram' AND ai_agent_id IS NULL;

-- Update Instagram platforms to use Kana (Grok)
UPDATE connected_platforms 
SET ai_agent_id = (SELECT id FROM ai_agents WHERE name = 'Kana')
WHERE platform = 'Instagram' AND ai_agent_id IS NULL;

-- ========================================
-- STEP 7: Insert AI Configurations
-- ========================================

-- OpenAI Configuration
INSERT INTO ai_configurations (ai_engine, token, endpoint, model, active, created_at, updated_at)
VALUES (
    'openai',
    '',
    'https://api.openai.com/v1/chat/completions',
    'gpt-4o-mini',
    true,
    NOW(),
    NOW()
) ON CONFLICT (ai_engine) DO UPDATE SET updated_at = NOW();

-- DeepSeek Configuration
INSERT INTO ai_configurations (ai_engine, token, endpoint, model, active, created_at, updated_at)
VALUES (
    'deepseek',
    '',
    'https://api.deepseek.com/v1/chat/completions',
    'deepseek-chat',
    true,
    NOW(),
    NOW()
) ON CONFLICT (ai_engine) DO UPDATE SET updated_at = NOW();

-- Grok Configuration (xAI)
INSERT INTO ai_configurations (ai_engine, token, endpoint, model, active, created_at, updated_at)
VALUES (
    'grok',
    '',
    'https://api.x.ai/v1/chat/completions',
    'grok-beta',
    true,
    NOW(),
    NOW()
) ON CONFLICT (ai_engine) DO UPDATE SET updated_at = NOW();

-- Gemini Configuration
INSERT INTO ai_configurations (ai_engine, token, endpoint, model, active, created_at, updated_at)
VALUES (
    'gemini',
    '',
    'https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent',
    'gemini-pro',
    true,
    NOW(),
    NOW()
) ON CONFLICT (ai_engine) DO UPDATE SET updated_at = NOW();

-- ========================================
-- STEP 8: Create Vector Embedding Tables
-- ========================================

-- Chat Histories with message and response embeddings
CREATE TABLE IF NOT EXISTS chat_histories (
    id BIGSERIAL PRIMARY KEY,
    contact_id BIGINT NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    response TEXT,
    message_embedding vector(1536),
    response_embedding vector(1536),
    platform VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Knowledge Base with embeddings
CREATE TABLE IF NOT EXISTS knowledge_bases (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    embedding vector(1536),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- FAQ Embeddings
CREATE TABLE IF NOT EXISTS faq_embeddings (
    id BIGSERIAL PRIMARY KEY,
    question VARCHAR(500) NOT NULL,
    answer TEXT NOT NULL,
    embedding vector(1536),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Product Embeddings
CREATE TABLE IF NOT EXISTS product_embeddings (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    embedding vector(1536),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ========================================
-- STEP 9: Create Indexes for Vector Search
-- ========================================

CREATE INDEX IF NOT EXISTS idx_chat_history_contact_id ON chat_histories(contact_id);
CREATE INDEX IF NOT EXISTS idx_knowledge_base_active ON knowledge_bases(active);
CREATE INDEX IF NOT EXISTS idx_faq_active ON faq_embeddings(active);
CREATE INDEX IF NOT EXISTS idx_product_embedding_product_id ON product_embeddings(product_id);

-- Vector indexes for similarity search (pgvector)
CREATE INDEX IF NOT EXISTS idx_knowledge_base_embedding ON knowledge_bases USING ivfflat (embedding vector_cosine_ops)
    WITH (lists = 100);

CREATE INDEX IF NOT EXISTS idx_chat_history_message_embedding ON chat_histories USING ivfflat (message_embedding vector_cosine_ops)
    WITH (lists = 100);

CREATE INDEX IF NOT EXISTS idx_chat_history_response_embedding ON chat_histories USING ivfflat (response_embedding vector_cosine_ops)
    WITH (lists = 100);

CREATE INDEX IF NOT EXISTS idx_faq_embedding ON faq_embeddings USING ivfflat (embedding vector_cosine_ops)
    WITH (lists = 100);

CREATE INDEX IF NOT EXISTS idx_product_embedding_embedding ON product_embeddings USING ivfflat (embedding vector_cosine_ops)
    WITH (lists = 100);

-- ========================================
-- Verification Queries (uncomment to verify)
-- ========================================

-- SELECT id, name, ai_engine FROM ai_agents ORDER BY name;
-- SELECT ai_engine, model FROM ai_configurations ORDER BY ai_engine;
-- SELECT p.id, p.platform, a.name, a.ai_engine 
-- FROM connected_platforms p
-- LEFT JOIN ai_agents a ON p.ai_agent_id = a.id
-- ORDER BY p.platform;