package cmd

import (
	"log"
	"time"

	"github.com/google/go-github/v60/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chia-network/github-bot/internal/config"
	"github.com/chia-network/github-bot/internal/label"
)

// labelPRsCmd represents the labelPRs command
var labelPRsCmd = &cobra.Command{
	Use:   "label-prs",
	Short: "Adds community and internal labels to pull requests in designated repos",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(viper.GetString("config"))
		if err != nil {
			log.Fatalf("error loading config: %s\n", err.Error())
		}
		client := github.NewClient(nil).WithAuthToken(cfg.GithubToken)

		loop := viper.GetBool("loop")
		loopDuration := viper.GetDuration("loop-time")
		for {
			log.Println("Labeling Pull Requests")
			err = label.PullRequests(client, cfg.InternalTeam, cfg.LabelConfig)
			if err != nil {
				log.Fatalln(err.Error())
			}

			if !loop {
				break
			}

			log.Printf("Waiting %s for next iteration\n", loopDuration.String())
			time.Sleep(loopDuration)
		}
	},
}

func init() {
	rootCmd.AddCommand(labelPRsCmd)
}
