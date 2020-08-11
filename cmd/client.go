package cmd

import (
	"context"
	"fmt"
	"github.com/leighmacdonald/discord_log_relay/client"
	"github.com/leighmacdonald/discord_log_relay/consts"
	"log"

	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		if err := client.New(ctx, "blah", "", fmt.Sprintf("localhost:%d", consts.ListenPort)); err != nil {
			log.Fatalf("Client existed early: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
