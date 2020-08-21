package pkg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

type Notifier struct {
	AppName string
	AppIcon string
}

type Notification struct {
	Args []string
}

func (n *Notifier) Push(subtitle string, message string, open string) {
	notification := newNotification()
	notification.addArg("title", n.AppName)
	notification.addArg("message", message)
	notification.addArg("subtitle", subtitle)
	notification.addArg("appIcon", n.AppIcon)
	notification.addArg("open", open)

	err := notification.ToCommand().Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (n *Notifier) PushError(message string, err error) {
	notification := newNotification()
	notification.addArg("title", n.AppName)
	notification.addArg("subtitle", "ERROR")
	notification.addArg("message", fmt.Sprintf("%s: %s", message, err.Error()))

	err = notification.ToCommand().Run()
	if err != nil {
		log.Fatal(err)
	}
}

func newNotification() *Notification {
	return &Notification{[]string{}}
}

func (n *Notification) addArg(key string, value string) {
	n.Args = append(n.Args, fmt.Sprintf("-%s", key), value)
}

func (n *Notification) ToCommand() *exec.Cmd {
	cmd := exec.Command("terminal-notifier")

	for _, arg := range n.Args {
		cmd.Args = append(cmd.Args, arg)
	}

	return cmd
}
