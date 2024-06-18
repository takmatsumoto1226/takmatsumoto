package cmd

import (
	"lottery/config"
	"lottery/model/df"
	"lottery/model/ftn"
	"lottery/model/interf"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ftnCmd = &cobra.Command{
	Use:   "ftn",
	Short: "FTN Service",
	RunE:  ftnCommand,
}

func init() {
	RootCmd.AddCommand(ftnCmd)
	ftnCmd.PersistentFlags().StringP("action", "a", "", "操作")
	ftnCmd.PersistentFlags().StringP("round", "r", "", "幾回")
	ftnCmd.PersistentFlags().StringP("valie", "v", "", "閾值")
	ftnCmd.PersistentFlags().StringP("sample", "s", "", "取樣倍數")
	ftnCmd.PersistentFlags().StringP("rtype", "t", "", "random type")
}

func ftnCommand(cmd *cobra.Command, args []string) error {

	act, _ := cmd.Flags().GetString("action")

	r, _ := cmd.Flags().GetString("round")
	ir, _ := strconv.Atoi(r)
	if ir == 0 {
		ir = 10
	}

	v, _ := cmd.Flags().GetString("value")
	iv, _ := strconv.Atoi(v)
	if iv == 0 {
		iv = 9
	}

	s, _ := cmd.Flags().GetString("sample")
	is, _ := strconv.ParseFloat(s, 64)
	if is == 0 {
		is = 5
	}

	logrus.Infof("action : %s", act)
	switch act {
	case "pre":
		config.LoadConfig("./config.yaml")

		var ar = ftn.FTNsManager{}
		ar.Prepare()

		start := 0
		//
		df.DisableFilters([]int{df.FilterOddCount, df.FilterEvenCount, df.FilterTailDigit})
		// df.DisableFilters([]int{df.FilterTailDigit})
		th := interf.Threshold{
			Round:      ir,
			Value:      iv,
			SampleTime: is,
			Sample:     len(ar.Combinations),
			Interval: interf.Interval{
				Index:  start,
				Length: len(ar.List) / 4,
			},
			Smart: interf.Smart{
				Enable: true,
				Type:   interf.RangeTypeLatestRange,
			},
			Randomer: 1,
			Match:    false,
		}

		ar.JSONGenerateTopPriceNumber(th)
		ar.SaveBTs()
	case "bak":
	default:
	}

	os.Exit(0)
	return nil
}
