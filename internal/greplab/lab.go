package greplab

import (
	"github.com/petuhovskiy/greplab/internal/conf"
	"github.com/xanzy/go-gitlab"
	"gorm.io/gorm"
)

type Lab struct {
	cli    *gitlab.Client
	config *conf.App
	db     *gorm.DB
}

func NewLab(cli *gitlab.Client, config *conf.App, db *gorm.DB) *Lab {
	return &Lab{
		cli:    cli,
		config: config,
		db:     db,
	}
}
