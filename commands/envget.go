package commands

import (
	"fmt"

	djbot "github.com/sunho/sdbx-discord-dj-bot"
	"github.com/sunho/sdbx-discord-dj-bot/msg"
	"github.com/sunho/sdbx-discord-dj-bot/stypes"
)

type EnvGet struct {
}

func (eg *EnvGet) Handle(sess *djbot.Session, parms []interface{}) {
	if !sess.AdminCheck() {
		return
	}

	list := [][]string{}
	for key, vars := range sess.GetEnvServer().Env {
		list = append(list, []string{key, fmt.Sprint(vars.Var)})
	}
	msg.LabeledListMsg("Env list", list, sess.UserID, sess.ChannelID, sess.Session)
}

func (eg *EnvGet) Description() string {
	return msg.DescriptionEnvGet
}

func (eg *EnvGet) Types() []stypes.Type {
	return []stypes.Type{}
}
