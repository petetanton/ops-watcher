package main

import (
	"log"
	"runtime"
	"strings"

	"gopkg.in/robfig/cron.v2"

	"github.com/petetanton/ops-watcher/pkg"
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

	var watchers []pkg.Watcher

	if config.JiraEnabled {
		jiraWatchers, err := pkg.NewJiraWatchers(config)
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

	log.Printf("Setup: %s", msg)

	c := cron.New()
	c.AddFunc("* * * * *", func() { runWatchers(watchers, notifier) })
	c.Start()

	runtime.Goexit()

	//time.Sleep(10 * time.Second)
	//c.Stop()
}

func runWatchers(watchers []pkg.Watcher, notifier pkg.Notifier) {
	log.Println("Running watchers")
	var ids string
	for _, watcher := range watchers {
		notifications, err := watcher.Watch()
		if err != nil {
			notifier.PushError("error when creating notifications", err)
		} else {
			for _, notification := range notifications {
				if strings.Contains(ids, notification.Id) == false {
					err := notification.ToCommand().Run()
					if err != nil {
						log.Fatal(err)
					}
					ids += notification.Id
				}
			}
		}
	}
}
