-- ========================================
-- Divine CRM Base Schema
-- ========================================
-- This creates all core tables that setup-multi-ai.sql depends on

-- ========================================
-- STEP 1: Enable Extensions
-- ========================================

CREATE EXTENSION IF NOT EXISTS vector;

-- ========================================
-- STEP 2: Create Core Tables
-- ========================================

-- Contacts Table
CREATE TABLE IF NOT EXISTS contacts (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(20),
    whatsapp_id VARCHAR(255),
    telegram_id VARCHAR(255),
    instagram_id VARCHAR(255),
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(email),
    UNIQUE(phone)
);

-- AI Agents Table
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

-- AI Configurations Table
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

-- Connected Platforms Table
CREATE TABLE IF NOT EXISTS connected_platforms (
    id BIGSERIAL PRIMARY KEY,
    contact_id BIGINT NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    platform VARCHAR(50) NOT NULL,
    platform_id VARCHAR(255) NOT NULL,
    platform_name VARCHAR(255),
    status VARCHAR(50) DEFAULT 'active',
    ai_agent_id BIGINT REFERENCES ai_agents(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(platform, platform_id)
);

-- Products Table
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(12, 2),
    sku VARCHAR(100) UNIQUE,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Conversations Table
CREATE TABLE IF NOT EXISTS conversations (
    id BIGSERIAL PRIMARY KEY,
    contact_id BIGINT NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    platform VARCHAR(50),
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Messages Table
CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_type VARCHAR(20), -- 'customer' or 'ai' or 'human'
    sender_id VARCHAR(255),
    content TEXT NOT NULL,
    message_type VARCHAR(50) DEFAULT 'text', -- 'text', 'image', 'document', etc.
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ========================================
-- STEP 3: Create Indexes for Performance
-- ========================================

CREATE INDEX IF NOT EXISTS idx_contacts_email ON contacts(email);
CREATE INDEX IF NOT EXISTS idx_contacts_phone ON contacts(phone);
CREATE INDEX IF NOT EXISTS idx_contacts_whatsapp_id ON contacts(whatsapp_id);
CREATE INDEX IF NOT EXISTS idx_contacts_telegram_id ON contacts(telegram_id);
CREATE INDEX IF NOT EXISTS idx_contacts_instagram_id ON contacts(instagram_id);

CREATE INDEX IF NOT EXISTS idx_ai_agents_engine ON ai_agents(ai_engine);
CREATE INDEX IF NOT EXISTS idx_ai_agents_active ON ai_agents(active);

CREATE INDEX IF NOT EXISTS idx_ai_configurations_engine ON ai_configurations(ai_engine);
CREATE INDEX IF NOT EXISTS idx_ai_configurations_active ON ai_configurations(active);

CREATE INDEX IF NOT EXISTS idx_connected_platforms_contact ON connected_platforms(contact_id);
CREATE INDEX IF NOT EXISTS idx_connected_platforms_platform ON connected_platforms(platform);
CREATE INDEX IF NOT EXISTS idx_connected_platforms_agent ON connected_platforms(ai_agent_id);

CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);

CREATE INDEX IF NOT EXISTS idx_conversations_contact ON conversations(contact_id);
CREATE INDEX IF NOT EXISTS idx_conversations_platform ON conversations(platform);

CREATE INDEX IF NOT EXISTS idx_messages_conversation ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);

-- ========================================
-- End of Base Schema
-- ========================================