package relay

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	session   *discordgo.Session
	channelID string
)

func SendPayload(p Payload) error {
	team := ""
	if p.SayTeam {
		team = "(team) "
	}
	format := "`[%s] %d` **%s**: %s%s"
	return SendMsg(channelID, fmt.Sprintf(format, p.Server, p.SteamID.Int64(), p.Username, team, p.Message))
}

func SendMsg(channel string, msg string) error {
	if _, err := session.ChannelMessageSend(channel, msg); err != nil {
		log.Errorf("Failed to send message to channel: %v", err)
		return err
	}
	return nil
}

func onConnect(s *discordgo.Session, _ *discordgo.Connect) {
	session = s
	log.Info("Connected to session ws API")
	d := discordgo.UpdateStatusData{
		Game: &discordgo.Game{
			Name:    `Uncletopia`,
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

func StartDiscord(ctx context.Context, token string, cID string) error {
	if token == "" {
		log.Fatalf("No TOKEN specified")
	}
	if cID == "" {
		log.Fatalf("No CHANNEL_ID specified")
	}
	channelID = cID
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return err
	}
	dg.AddHandler(onConnect)
	dg.AddHandler(onDisconnect)
	dg.AddHandler(onMessageCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening Discord session: ", err)
		return err
	}
	listenHost := os.Getenv("LISTEN")
	if listenHost == "" {
		listenHost = ":5555"
	}
	<-ctx.Done()
	// Cleanly close down the Discord session.
	if err := session.Close(); err != nil {
		log.Errorf("Failed to cleanly shut down the session connection: %s", err)
		return err
	}
	return nil
}
