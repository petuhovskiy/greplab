package greplab

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

func (l *Lab) ListJobs(project *gitlab.Project) error {
	jobs, _, err := l.cli.Jobs.ListProjectJobs(project.ID, &gitlab.ListJobsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 100,
		},
		IncludeRetried: true,
	})
	if err != nil {
		return err
	}

	for _, job := range jobs {
		err := l.ProcessJob(project, job)
		if err != nil {
			log.WithError(err).Errorf("Failed to process job %d", job.ID)
		}
	}

	return nil
}

func (l *Lab) ProcessJob(project *gitlab.Project, job *gitlab.Job) error {
	fmt.Printf("Got job: %v, url=%v\n", job.Name, job.WebURL)

	if job.Name != l.config.JobNameFilter {
		return nil
	}

	logs, _, err := l.cli.Jobs.GetTraceFile(project.ID, job.ID)
	if err != nil {
		return err
	}

	logsBytes := make([]byte, logs.Size())
	_, err = logs.Read(logsBytes)

	grepBytes := []byte(l.config.GrepLine)

	var lines []string
	for len(logsBytes) > 0 {
		lineEnd := bytes.IndexByte(logsBytes, '\n')
		if lineEnd == -1 {
			lineEnd = len(logsBytes)
		}

		line := logsBytes[:lineEnd]
		index := bytes.Index(line, grepBytes)
		if index != -1 {
			lines = append(lines, string(line))
		}

		logsBytes = logsBytes[lineEnd+1:]
	}

	if len(lines) == 0 {
		return nil
	}

	return l.HandleLines(project, job, lines)
}
