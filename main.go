package main

import (
	"github.com/petuhovskiy/greplab/internal/greplab"
	"github.com/xanzy/go-gitlab"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"

	"github.com/petuhovskiy/greplab/internal/conf"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	cfg, err := conf.ParseEnv()
	if err != nil {
		log.WithError(err).Fatal("failed to parse config from env")
	}

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(cfg.PrometheusBind, mux)
		if err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("prometheus server error")
		}
	}()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseConn))
	if err != nil {
		log.WithError(err).Fatal("failed to connect to database")
	}

	// automigrate all the tables
	err = db.AutoMigrate(&greplab.Grade{})
	if err != nil {
		log.WithError(err).Fatal("failed to automigrate database")
	}

	cli, err := gitlab.NewClient(cfg.GitlabToken)
	if err != nil {
		log.WithError(err).Fatal("failed to create gitlab client")
	}

	lab := greplab.NewLab(cli, cfg, db)

	err = lab.ListProjects()
	if err != nil {
		log.WithError(err).Fatal("failed to list projects")
	}
}
