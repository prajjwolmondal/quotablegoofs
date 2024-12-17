package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JokeType string

const (
	OneLiner   JokeType = "oneLiner"
	MultiLine  JokeType = "multiLiner"
	KnockKnock JokeType = "knockKnock"
)

func (j JokeType) IsValid() bool {
	switch j {
	case OneLiner, MultiLine, KnockKnock:
		return true
	default:
		return false
	}
}

type Joke struct {
	Id        int       `json:"id"`
	JokeType  JokeType  `json:"joke_type"`
	Content   []string  `json:"content"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (joke *Joke) Validate() []error {
	var errs []error

	if joke.Id != 0 {
		errs = append(errs, errors.New("field 'id' is not allowed in the request body"))
	}

	if !joke.CreatedAt.IsZero() {
		errs = append(errs, errors.New("field 'createdAt' is not allowed in the request body"))
	}

	if !joke.UpdatedAt.IsZero() {
		errs = append(errs, errors.New("field 'updatedAt' is not allowed in the request body"))
	}

	if len(joke.Content) == 0 {
		errs = append(errs, errors.New("field 'content' cannot be empty"))
	}

	if !joke.JokeType.IsValid() {
		errs = append(errs, errors.New("invalid 'joke_type' provided. Must be one of [oneLiner, multiLiner, knockKnock]"))
	}

	return errs
}

type JokeModel struct {
	DbPool *pgxpool.Pool
}

func (j *JokeModel) Insert(joke Joke) (Joke, error) {
	sqlStatement := `INSERT INTO jokes(joke_type, content, source, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id, joke_type, content, source, created_at, updated_at`

	var jokeFromDb Joke
	err := j.DbPool.QueryRow(context.Background(), sqlStatement, joke.JokeType, joke.Content, joke.Source, joke.CreatedAt, joke.UpdatedAt).Scan(&jokeFromDb.Id, &jokeFromDb.JokeType, &jokeFromDb.Content, &jokeFromDb.Source, &jokeFromDb.CreatedAt, &jokeFromDb.UpdatedAt)
	if err != nil {
		return Joke{}, err
	}

	return jokeFromDb, nil
}

func (j *JokeModel) Get(id int) (Joke, error) {
	sqlStatement := `SELECT id, joke_type, content, source, created_at, updated_at
	FROM jokes
	WHERE id = $1`

	var joke Joke
	err := j.DbPool.QueryRow(context.Background(), sqlStatement, id).Scan(&joke.Id, &joke.JokeType, &joke.Content, &joke.Source, &joke.CreatedAt, &joke.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Joke{}, ErrNoRecord
		} else {
			return Joke{}, err
		}
	}

	return joke, nil
}

func (j *JokeModel) GetRandomJokes(limit int) ([]Joke, error) {
	sqlStatment := `SELECT * FROM jokes ORDER BY random() LIMIT $1;`

	var jokePointers []*Joke

	err := pgxscan.Select(context.Background(), j.DbPool, &jokePointers, sqlStatment, limit)
	if err != nil {
		return []Joke{}, err
	}

	jokes := make([]Joke, len(jokePointers))
	for i, jp := range jokePointers {
		jokes[i] = *jp
	}

	return jokes, nil
}
