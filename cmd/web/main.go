package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql" // _ is used, if we are not using the import directly
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"snippetbox.amsakib.com/internal/models"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

type application struct {
	logger         *slog.Logger
	cfg            config
	snippetService models.SnippetService
	templateCache  map[string]*template.Template
}

func main() {

	// define flag
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "http service address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "static dir")
	flag.StringVar(&cfg.dsn, "dsn", "root:@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// add logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// automatically close the database connection
	defer db.Close()
	//snippetService, err := models.NewSnippetService(db)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	//defer snippetService.InsertStatement.Close()
	//defer snippetService.GetStatement.Close()
	//defer snippetService.LatestStatement.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:         logger,
		cfg:            cfg,
		snippetService: models.SnippetService{DB: db},
		templateCache:  templateCache,
	}

	logger.Info("Starting server on", slog.Any(cfg.addr, ":4000"))

	err = http.ListenAndServe(cfg.addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}