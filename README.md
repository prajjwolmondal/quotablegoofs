# quotablegoofs

This is a Go webserver that will serve as a backend from a couple of mobile apps that I'm developing. My aim is not to use any frameworks and instead rely on native go modules for the server functionality. 

Currently these are the goals of this project:
1. Provide jokes
    - One-line jokes
    - Multi-line jokes
    - Knock knock jokes
2. Provide quotes
3. Provide end-users a way to upvote/downvote the jokes & quotes so that unpopular ones aren't shown as frequently

## Tools used
- Go/golang (v1.23.2)
- [air-verse/air](https://github.com/air-verse/air) - provide live reload functionality while developing locally
- PostgreSQL
	- Used `psql` and [DBeaver](https://dbeaver.com/) for DB commands and GUI.
- [jackc/pgx](https://github.com/jackc/pgx) - PostgreSQL driver for Go. 
	- This project utilizes the [pgxpool](https://pkg.go.dev/github.com/jackc/pgx/v5@v5.7.1/pgxpool) package to have a concuurency safe connections to the DB via a connection pool.
- Docker 
- Google Cloud Platform
	- Artifact Registry - holds the docker images
	- Secret Manager - holds sensitive information required by the application (e.g. DB credentials)
	- Cloud SQL - hosts the PostgreSQL server
	- Cloud Run - deploys the server in a container using the images in the Artifact Registry.
- [REST Client VS code extension](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) - for manual testing of the API (see `api_test.txt` for examples)
- [`go vet`](https://golang.google.cn/cmd/vet/), [`staticcheck`](https://staticcheck.dev/), and [`golangci-lint`](https://golangci-lint.run/) - Used to analyse the codebase. 
- [`govulncheck`](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) - to check for vulnerabilities.

## Architecture

![Architecture diagram of quotablegoofs](https://docs.google.com/drawings/d/e/2PACX-1vStlM5h46sGZBBEFQk08ugp1uL74L3WAXiVg6iF6OcFhctIKk2EYvU0N2w9YIbT11jQQFsgd6GQyoTW/pub?w=1288&h=688)

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