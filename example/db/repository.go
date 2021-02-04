package db

import (
	// "context"
	"database/sql"
	"log"
	"os"

	ginopentracing "github.com/Bose/go-gin-opentracing"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" //need this here
	"github.com/opentracing/opentracing-go"
)

var dbDriver = os.Getenv("DB_DRIVER")
var dbSource = os.Getenv("DB_SOURCE")

// Repository retrieves information about people.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new Repository backed by MySQL database.
func NewRepository() *Repository {
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Cannot ping the db: %v", err)
	}
	return &Repository{
		db: db,
	}
}

// GetBooks returns the list of books that match the given genre
func (r *Repository) GetBooks(
	ctx *gin.Context,
	genre string,
) ([]Book, error) {
	query := `SELECT name, author, genre FROM book
	WHERE genre = $1`

	var span opentracing.Span
	if cspan, ok := ctx.Get("tracing-context"); ok {
		span = ginopentracing.StartDBSpanWithParent(cspan.(opentracing.Span).Context(), "returnBook", "postgres", "sql", "select")
	} else {
		span = ginopentracing.StartSpanWithHeader(&ctx.Request.Header, "returnBook", ctx.Request.Method, ctx.Request.URL.Path)
	}
	defer span.Finish()

	rows, err := r.db.QueryContext(ctx, query, genre)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	books := []Book{}
	for rows.Next() {
		var i Book
		if err := rows.Scan(
			&i.Name,
			&i.Author,
			&i.Genre,
		); err != nil {
			return nil, err
		}
		books = append(books, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

// Close calls close on the underlying db connection.
func (r *Repository) Close() {
	r.db.Close()
}
