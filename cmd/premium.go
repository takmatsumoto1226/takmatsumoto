package cmd

import (
	"errors"
	"fmt"
	"lottery/api"
	"lottery/config"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// 手動簽發付費憑證, 例 : go run . premium -u user123 -d 30
var premiumCmd = &cobra.Command{
	Use:   "premium",
	Short: "簽發付費憑證",
	RunE:  premiumCommand,
}

func init() {
	RootCmd.AddCommand(premiumCmd)
	premiumCmd.Flags().StringP("user", "u", "", "使用者 ID (A-Z a-z 0-9 _ -)")
	premiumCmd.Flags().IntP("days", "d", 30, "有效天數")
}

func premiumCommand(cmd *cobra.Command, args []string) error {
	user, _ := cmd.Flags().GetString("user")
	days, _ := cmd.Flags().GetInt("days")
	if user == "" || days <= 0 {
		return errors.New("必須指定 --user 與正數的 --days")
	}

	config.LoadConfig("./config.yaml")
	expires := time.Now().AddDate(0, 0, days)
	token, err := api.IssuePremiumToken(config.Config.Server.PremiumSecret, user, expires)
	if err != nil {
		return err
	}
	fmt.Printf("user    : %s\nexpires : %s\ntoken   : %s\n", user, expires.Format(time.RFC3339), token)
	os.Exit(0)
	return nil
}
