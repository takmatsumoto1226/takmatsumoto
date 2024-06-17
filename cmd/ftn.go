package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ftnCmd = &cobra.Command{
	Use:   "lot",
	Short: "Lot Service",
	RunE:  ftnCommand,
}

func init() {
	RootCmd.AddCommand(ftnCmd)
	ftnCmd.PersistentFlags().StringP("action", "a", "", "操作")
}

func ftnCommand(cmd *cobra.Command, args []string) error {

	act, _ := cmd.Flags().GetString("action")

	logrus.Infof("action : %s", act)
	switch act {
	case "pre":
	case "bak":
	default:
	}
	return nil
}
