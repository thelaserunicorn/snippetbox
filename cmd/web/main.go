package main

import (
    "database/sql"
  "crypto/tls"
    "flag"
    "log/slog"
    "net/http"
    "os"
  "time"
    "html/template"
      "github.com/alexedwards/scs/mysqlstore" // New import
    "github.com/alexedwards/scs/v2"
"github.com/go-playground/form/v4"
    "thelaserunicorn.snippetbox/internal/models" 

    _ "github.com/go-sql-driver/mysql"

)

type application struct {
    logger        *slog.Logger
    snippets       *models.SnippetModel
    users          *models.UserModel
    templateCache  map[string]*template.Template
    formDecoder    *form.Decoder
    sessionManager *scs.SessionManager
}

func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")
    dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
    flag.Parse()

    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    db, err := openDB(*dsn)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }
    defer db.Close()

    templateCache, err := newTemplateCache()
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }
    formDecoder := form.NewDecoder()

    // Use the scs.New() function to initialize a new session manager. Then we
    // configure it to use our MySQL database as the session store, and set a
    // lifetime of 12 hours (so that sessions automatically expire 12 hours
    // after first being created).
    sessionManager := scs.New()
    sessionManager.Store = mysqlstore.New(db)
    sessionManager.Lifetime = 12 * time.Hour
   sessionManager.Cookie.Secure = true


    // And add the session manager to our application dependencies.
    app := &application{
        logger:         logger,
        snippets:       &models.SnippetModel{DB: db},
        users:          &models.UserModel{DB: db},
        templateCache:  templateCache,
        formDecoder:    formDecoder,
        sessionManager: sessionManager,
    }

    tlsConfig := &tls.Config{
        CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
    }
      srv := &http.Server{
        Addr:    *addr,
        Handler: app.routes(),
    ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
    TLSConfig: tlsConfig,
            IdleTimeout:  time.Minute,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    logger.Info("starting server", "addr", srv.Addr)

    // Call the ListenAndServe() method on our new http.Server struct to start 
    // the server.
  err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
