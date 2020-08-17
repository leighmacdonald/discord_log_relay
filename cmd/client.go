package cmd

import (
	"context"
	"github.com/leighmacdonald/discord_log_relay/client"
	"github.com/leighmacdonald/discord_log_relay/relay"
	"log"

	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:     "client",
	Short:   "",
	Long:    ``,
	Version: relay.BuildVersion,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		logPath := cmd.Flag("logpath").Value.String()
		name := cmd.Flag("name").Value.String()
		host := cmd.Flag("host").Value.String()
		if err := client.New(ctx, name, logPath, host); err != nil {
			log.Fatalf("Client existed early: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringP("logpath", "l", "/home/tf2server/serverfiles/tf/logs", "Help message for toggle")
	clientCmd.Flags().StringP("name", "n", "server-1", "Shorthand server id used to identify it uniquely")
	clientCmd.Flags().StringP("host", "H", "localhost:6666", "Shorthand server id used to identify it uniquely")
}
