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
) (Book, error) {
	query := `SELECT name, author, genre FROM book
	WHERE name = $1 LIMIT 1`

	
	var span opentracing.Span
	if cspan, ok := ctx.Get("tracing-context"); ok {
		span = ginopentracing.StartDBSpanWithParent(cspan.(opentracing.Span).Context(), "getBook", "posgres", "sql", "select")	
	} else {
		span = ginopentracing.StartSpanWithHeader(&ctx.Request.Header, "getBook", ctx.Request.Method, ctx.Request.URL.Path)
	}
	defer span.Finish()

	
	rows := r.db.QueryRowContext(ctx, query, genre)

	var i Book
	err := rows.Scan(
		&i.Name,
		&i.Author,
		&i.Genre,
	)
	return i, err
}

// Close calls close on the underlying db connection.
func (r *Repository) Close() {
	r.db.Close()
}
