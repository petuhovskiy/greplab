package greplab

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"strconv"
	"strings"
	"time"
)

type Grade struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	ProjectID int `gorm:"primaryKey"`
	JobID     int `gorm:"primaryKey"`
	JobURL    string

	Author        string
	AuthorEmail   string
	CommitMessage string
	CommitTime    *time.Time

	FinalGrade int
}

func (l *Lab) HandleLines(project *gitlab.Project, job *gitlab.Job, lines []string) error {
	fmt.Printf("Found lines: %s\n", lines)
	finalGrade := 0

	for _, line := range lines {
		line = strings.SplitN(line, l.config.GrepLine, 2)[1]
		line = strings.TrimSpace(line)
		grade, err := strconv.Atoi(line)
		if err != nil {
			log.Errorf("failed to parse grade in string")
			return err
		}

		finalGrade += grade
	}

	grade := Grade{
		ProjectID:     project.ID,
		JobID:         job.ID,
		JobURL:        job.WebURL,
		Author:        job.Commit.AuthorName,
		AuthorEmail:   job.Commit.AuthorEmail,
		CommitMessage: job.Commit.Message,
		CommitTime:    job.Commit.CreatedAt,
		FinalGrade:    finalGrade,
	}

	return l.db.Create(&grade).Error
}
