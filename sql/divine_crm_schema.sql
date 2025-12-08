--
-- PostgreSQL database dump
--

\restrict xZiSob15ABqvCFDZZ3bNzhGbi7qBqQzczDO7IqDprN4PrIMFNUJcqJawA2X164o

-- Dumped from database version 18.0
-- Dumped by pg_dump version 18.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: ai_agents; Type: TABLE; Schema: public; Owner: divine_user
--

CREATE TABLE public.ai_agents (
    id bigint NOT NULL,
    name character varying(255),
    ai_engine character varying(100),
    basic_prompt text,
    active boolean DEFAULT true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.ai_agents OWNER TO divine_user;

--
-- Name: ai_agents_id_seq; Type: SEQUENCE; Schema: public; Owner: divine_user
--

CREATE SEQUENCE public.ai_agents_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ai_agents_id_seq OWNER TO divine_user;

--
-- Name: ai_agents_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: divine_user
--

ALTER SEQUENCE public.ai_agents_id_seq OWNED BY public.ai_agents.id;


--
-- Name: ai_configurations; Type: TABLE; Schema: public; Owner: divine_user
--

CREATE TABLE public.ai_configurations (
    id bigint NOT NULL,
    ai_engine character varying(100),
    token text,
    endpoint text,
    model character varying(100),
    active boolean DEFAULT true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.ai_configurations OWNER TO divine_user;

--
-- Name: ai_configurations_id_seq; Type: SEQUENCE; Schema: public; Owner: divine_user
--

CREATE SEQUENCE public.ai_configurations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ai_configurations_id_seq OWNER TO divine_user;

--
-- Name: ai_configurations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: divine_user
--

ALTER SEQUENCE public.ai_configurations_id_seq OWNED BY public.ai_configurations.id;


--
-- Name: broadcast_templates; Type: TABLE; Schema: public; Owner: divine_user
--

CREATE TABLE public.broadcast_templates (
    id bigint NOT NULL,
    name character varying(255),
    content text,
    type character varying(50),
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.broadcast_templates OWNER TO divine_user;

--
-- Name: broadcast_templates_id_seq; Type: SEQUENCE; Schema: public; Owner: divine_user
--

CREATE SEQUENCE public.broadcast_templates_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.broadcast_templates_id_seq OWNER TO divine_user;

--
-- Name: broadcast_templates_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: divine_user
--

ALTER SEQUENCE public.broadcast_templates_id_seq OWNED BY public.broadcast_templates.id;


--
-- Name: chat_labels; Type: TABLE; Schema: public; Owner: divine_user
--

CREATE TABLE public.chat_labels (
    id bigint NOT NULL,
    label character varying(100),
    description character varying(500),
    color character varying(50),
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.chat_labels OWNER TO divine_user;

--
-- Name: chat_labels_id_seq; Type: SEQUENCE; Schema: public; Owner: divine_user
--

CREATE SEQUENCE public.chat_labels_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.chat_labels_id_seq OWNER TO divine_user;

--
-- Name: chat_labels_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: divine_user
--

ALTER SEQUENCE public.chat_labels_id_seq OWNED BY public.chat_labels.id;


--
-- Name: chat_messages; Type: TABLE; Schema: public; Owner: divine_user
--

CREATE TABLE public.chat_messages (
    id bigint NOT NULL,
    contact_id bigint,
    contact_name character varying(255),
    message text,
    response text,
    status character varying(50) DEFAULT 'Unassigned'::character varying,
    assigned_to character varying(255),
    assigned_agent character varying(255),
    channel character varying(50),
    labels text,
    tokens_used bigint DEFAULT 0,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.chat_messages OWNER TO divine_user;

--
-- Name: chat_messages_id_seq; Type: SEQUENCE; Schema: public; Owner: divine_user
--

CREATE SEQUENCE public.chat_messages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.chat_messages_id_seq OWNER TO divine_user;

--
-- Name: chat_messages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: divine_user
--

ALTER SEQUENCE public.chat_messages_id_seq OWNED BY public.chat_messages.id;


--
-- Name: connected_platforms; Type: TABLE; Schema: public; Owner: divine_user
--

CREATE TABLE public.connected_platforms (
    id bigint NOT NULL,
    platform character varying(100),
    token text,
    client_id character varying(255),
    client_secret text,
    webhook_url text,
    phone_number_id character varying(255),
    active boolean DEFAULT true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.connected_platforms OWNER TO divine_user;

--
-- Name: connected_platforms_id_seq; Type: SEQUENCE; Schema: public; Owner: divine_user
--

CREATE SEQUENCE public.connected_platforms_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.connected_platforms_id_seq OWNER TO divine_user;

--
-- Name: connected_platforms_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: divine_user
--

ALTER SEQUENCE public.connected_platforms_id_seq OWNED BY public.connected_platforms.id;


--
-- Name: contacts; Type: TABLE; Schema: public; Owner: divine_user
--

CREATE TABLE public.contacts (
    id bigint NOT NULL,
    code character varying(50),
    channel character varying(50),
    channel_id character varying(255),
    name character varying(255),
    temperature character varying(20),
    first_contact timestamp with time zone,
    last_contact timestamp with time zone,
    last_agent character varying(255),
    last_agent_type character varying(50),
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.contacts OWNER TO divine_user;

--
-- Name: contacts_id_seq; Type: SEQUENCE; Schema: public; Owner: divine_user
--

CREATE SEQUENCE public.contacts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.contacts_id_seq OWNER TO divine_user;

--
-- Name: contacts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: divine_user
--

ALTER SEQUENCE public.contacts_id_seq OWNED BY public.contacts.id;


--
-- Name: leads; Type: TABLE; Schema: public; Owner: divine_user
--

CREATE TABLE public.leads (
    id bigint NOT NULL,
    code character varying(50),
    channel character varying(50),
    channel_id character varying(255),
    temperature character varying(20),
    first_contact timestamp with time zone,
    last_contact timestamp with time zone,
    last_agent character varying(255),
    last_agent_type character varying(50),
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.leads OWNER TO divine_user;

--
-- Name: leads_id_seq; Type: SEQUENCE; Schema: public; Owner: divine_user
--

CREATE SEQUENCE public.leads_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.leads_id_seq OWNER TO divine_user;

--
-- Name: leads_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: divine_user
--

ALTER SEQUENCE public.leads_id_seq OWNED BY public.leads.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: divine_user
--

CREATE TABLE public.products (
    id bigint NOT NULL,
    code character varying(50),
    name character varying(255),
    price numeric,
    stock bigint,
    uploaded_by character varying(255),
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.products OWNER TO divine_user;

--
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: divine_user
--

CREATE SEQUENCE public.products_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.products_id_seq OWNER TO divine_user;

--
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: divine_user
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- Name: quick_replies; Type: TABLE; Schema: public; Owner: divine_user
--

CREATE TABLE public.quick_replies (
    id bigint NOT NULL,
    trigger character varying(255),
    response text,
    active boolean DEFAULT true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.quick_replies OWNER TO divine_user;

--
-- Name: quick_replies_id_seq; Type: SEQUENCE; Schema: public; Owner: divine_user
--

CREATE SEQUENCE public.quick_replies_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.quick_replies_id_seq OWNER TO divine_user;

--
-- Name: quick_replies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: divine_user
--

ALTER SEQUENCE public.quick_replies_id_seq OWNED BY public.quick_replies.id;


--
-- Name: token_balances; Type: TABLE; Schema: public; Owner: divine_user
--

CREATE TABLE public.token_balances (
    id bigint NOT NULL,
    ai_engine character varying(100),
    total_tokens bigint,
    used_tokens bigint DEFAULT 0,
    remaining_tokens bigint,
    last_reset timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.token_balances OWNER TO divine_user;

--
-- Name: token_balances_id_seq; Type: SEQUENCE; Schema: public; Owner: divine_user
--

CREATE SEQUENCE public.token_balances_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.token_balances_id_seq OWNER TO divine_user;

--
-- Name: token_balances_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: divine_user
--

ALTER SEQUENCE public.token_balances_id_seq OWNED BY public.token_balances.id;


--
-- Name: ai_agents id; Type: DEFAULT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.ai_agents ALTER COLUMN id SET DEFAULT nextval('public.ai_agents_id_seq'::regclass);


--
-- Name: ai_configurations id; Type: DEFAULT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.ai_configurations ALTER COLUMN id SET DEFAULT nextval('public.ai_configurations_id_seq'::regclass);


--
-- Name: broadcast_templates id; Type: DEFAULT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.broadcast_templates ALTER COLUMN id SET DEFAULT nextval('public.broadcast_templates_id_seq'::regclass);


--
-- Name: chat_labels id; Type: DEFAULT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.chat_labels ALTER COLUMN id SET DEFAULT nextval('public.chat_labels_id_seq'::regclass);


--
-- Name: chat_messages id; Type: DEFAULT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.chat_messages ALTER COLUMN id SET DEFAULT nextval('public.chat_messages_id_seq'::regclass);


--
-- Name: connected_platforms id; Type: DEFAULT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.connected_platforms ALTER COLUMN id SET DEFAULT nextval('public.connected_platforms_id_seq'::regclass);


--
-- Name: contacts id; Type: DEFAULT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.contacts ALTER COLUMN id SET DEFAULT nextval('public.contacts_id_seq'::regclass);


--
-- Name: leads id; Type: DEFAULT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.leads ALTER COLUMN id SET DEFAULT nextval('public.leads_id_seq'::regclass);


--
-- Name: products id; Type: DEFAULT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- Name: quick_replies id; Type: DEFAULT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.quick_replies ALTER COLUMN id SET DEFAULT nextval('public.quick_replies_id_seq'::regclass);


--
-- Name: token_balances id; Type: DEFAULT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.token_balances ALTER COLUMN id SET DEFAULT nextval('public.token_balances_id_seq'::regclass);


--
-- Name: ai_agents ai_agents_pkey; Type: CONSTRAINT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.ai_agents
    ADD CONSTRAINT ai_agents_pkey PRIMARY KEY (id);


--
-- Name: ai_configurations ai_configurations_pkey; Type: CONSTRAINT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.ai_configurations
    ADD CONSTRAINT ai_configurations_pkey PRIMARY KEY (id);


--
-- Name: broadcast_templates broadcast_templates_pkey; Type: CONSTRAINT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.broadcast_templates
    ADD CONSTRAINT broadcast_templates_pkey PRIMARY KEY (id);


--
-- Name: chat_labels chat_labels_pkey; Type: CONSTRAINT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.chat_labels
    ADD CONSTRAINT chat_labels_pkey PRIMARY KEY (id);


--
-- Name: chat_messages chat_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.chat_messages
    ADD CONSTRAINT chat_messages_pkey PRIMARY KEY (id);


--
-- Name: connected_platforms connected_platforms_pkey; Type: CONSTRAINT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.connected_platforms
    ADD CONSTRAINT connected_platforms_pkey PRIMARY KEY (id);


--
-- Name: contacts contacts_pkey; Type: CONSTRAINT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.contacts
    ADD CONSTRAINT contacts_pkey PRIMARY KEY (id);


--
-- Name: leads leads_pkey; Type: CONSTRAINT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.leads
    ADD CONSTRAINT leads_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: quick_replies quick_replies_pkey; Type: CONSTRAINT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.quick_replies
    ADD CONSTRAINT quick_replies_pkey PRIMARY KEY (id);


--
-- Name: token_balances token_balances_pkey; Type: CONSTRAINT; Schema: public; Owner: divine_user
--

ALTER TABLE ONLY public.token_balances
    ADD CONSTRAINT token_balances_pkey PRIMARY KEY (id);


--
-- Name: idx_ai_agents_name; Type: INDEX; Schema: public; Owner: divine_user
--

CREATE UNIQUE INDEX idx_ai_agents_name ON public.ai_agents USING btree (name);


--
-- Name: idx_ai_configurations_ai_engine; Type: INDEX; Schema: public; Owner: divine_user
--

CREATE UNIQUE INDEX idx_ai_configurations_ai_engine ON public.ai_configurations USING btree (ai_engine);


--
-- Name: idx_channel_contact; Type: INDEX; Schema: public; Owner: divine_user
--

CREATE INDEX idx_channel_contact ON public.contacts USING btree (channel_id);


--
-- Name: idx_chat_messages_contact_id; Type: INDEX; Schema: public; Owner: divine_user
--

CREATE INDEX idx_chat_messages_contact_id ON public.chat_messages USING btree (contact_id);


--
-- Name: idx_chat_messages_created_at; Type: INDEX; Schema: public; Owner: divine_user
--

CREATE INDEX idx_chat_messages_created_at ON public.chat_messages USING btree (created_at);


--
-- Name: idx_chat_messages_status; Type: INDEX; Schema: public; Owner: divine_user
--

CREATE INDEX idx_chat_messages_status ON public.chat_messages USING btree (status);


--
-- Name: idx_connected_platforms_platform; Type: INDEX; Schema: public; Owner: divine_user
--

CREATE UNIQUE INDEX idx_connected_platforms_platform ON public.connected_platforms USING btree (platform);


--
-- Name: idx_contacts_code; Type: INDEX; Schema: public; Owner: divine_user
--

CREATE UNIQUE INDEX idx_contacts_code ON public.contacts USING btree (code);


--
-- Name: idx_contacts_last_contact; Type: INDEX; Schema: public; Owner: divine_user
--

CREATE INDEX idx_contacts_last_contact ON public.contacts USING btree (last_contact);


--
-- Name: idx_contacts_temperature; Type: INDEX; Schema: public; Owner: divine_user
--

CREATE INDEX idx_contacts_temperature ON public.contacts USING btree (temperature);


--
-- Name: idx_leads_code; Type: INDEX; Schema: public; Owner: divine_user
--

CREATE UNIQUE INDEX idx_leads_code ON public.leads USING btree (code);


--
-- Name: idx_products_code; Type: INDEX; Schema: public; Owner: divine_user
--

CREATE UNIQUE INDEX idx_products_code ON public.products USING btree (code);


--
-- Name: idx_quick_replies_trigger; Type: INDEX; Schema: public; Owner: divine_user
--

CREATE UNIQUE INDEX idx_quick_replies_trigger ON public.quick_replies USING btree (trigger);


--
-- Name: idx_token_balances_ai_engine; Type: INDEX; Schema: public; Owner: divine_user
--

CREATE UNIQUE INDEX idx_token_balances_ai_engine ON public.token_balances USING btree (ai_engine);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: pg_database_owner
--

GRANT ALL ON SCHEMA public TO divine_user;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON TABLES TO divine_user;


--
-- PostgreSQL database dump complete
--

\unrestrict xZiSob15ABqvCFDZZ3bNzhGbi7qBqQzczDO7IqDprN4PrIMFNUJcqJawA2X164o

