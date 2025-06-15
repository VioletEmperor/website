DROP SCHEMA IF EXISTS public CASCADE;

CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE public.posts (
    id serial primary key,
    title varchar(80),
    author varchar(80) default 'Adam Shkolnik',
    created timestamp default CURRENT_TIMESTAMP,
    edited timestamp default  CURRENT_TIMESTAMP,
    body varchar(80)
);

INSERT INTO public.posts VALUES
(DEFAULT, 'POST A', DEFAULT, DEFAULT, DEFAULT, 'a.html'),
(DEFAULT, 'POST B', DEFAULT, DEFAULT, DEFAULT, 'a.html');