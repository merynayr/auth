package pg

import (
	"context"
	"fmt"
	"log"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/merynayr/auth/internal/client/db"
	"github.com/merynayr/auth/internal/client/db/prettier"
)

// pg структура с подключением к Postgres
type pg struct {
	dbc *pgxpool.Pool
}

// NewDB возвращает объект с подключением к Postgres
func NewDB(dbc *pgxpool.Pool) db.DB {
	return &pg{
		dbc: dbc,
	}
}

// ScanOneContext обертка над методом
func (p *pg) ScanOneContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	logQuery(ctx, q, args...)

	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

// ScanAllContext обертка над методом
func (p *pg) ScanAllContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	logQuery(ctx, q, args...)

	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, row)
}

// ExecContext обертка над методом
func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...interface{}) (pgconn.CommandTag, error) {
	logQuery(ctx, q, args...)

	return p.dbc.Exec(ctx, q.QueryRaw, args...)
}

// QueryContext обертка над методом
func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...interface{}) (pgx.Rows, error) {
	logQuery(ctx, q, args...)

	return p.dbc.Query(ctx, q.QueryRaw, args...)
}

// QueryRowContext обертка над методом
func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row {
	logQuery(ctx, q, args...)

	return p.dbc.QueryRow(ctx, q.QueryRaw, args...)
}

// Ping обертка над методом
func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

// Close обертка над методом
func (p *pg) Close() {
	p.dbc.Close()
}

func logQuery(ctx context.Context, q db.Query, args ...interface{}) {
	prettyQuery := prettier.Pretty(q.QueryRaw, prettier.PlaceholderDollar, args...)
	log.Println(
		ctx,
		fmt.Sprintf("sql: %s", q.Name),
		fmt.Sprintf("query: %s", prettyQuery),
	)
}
