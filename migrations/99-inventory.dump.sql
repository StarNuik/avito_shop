--
-- PostgreSQL database dump
--

-- Dumped from database version 13.19 (Debian 13.19-1.pgdg120+1)
-- Dumped by pg_dump version 16.6 (Ubuntu 16.6-0ubuntu0.24.04.1)

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

--
-- Data for Name: inventory; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.inventory (id, name, price) FROM stdin;
195	t-shirt	80
196	cup	20
197	book	50
198	pen	10
199	powerbank	200
200	hoody	300
201	umbrella	200
202	socks	10
203	wallet	50
204	pink-hoody	500
\.


--
-- Name: inventory_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.inventory_id_seq', 204, true);


--
-- PostgreSQL database dump complete
--

