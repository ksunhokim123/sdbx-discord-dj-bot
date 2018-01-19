package commands

import (
	"bufio"
	"encoding/binary"
	"os/exec"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type MusicStart struct {
	Music *Music
}

func (mc *MusicStart) Handle(sess *djbot.Session, parms []interface{}) {
	mc.Music.GetServer(sess.ServerID).Start(sess)
}

func (vc *MusicStart) Description() string {
	return msg.DescriptionMusicStart
}

func (vc *MusicStart) Types() []stypes.Type {
	return []stypes.Type{}
}

func (m *MusicServer) PlayOne(sess *djbot.Session) {
	url := m.Songs[0].Url
	ytdl := exec.Command("./youtube-dl", "-v", "-f", "bestaudio", "-o", "-", url)
	ytdlout, err := ytdl.StdoutPipe()
	if err != nil {
		sess.Send(err)
		return
	}
	ffmpeg := exec.Command("./ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpegout, err := ffmpeg.StdoutPipe()
	ffmpeg.Stdin = ytdlout
	if err != nil {
		sess.Send(err)
		return
	}
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	dca := exec.Command("./dca")
	dca.Stdin = ffmpegbuf
	dcaout, err := dca.StdoutPipe()
	if err != nil {
		sess.Send(err)
		return
	}
	dcabuf := bufio.NewReaderSize(dcaout, 16384)
	err = ytdl.Start()
	if err != nil {
		sess.Send(err)
		return
	}
	defer func() {
		go ytdl.Wait()
	}()
	err = ffmpeg.Start()

	if err != nil {
		sess.Send(err)
		return
	}
	defer func() {
		go ffmpeg.Wait()
	}()

	err = dca.Start()
	if err != nil {
		sess.Send(err)
		return
	}
	defer func() {
		go dca.Wait()
	}()
	var opuslen int16
	sess.VoiceConnection.Speaking(true)
	defer sess.VoiceConnection.Speaking(false)
	done := true
	for done {
		select {
		case <-m.SkipChan:
			done = false
		default:
			err = binary.Read(dcabuf, binary.LittleEndian, &opuslen)
			if err != nil {
				done = false
				break
			}
			opus := make([]byte, opuslen)
			err = binary.Read(dcabuf, binary.LittleEndian, &opus)
			if err != nil {
				done = false
				break
			}
			if sess.VoiceConnection == nil {
				done = false
				break
			}
			sess.VoiceConnection.OpusSend <- opus
		}
	}
}

func (m *MusicServer) Start(sess *djbot.Session) {
	if m.State == Playing {
		sess.Send(msg.NoJustATrick)
		return
	}
	if len(m.Songs) == 0 {
		sess.Send(msg.NoQueue)
		return
	}
	m.State = Playing
	for {
		if sess.VoiceConnection == nil {
			break
		}
		if len(m.Songs) == 0 {
			break
		}
		m.PlayOne(sess)
		if len(m.Songs) != 0 {
			m.Next()
		} else {
			break
		}
	}
	m.State = NotPlaying
}
