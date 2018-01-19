package commands

import (
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type MusicQueue struct {
	Music *Music
}

//TODO: place this to msg
func (mc *MusicQueue) Handle(sess *djbot.Session, parms []interface{}) {
	list := []string{}
	songs := mc.Music.GetServer(sess.ServerID).Songs
	for i := 0; i < len(songs); i++ {
		list = append(list, "`"+songs[i].Name+"`  **"+songs[i].Duration.String()+"**  Requested by "+songs[i].Requester)
	}
	msg.ListMsg(list, sess.UserID, sess.ChannelID, sess.Session)
}

func (vc *MusicQueue) Description() string {
	return msg.DescriptionMusicQueue
}

func (vc *MusicQueue) Types() []stypes.Type {
	return []stypes.Type{}
}
