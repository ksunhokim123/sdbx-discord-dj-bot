package commands

import (
	"github.com/bwmarrin/discordgo"
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/envs"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type VoiceConnect struct {
}

func (vc *VoiceConnect) Handle(sess *djbot.Session, parms []interface{}) {
	if sess.GetEnvServer().GetEnv(envs.VOICECHANNELONLY).(bool) {
		channel := sess.GetEnvServer().GetEnv(envs.VOICECHANNEL)
		vc.Connect(sess, channel)
		return
	}

	gd, _ := sess.Guild(sess.ServerID)
	slist := []string{}
	dlist := []interface{}{}
	for _, ch := range gd.Channels {
		if ch.Type == discordgo.ChannelTypeGuildVoice {
			dlist = append(dlist, ch.ID)
			slist = append(slist, ch.Name)
		}
	}

	sess.DJBot.RequestManager.Set(sess, &djbot.Request{
		List:     slist,
		DataList: dlist,
		CallBack: vc.Connect,
	})
}
func (vc *VoiceConnect) Connect(sess *djbot.Session, id interface{}) {
	id2, ok := id.(string)
	if ok {
		ch, err := sess.Channel(id.(string))
		if err != nil {
			sess.Send(msg.NoJustATrick)
			return
		}
		if ch.Type != discordgo.ChannelTypeGuildVoice {
			sess.Send(msg.NoJustATrick)
			return
		}
		vc, err := sess.ChannelVoiceJoin(sess.ServerID, id2, false, true)
		if err != nil {
			sess.Send(msg.NoJustATrick)
			return
		}
		sess.DJBot.VoiceConnections[sess.ServerID] = vc
	}
}
func (vc *VoiceConnect) Description() string {
	return "this will make the bot connect to your voice channel"
}
func (vc *VoiceConnect) Types() []stypes.Type {
	return []stypes.Type{}
}
