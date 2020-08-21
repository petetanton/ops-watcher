package main

import (
	"github.com/petetanton/ops-watcher/pkg"
	"log"
)

func main() {
	notifier := pkg.Notifier{
		AppName: "OPS Watcher",
		AppIcon: "danger.png",
	}

	config, err := pkg.NewConfig()
	if err != nil {
		notifier.PushError("got an error parsing config", err)
	}

	if config.JiraEnabled {
		watcher, err := pkg.NewJiraWatcher(config)
		if err != nil {
			notifier.PushError("error when creating JiraWatcher", err)
		}

		notifications, err := watcher.Watch()
		if err != nil {
			notifier.PushError("error when creating notifcations", err)
		} else {
			for _, notification := range notifications {
				err := notification.ToCommand().Run()
				if err != nil {
					log.Fatal(err)
				}
			}
		}

	}

	msg := "setting up. JIRA: "
	if config.JiraEnabled {
		msg += "ENABLED"
	} else {
		msg += "DISABLED"
	}

	notifier.Push("Setup", msg, "https://github.com/petetanton/ops-watcher")

}
