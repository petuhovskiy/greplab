package greplab

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

func (l *Lab) ListProjects() error {
	const pageSize = 20

	for page := 1; ; page++ {
		projects, _, err := l.cli.Projects.ListProjects(&gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    page,
				PerPage: pageSize,
			},
			Visibility: gitlab.Visibility(gitlab.PrivateVisibility),
		})
		if err != nil {
			return err
		}

		for _, project := range projects {
			err := l.ProcessProject(project)
			if err != nil {
				log.WithError(err).Errorf("failed to process project %s", project.PathWithNamespace)
			}
		}

		if len(projects) < pageSize {
			break
		}
	}

	return nil
}

func (l *Lab) ProcessProject(project *gitlab.Project) error {
	fmt.Printf("Got project: %s\n", project.Name)
	return l.ListJobs(project)
}
