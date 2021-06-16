package database

import (
	"context"
	"net/url"
	"reflect"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	postgreConn *sqlx.DB
	postgreOnce sync.Once
)

type PostgreConfig struct {
	User         string
	Password     string
	Host         string
	Name         string
	MaxIdleConns int
	MaxOpenConns int
	DisableTLS   bool
}

// NewPostgreSQL returns connection to postgreSQL.
func NewPostgreSQL(cfg PostgreConfig) (*sqlx.DB, error) {
	var connErr error

	postgreOnce.Do(func() {

		q := make(url.Values)
		q.Set("sslmode", "disable")

		u := url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword(cfg.User, cfg.Password),
			Host:     cfg.Host,
			RawQuery: q.Encode(),
		}

		conn, err := sqlx.Connect("postgres", u.String())
		if err != nil {
			connErr = errors.Wrap(err, "db")
			return
		}

		conn.SetMaxIdleConns(cfg.MaxIdleConns)
		conn.SetMaxOpenConns(cfg.MaxOpenConns)

		postgreConn = conn
	})
	return postgreConn, connErr
}

// QueryStruct queries a single record and puts it into a struct.
func QueryStruct(ctx context.Context, db *sqlx.DB, query string, data interface{}, dest interface{}) error {
	row, err := db.NamedQueryContext(ctx, query, data)
	if err != nil {
		return err
	}
	if !row.Next() {
		return ErrNotFound
	}

	if err := row.StructScan(dest); err != nil {
		return err
	}

	return nil
}

// QuerySlice queries multiple records and puts them into a slice.
func QuerySlice(ctx context.Context, db *sqlx.DB, query string, data interface{}, dest interface{}) error {
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		return errors.New("must provide a pointer to a slice")
	}

	rows, err := db.NamedQueryContext(ctx, query, data)
	if err != nil {
		return err
	}

	slice := val.Elem()
	for rows.Next() {
		v := reflect.New(slice.Type().Elem())
		if err := rows.StructScan(v.Interface()); err != nil {
			return err
		}
		slice.Set(reflect.Append(slice, v.Elem()))
	}

	return nil
}
