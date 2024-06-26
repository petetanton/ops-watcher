package main

import (
	"context"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"

	"github.com/petetanton/ops-watcher/pkg"
)

func getLogger() logrus.FieldLogger {
	logger := logrus.New()

	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = logrus.InfoLevel
		logger.Infof("log level not set so defaulting to INFO. To set a different log level, set LOG_LEVEL in your environment")
	}
	logger.SetLevel(level)

	return logger
}

func main() {
	ctx := context.Background()

	var err error
	logger := getLogger()

	notifier := pkg.Notifier{
		AppName: "OPS Watcher",
		AppIcon: "face.png",
	}

	config, err := pkg.NewConfig()
	if err != nil {
		notifier.PushError("got an error parsing config", err)
		logger.Fatal("got an error parsing config", err)
	}

	var watchers []pkg.Watcher

	if config.JiraEnabled {
		jiraWatchers, err2 := pkg.NewJiraWatchers(ctx, config, logger)
		if err2 != nil {
			logger.Error(err2)
			notifier.PushError("error when creating JiraWatcher", err2)
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
	_, err = c.AddFunc("* * * * *", func() { runWatchers(ctx, watchers, notifier, logger) })
	if err != nil {
		logger.Fatal(err)
	}
	c.Start()

	runtime.Goexit()
}

func runWatchers(ctx context.Context, watchers []pkg.Watcher, notifier pkg.Notifier, logger logrus.FieldLogger) {
	logger.Info("Running watchers")
	var ids string
	for _, watcher := range watchers {
		notifications, err := watcher.Watch(ctx)
		if err != nil {
			notifier.PushError("error when creating notifications", err)
			logger.Error(err)
		} else {
			for _, notification := range notifications {
				if !strings.Contains(ids, notification.Id) {
					command := notification.ToCommand()
					logger.Debugf("running: %s", command.String())
					err := command.Run()
					if err != nil {
						logger.Fatal(err)
					}
					ids += notification.Id
				}
			}
		}
	}
}
