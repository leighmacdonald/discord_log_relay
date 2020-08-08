package cmd

import (
	"context"
	"fmt"
	"github.com/leighmacdonald/discord_log_relay/consts"
	"github.com/leighmacdonald/discord_log_relay/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		if err := server.Server(ctx, fmt.Sprintf(":%d", consts.ListenPort); err != nil {
			log.Errorf("Failed to close server cleanly: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
