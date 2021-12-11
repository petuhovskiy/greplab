package conf

import (
	"github.com/caarlos0/env/v6"
)

type App struct {
	PrometheusBind string `env:"PROMETHEUS_BIND" envDefault:":2112"`
	GitlabToken    string `env:"GITLAB_TOKEN"`
	DatabaseConn   string `env:"DATABASE_CONN"` // rename to PostgresDSN
	JobNameFilter  string `env:"JOB_NAME_FILTER"`
	GrepLine       string `env:"GREP_LINE"`
}

func ParseEnv() (*App, error) {
	cfg := App{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
