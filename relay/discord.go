package relay

import (
	"context"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	token     string
	session   *discordgo.Session
	channelID string
	ctx       = context.Background()
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func sendMsg(s *discordgo.Session, c string, msg string) {
	if _, err := s.ChannelMessageSend(c, msg); err != nil {
		log.Errorf("Failed to send message to channel: %s", err.Error())
	}
}

func onConnect(s *discordgo.Session, _ *discordgo.Connect) {
	session = s
	log.Info("Connected to session ws API")
	d := discordgo.UpdateStatusData{
		Game: &discordgo.Game{
			Name:    `xxx`,
			URL:     "git@github.com/leighmacdonald/discord_log_relay",
			Details: "Pew Pew",
		},
	}
	if err := s.UpdateStatusComplex(d); err != nil {
		log.WithError(err).Errorf("Failed to update status complex")
	}
}

func onDisconnect(_ *discordgo.Session, _ *discordgo.Disconnect) {
	log.Info("Disconnected from session ws API")
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	log.Infof(m.Content)
}

func StartBot(ctx context.Context, t string, cID string) {
	if token == "" {
		log.Fatalf("No TOKEN specified")
	}
	if channelID == "" {
		log.Fatalf("No CHANNEL_ID specified")
	}
	channelID = cID
	token = t
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	dg.AddHandler(onConnect)
	dg.AddHandler(onDisconnect)
	dg.AddHandler(onMessageCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening Discord session: ", err)
	}
	listenHost := os.Getenv("LISTEN")
	if listenHost == "" {
		listenHost = ":5555"
	}
	<-ctx.Done()
	// Cleanly close down the Discord session.
	if err := session.Close(); err != nil {
		log.Errorf("Failed to cleanly shut down the session connection: %s", err)
	}
}
