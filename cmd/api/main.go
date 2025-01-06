package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/mnabil1718/blog.mnabil.dev/internal/config"
	"github.com/mnabil1718/blog.mnabil.dev/internal/data"
	"github.com/mnabil1718/blog.mnabil.dev/internal/jsonlog"
	"github.com/mnabil1718/blog.mnabil.dev/internal/mailer"
	"github.com/mnabil1718/blog.mnabil.dev/internal/storage"
	"github.com/spf13/viper"
)

var (
	version string = "1.0.0"
)

type application struct {
	logger  *jsonlog.Logger
	config  config.Config
	models  data.Models
	wg      sync.WaitGroup
	mailer  mailer.Mailer
	storage storage.ImageStorage
}

func main() {

	var cfg config.Config

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		logger.PrintFatal(err, nil)
	}

	if err := config.LoadConfig(&cfg); err != nil {
		logger.PrintFatal(err, nil)
		os.Exit(1)
	}

	if viper.GetBool("DISPLAY_VERSION") {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	logger.PrintInfo("database connection pool established successfully.", nil)

	storage, err := storage.New(cfg.Upload.Path, cfg.Upload.TempPath)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	app := application{
		logger:  logger,
		config:  cfg,
		models:  data.NewModels(db),
		mailer:  mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender),
		storage: *storage,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

}
