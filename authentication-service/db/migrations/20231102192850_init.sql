-- +goose Up
-- +goose StatementBegin

begin;

	create sequence public.user_id_seq
	start with 1
	increment by 1
	no minvalue
	no maxvalue 
	cashe 1;
	
	alter table public.user_id_seq owner to postgres;
	
	set default_tablespace = '';
	
	set default_table_access_method = heap;
	
	create table public.users (
		id integer DEFAULT nextval('public.user_id_seq'::regclass) NOT NULL,
		email character varying(255),
		first_name character varying(255),
		last_name character varying(255),
		password character varying(60),
		user_active integer DEFAULT 0,
		created_at timestamp without time zone,
		updated_at timestamp without time zone
	);
	
	alter table public.users OWNER TO postgres;
	
	select pg_catalog.setval('public.user_id_seq', 1, true);
	
	alter table only public.users
	add constraint users_pkey PRIMARY KEY (id);
	
	insert into "public"."users"("email","first_name","last_name","password","user_active","created_at","updated_at")
	values (E'admin@example.com',E'Admin',E'User',E'$2a$12$1zGLuYDDNvATh4RA4avbKuheAMpb1svexSzrQm7up.bnpwQHs0jNe',1,E'2022-03-14 00:00:00',E'2022-03-14 00:00:00');
	
	commit;

-- +goose StatementEnd
