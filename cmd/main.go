package main

import (
	"github-issue-schedule/internal/pkg/configs"
	"github-issue-schedule/internal/pkg/github"
	"github-issue-schedule/internal/pkg/utils"
	"log"
	t "time"
)

var AppVersion = ""

const AppName = "GitHub Issue Scheduler"

func init() {
	if AppVersion == "" {
		AppVersion = "Unknown" // expect the value to be set during the build
	}
}

func main() {
	log.Printf("%v version: %v", AppName, AppVersion)

	// capture all errors
	var errors []error

	// in case of any errors, be ready to handle the case
	defer func() {
		if len(errors) > 0 {
			log.Fatalf("caught error :%+v", errors)
		}
	}()

	// read the configuration file, it has information of which
	// project, which Git repository to be used and the deadline
	// with the template information
	config := configs.ReadConfiguration()

	// get a client object to perform any of the further
	// operations
	client := github.NewClient()

	// for each of the projects, read the configuration
	// identify if there are any reminders to be sent.
	for _, project := range config.Projects {
		// check if there are schedules to be handled
		for _, schedule := range project.Schedules {
			// if the schedule is within the buffer period
			requestedTime, err := utils.GetDate(schedule.Date)
			if err != nil {
				errors = append(errors, err)
				return
			}
			currentTime := t.Now()
			differenceTime := requestedTime.Sub(currentTime)
			// number of days from today to the actual schedule
			numberOfDays := int(differenceTime.Hours() / 24)

			log.Printf("current: %v requested: %v", currentTime, requestedTime)
			log.Printf("difference: %v configured: %v", numberOfDays, config.BufferWindowDays)

			// found a case where action is needed; both sides are ints now so
			// the comparison is straightforward and cannot wrap around.
			if numberOfDays == config.BufferWindowDays {
				// do as scheduled on this date
				projectInstance := client.GetProject(project)
				errInfo := projectInstance.CreateIssue(schedule)
				if errInfo != nil {
					errors = append(errors, errInfo)
					return
				}
			}
		}
	}
}
