package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"lincoln.boris/forum/pkg/models/sqlite3"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	categories *sqlite3.CategoryModel
	comments *sqlite3.CommentModel
	comment_votes *sqlite3.CommentVoteModel
	errorLog *log.Logger
	infoLog *log.Logger
	posts *sqlite3.PostModel
	post_votes *sqlite3.PostVoteModel
	post_categories *sqlite3.PostCategoryModel
	sessions *sqlite3.SessionModel
	templateCache map[string]*template.Template
	users *sqlite3.UserModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "mydb.db", "SQLite3 database connection")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		categories: &sqlite3.CategoryModel{DB: db},
		comments: &sqlite3.CommentModel{DB: db},
		comment_votes: &sqlite3.CommentVoteModel{DB: db},
		errorLog: errorLog,
		infoLog: infoLog,
		posts: &sqlite3.PostModel{DB: db},
		post_categories: &sqlite3.PostCategoryModel{DB: db},
		post_votes: &sqlite3.PostVoteModel{DB: db},
		sessions: &sqlite3.SessionModel{DB: db},
		templateCache: templateCache,
		users: &sqlite3.UserModel{DB: db},
	}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}
