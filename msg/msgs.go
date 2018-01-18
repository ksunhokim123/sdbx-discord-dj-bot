package msg

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/ksunhokim123/sdbx-discord-dj-bot/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/ksunhokim123/dgwidgets"
)

func timeOutMsg(sess *discordgo.Session, chID string, msgID string, t time.Duration) {
	timer := time.NewTimer(t)
	<-timer.C
	sess.ChannelMessageDelete(chID, msgID)
}

func TimeOutMsg(sess *discordgo.Session, chID string, msgID string, t time.Duration) {
	go timeOutMsg(sess, chID, msgID, t)
}

func ListMsg(list []string, sess *discordgo.Session, channel string) chan bool {
	p := dgwidgets.NewPaginator(sess, channel)
	for i := 0; i < len(list); i += 10 {
		str := ""
		for j := 0; j < utils.MinInt(10, len(list)-i); j++ {
			str += fmt.Sprintln(i+j, " "+list[i+j])
		}
		p.Add(&discordgo.MessageEmbed{Description: str, Color: 0xffff00})
	}
	p.SetPageFooters()
	p.Widget.Timeout = time.Second * 20
	p.ColourWhenDone = 0xfffff0
	go p.Spawn()
	return p.Widget.Close
}

type CmdList [][]string

func (list CmdList) Len() int {
	return len(list)
}

func (list CmdList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list CmdList) Less(i, j int) bool {
	return list[i][0] < list[j][0]
}

func HelpMsg(list CmdList, sess *discordgo.Session, channel string) {
	sort.Sort(list)
	fields := []*discordgo.MessageEmbedField{}
	for i := 0; i < len(list); i++ {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  list[i][0],
			Value: list[i][1],
		})
	}
	eb := &discordgo.MessageEmbed{
		Title:  "Commands list",
		Fields: fields,
		Color:  0xffff00,
	}
	sess.ChannelMessageSendEmbed(channel, eb)
}
func AddedToQueue(song []string, position int, userid string, channel string, sess *discordgo.Session) {
	usr, _ := sess.User(userid)
	eb := &discordgo.MessageEmbed{
		Title:       song[0],
		Description: "the song has been added to the queue successfully",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "type",
				Value: song[1],
			},
			{
				Name:   "length",
				Value:  song[2],
				Inline: true,
			},
			{
				Name:  "position",
				Value: strconv.Itoa(position),
			},
		},
		Color: 0xffff00,
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: usr.AvatarURL(""),
			Text:    usr.Username,
		},
	}
	sess.ChannelMessageSendEmbed(channel, eb)
}