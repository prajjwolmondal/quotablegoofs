# quotablegoofs

This is a Go webserver that will server as a backend from a couple of mobile apps that I'm developing.

Currently these are the goals of this project:
1. Provide jokes
    - One-line jokes
    - Multi-line jokes
    - Knock knock jokes
2. Provide quotes
3. Provide end-users a way to upvote/downvote the jokes & quotes so that unpopular ones aren't shown as frequently

## Tools used
- Go
- PostgreSQL
- [air-verse/air](https://github.com/air-verse/air) - provide live reload functionality while developing locally
- [jackc/pgx](https://github.com/jackc/pgx) - PostgreSQL driver

## Setup

DB setup (using psql):

```postgresql
CREATE DATABASE quotablegoofs;

CREATE USER quotablegoof WITH PASSWORD 'localdevpassword';

GRANT ALL PRIVILEGES ON DATABASE quotablegoofs TO quotablegoof;

\c quotablegoofs

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO quotablegoof;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO quotablegoof;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO quotablegoof;

ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO quotablegoof;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO quotablegoof;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON FUNCTIONS TO quotablegoof;
```

Table setup:
```postgresql
create table jokes (
	id SERIAL NOT NULL PRIMARY KEY ,
	joke_type VARCHAR(100) NOT NULL,
	content jsonb NOT NULL,
	source TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT null
);

create table quotes (
	id SERIAL NOT NULL PRIMARY KEY ,
	content jsonb NOT NULL,
	source TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT null
);
```