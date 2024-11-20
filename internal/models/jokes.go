package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
	Content   any       `json:"content"`
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

	if joke.Content == "" {
		errs = append(errs, errors.New("field 'content' cannot be empty"))
	}

	if joke.Source == "" {
		errs = append(errs, errors.New("field 'source' cannot be empty, if source isn't known then anonymous is valid"))
	}

	if !joke.JokeType.IsValid() {
		errs = append(errs, errors.New("invalid 'joke_type' provided. Must be one of [oneLiner, multiLiner, knockKnock]"))
	}

	return errs
}

type JokeModel struct {
	Dbpool *pgxpool.Pool
}

func (j *JokeModel) Insert(joke Joke) (int, error) {
	sqlStatement := `INSERT INTO jokes(joke_type, content, source, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id int
	err := j.Dbpool.QueryRow(context.Background(), sqlStatement, joke.JokeType, joke.Content, joke.Source, joke.CreatedAt, joke.UpdatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (j *JokeModel) Get(id int) (Joke, error) {
	sqlStatement := `SELECT id, joke_type, content, source, created_at, updated_at
	FROM jokes
	WHERE id = $1`

	var joke Joke
	err := j.Dbpool.QueryRow(context.Background(), sqlStatement, id).Scan(&joke.Id, &joke.JokeType, &joke.Content, &joke.Source, &joke.CreatedAt, &joke.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Joke{}, ErrNoRecord
		} else {
			return Joke{}, err
		}
	}

	return joke, nil
}
