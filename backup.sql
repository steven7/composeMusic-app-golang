--
-- PostgreSQL database dump
--

-- Dumped from database version 14.1 (Debian 14.1-1.pgdg110+1)
-- Dumped by pg_dump version 14.1 (Debian 14.1-1.pgdg110+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
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
-- Name: galleries; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.galleries (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id integer,
    title text
);


ALTER TABLE public.galleries OWNER TO postgres;

--
-- Name: galleries_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.galleries_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.galleries_id_seq OWNER TO postgres;

--
-- Name: galleries_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.galleries_id_seq OWNED BY public.galleries.id;


--
-- Name: o_auths; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.o_auths (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id integer NOT NULL,
    service text NOT NULL,
    access_token text,
    token_type text,
    refresh_token text,
    expiry timestamp with time zone
);


ALTER TABLE public.o_auths OWNER TO postgres;

--
-- Name: o_auths_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.o_auths_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.o_auths_id_seq OWNER TO postgres;

--
-- Name: o_auths_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.o_auths_id_seq OWNED BY public.o_auths.id;


--
-- Name: pw_resets; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pw_resets (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id integer NOT NULL,
    token_hash text NOT NULL
);


ALTER TABLE public.pw_resets OWNER TO postgres;

--
-- Name: pw_resets_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pw_resets_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.pw_resets_id_seq OWNER TO postgres;

--
-- Name: pw_resets_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pw_resets_id_seq OWNED BY public.pw_resets.id;


--
-- Name: tracks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tracks (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id integer,
    title text,
    description text,
    filename text
);


ALTER TABLE public.tracks OWNER TO postgres;

--
-- Name: tracks_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.tracks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tracks_id_seq OWNER TO postgres;

--
-- Name: tracks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tracks_id_seq OWNED BY public.tracks.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    email text NOT NULL,
    password_hash text NOT NULL,
    remember_hash text NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: galleries id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.galleries ALTER COLUMN id SET DEFAULT nextval('public.galleries_id_seq'::regclass);


--
-- Name: o_auths id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.o_auths ALTER COLUMN id SET DEFAULT nextval('public.o_auths_id_seq'::regclass);


--
-- Name: pw_resets id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pw_resets ALTER COLUMN id SET DEFAULT nextval('public.pw_resets_id_seq'::regclass);


--
-- Name: tracks id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tracks ALTER COLUMN id SET DEFAULT nextval('public.tracks_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: galleries; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.galleries (id, created_at, updated_at, deleted_at, user_id, title) FROM stdin;
\.


--
-- Data for Name: o_auths; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.o_auths (id, created_at, updated_at, deleted_at, user_id, service, access_token, token_type, refresh_token, expiry) FROM stdin;
\.


--
-- Data for Name: pw_resets; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.pw_resets (id, created_at, updated_at, deleted_at, user_id, token_hash) FROM stdin;
\.


--
-- Data for Name: tracks; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tracks (id, created_at, updated_at, deleted_at, user_id, title, description, filename) FROM stdin;
1	2022-01-20 21:39:06.765082+00	2022-01-20 21:39:06.77739+00	\N	1	test compseWithWebApp postman	test compseWithWebApp postman	
2	2022-01-21 00:31:15.739412+00	2022-01-21 00:31:15.75787+00	\N	1	test compseWithWebApp postman	test compseWithWebApp postman	
35	2022-01-21 00:45:59.316201+00	2022-01-21 00:45:59.33749+00	\N	1	test compseWithWebApp postman	test compseWithWebApp postman	
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, created_at, updated_at, deleted_at, email, password_hash, remember_hash) FROM stdin;
\.


--
-- Name: galleries_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.galleries_id_seq', 1, false);


--
-- Name: o_auths_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.o_auths_id_seq', 1, false);


--
-- Name: pw_resets_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.pw_resets_id_seq', 1, false);


--
-- Name: tracks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.tracks_id_seq', 35, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 1, false);


--
-- Name: galleries galleries_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.galleries
    ADD CONSTRAINT galleries_pkey PRIMARY KEY (id);


--
-- Name: o_auths o_auths_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.o_auths
    ADD CONSTRAINT o_auths_pkey PRIMARY KEY (id);


--
-- Name: pw_resets pw_resets_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pw_resets
    ADD CONSTRAINT pw_resets_pkey PRIMARY KEY (id);


--
-- Name: tracks tracks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tracks
    ADD CONSTRAINT tracks_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_galleries_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_galleries_deleted_at ON public.galleries USING btree (deleted_at);


--
-- Name: idx_galleries_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_galleries_user_id ON public.galleries USING btree (user_id);


--
-- Name: idx_o_auths_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_o_auths_deleted_at ON public.o_auths USING btree (deleted_at);


--
-- Name: idx_pw_resets_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_pw_resets_deleted_at ON public.pw_resets USING btree (deleted_at);


--
-- Name: idx_tracks_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_tracks_deleted_at ON public.tracks USING btree (deleted_at);


--
-- Name: idx_tracks_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_tracks_user_id ON public.tracks USING btree (user_id);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: uix_pw_resets_token_hash; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX uix_pw_resets_token_hash ON public.pw_resets USING btree (token_hash);


--
-- Name: uix_users_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX uix_users_email ON public.users USING btree (email);


--
-- Name: uix_users_remember_hash; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX uix_users_remember_hash ON public.users USING btree (remember_hash);


--
-- Name: user_id_service; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_id_service ON public.o_auths USING btree (user_id, service);


--
-- PostgreSQL database dump complete
--

