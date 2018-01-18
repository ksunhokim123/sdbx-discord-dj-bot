package djbot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (base *DJBot) HandleNewMessage(s *discordgo.Session, msg2 *discordgo.MessageCreate) {
	if msg2.Author.ID == s.State.User.ID {
		return
	}
	ch, err := s.Channel(msg2.ChannelID)
	if err != nil {
		fmt.Println("s.Channel(msg.ChannelID) something is wrong definitely:", err)
		return
	}
	var sess = &Session{
		Session:   s,
		ChannelID: msg2.ChannelID,
		ServerID:  ch.GuildID,
		DJBot:     base,
		Msg:       msg2,
		UserID:    msg2.Author.ID,
	}
	if vc, ok := base.VoiceConnections[sess.ServerID]; ok {
		sess.VoiceConnection = vc
	}
	if len(msg2.Content) != 0 {
		/*go*/ base.CommandMannager.HandleMessage(sess, msg2) // discord go already goed this (go eh.eventHandler.Handle(s, i))
		/*go*/ base.RequestManager.HandleMessage(sess, msg2)
	}

}