package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Quote struct {
	Id        int       `json:"id"`
	Content   []string  `json:"content"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (quote *Quote) Validate() []error {
	var errs []error

	if quote.Id != 0 {
		errs = append(errs, errors.New("field 'id' is not allowed in the request body"))
	}

	if !quote.CreatedAt.IsZero() {
		errs = append(errs, errors.New("field 'createdAt' is not allowed in the request body"))
	}

	if !quote.UpdatedAt.IsZero() {
		errs = append(errs, errors.New("field 'updatedAt' is not allowed in the request body"))
	}

	if len(quote.Content) == 0 {
		errs = append(errs, errors.New("field 'content' cannot be empty"))
	}

	return errs
}

type QuoteModel struct {
	DbPool *pgxpool.Pool
}

func (q *QuoteModel) Insert(quote Quote) (Quote, error) {
	sqlStatement := `INSERT INTO quotes(content, source, created_at, updated_at) 
	VALUES ($1, $2, $3, $4) RETURNING id`

	var id int
	err := q.DbPool.QueryRow(context.Background(), sqlStatement, quote.Content, quote.Source, quote.CreatedAt, quote.UpdatedAt).Scan(&id)
	if err != nil {
		return Quote{}, err
	}

	return Quote{
		Id:        id,
		Content:   quote.Content,
		Source:    quote.Source,
		CreatedAt: quote.CreatedAt,
		UpdatedAt: quote.UpdatedAt,
	}, nil
}

func (q *QuoteModel) Get(id int) (Quote, error) {
	sqlStatement := `SELECT id, content, source, created_at, updated_at
	FROM quotes
	WHERE id = $1`

	var quote Quote
	err := q.DbPool.QueryRow(context.Background(), sqlStatement, id).Scan(&quote.Id, &quote.Content, &quote.Source, &quote.CreatedAt, &quote.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Quote{}, ErrNoRecord
		} else {
			return Quote{}, err
		}
	}

	return quote, nil
}

func (q *QuoteModel) GetRandomQuotes(limit int) ([]Quote, error) {
	sqlStatment := `SELECT * FROM quotes ORDER BY random() LIMIT $1;`

	var quotePointers []*Quote

	err := pgxscan.Select(context.Background(), q.DbPool, &quotePointers, sqlStatment, limit)
	if err != nil {
		return []Quote{}, err
	}

	quotes := make([]Quote, len(quotePointers))
	for i, qp := range quotePointers {
		quotes[i] = *qp
	}

	return quotes, nil
}
