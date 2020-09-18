package pkg

import (
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type Notifier struct {
	AppName string
	AppIcon string
	Logger  logrus.FieldLogger
}

type Notification struct {
	Id   string
	Args []string
}

func (n *Notifier) Push(subtitle string, message string, open string) {
	notification := newNotification(subtitle)
	notification.addArg("title", n.AppName)
	notification.addArg("message", message)
	notification.addArg("subtitle", subtitle)
	//notification.addArg("appIcon", n.AppIcon)
	notification.addArg("open", open)

	err := notification.ToCommand().Run()
	if err != nil {
		n.Logger.Fatal(err)
	}
}

func (n *Notifier) PushError(message string, err error) {
	notification := newNotification("err")
	notification.addArg("title", n.AppName)
	notification.addArg("subtitle", "ERROR")
	notification.addArg("message", fmt.Sprintf("%s: %s", message, err.Error()))

	err = notification.ToCommand().Run()
	if err != nil {
		n.Logger.Fatal(err)
	}
}

func newNotification(id string) *Notification {
	return &Notification{id, []string{"-appIcon", "face.png"}}
}

func (n *Notification) addArg(key string, value string) {
	n.Args = append(n.Args, fmt.Sprintf("-%s", key), value)
}

func (n *Notification) ToCommand() *exec.Cmd {
	cmd := exec.Command("terminal-notifier")
	cmd.Args = append(cmd.Args, n.Args...)

	return cmd
}
