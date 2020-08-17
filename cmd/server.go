package cmd

import (
	"context"
	"fmt"
	"github.com/leighmacdonald/discord_log_relay/consts"
	"github.com/leighmacdonald/discord_log_relay/relay"
	"github.com/leighmacdonald/discord_log_relay/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "",
	Long:    ``,
	Version: relay.BuildVersion,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		token := cmd.Flag("token").Value.String()
		channel := cmd.Flag("channel").Value.String()
		go func() {
			if err := relay.StartDiscord(ctx, token, channel); err != nil {
				log.Errorf("Bot shutdown uncleanly: %v", err)
			}
		}()
		if err := server.Server(ctx, fmt.Sprintf(":%d", consts.ListenPort)); err != nil {
			log.Errorf("Failed to close server cleanly: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringP("token", "t", "", "Discord bot token")
	serverCmd.Flags().StringP("channel", "c", "", "Discord channel ID to send messages")
}
