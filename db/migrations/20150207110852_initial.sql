
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE users (
  id integer not null,
  email character varying(250) not null,
  token character varying(50) not null,
  ttl timestamp with time zone not null,
  originurl character varying(250)
);

-- Sequences
CREATE SEQUENCE users_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;


ALTER SEQUENCE users_id_seq OWNED BY users.id;

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);

-- Constraints
ALTER TABLE ONLY users
ADD CONSTRAINT users_pkey PRIMARY KEY (id),
ADD CONSTRAINT users_email_uniq UNIQUE (email);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE users CASCADE;
