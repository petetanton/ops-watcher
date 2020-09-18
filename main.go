package main

import (
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"

	"github.com/petetanton/ops-watcher/pkg"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	notifier := pkg.Notifier{
		AppName: "OPS Watcher",
		AppIcon: "danger.png",
	}

	config, err := pkg.NewConfig()
	if err != nil {
		notifier.PushError("got an error parsing config", err)
	}

	var watchers []pkg.Watcher

	if config.JiraEnabled {
		jiraWatchers, err := pkg.NewJiraWatchers(config, logger)
		if err != nil {
			notifier.PushError("error when creating JiraWatcher", err)
		}

		for _, jiraWatcher := range jiraWatchers {
			watchers = append(watchers, jiraWatcher)
		}

	}

	msg := "setting up. JIRA: "
	if config.JiraEnabled {
		msg += "ENABLED"
	} else {
		msg += "DISABLED"
	}

	logger.Infof("Setup: %s", msg)

	c := cron.New()
	_, err = c.AddFunc("* * * * *", func() { runWatchers(watchers, notifier, logger) })
	if err != nil {
		logger.Fatal(err)
	}
	c.Start()

	runtime.Goexit()
}

func runWatchers(watchers []pkg.Watcher, notifier pkg.Notifier, logger logrus.FieldLogger) {
	logger.Info("Running watchers")
	var ids string
	for _, watcher := range watchers {
		notifications, err := watcher.Watch()
		if err != nil {
			notifier.PushError("error when creating notifications", err)
		} else {
			for _, notification := range notifications {
				if !strings.Contains(ids, notification.Id) {
					err := notification.ToCommand().Run()
					if err != nil {
						logger.Fatal(err)
					}
					ids += notification.Id
				}
			}
		}
	}
}
