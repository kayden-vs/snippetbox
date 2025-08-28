package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/form"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kayden-vs/snippetbox/internal/models"
)

type application struct {
	errorLog    *log.Logger
	infoLog     *log.Logger
	snippets    *models.SnippetModel
	formDecoder *form.Decoder
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP server port")
	dsn := flag.String("dsn", "rohit:eren@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	formDecoder := form.NewDecoder()

	app := &application{
		errorLog:    errorLog,
		infoLog:     infoLog,
		snippets:    &models.SnippetModel{DB: db},
		formDecoder: formDecoder,
	}

	infoLog.Printf("Starting server on %s", *addr)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
