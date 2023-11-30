package changelog

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upInit, downInit)
}

func upInit(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create schema service;

		create sequence service.user_id_seq
		start with 1
		increment by 1
		no minvalue
		no maxvalue
		cache 1;

		alter table service.user_id_seq owner to postgres;

		set default_tablespace = '';

		set default_table_access_method = heap;

		create table service.users (
			id integer default nextval('service.user_id_seq'::regclass) NOT NULL,
			email character varying(255),
			first_name character varying(255),
			last_name character varying(255),
			password character varying(60),
			user_active integer default 0,
			created_at timestamp without time zone,
			updated_at timestamp without time zone
		);

		alter table service.users OWNER TO postgres;

		select pg_catalog.setval('service.user_id_seq', 1, true);

		alter table only service.users
		add constraint users_pkey PRIMARY KEY (id);

		insert into "service"."users"("email","first_name","last_name","password","user_active","created_at","updated_at")
		values (E'admin@example.com',E'Admin',E'User',E'$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK',1,E'2022-03-14 00:00:00',E'2022-03-14 00:00:00');
     `)
	if err != nil {
		return err
	}

	return nil
}

func downInit(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
