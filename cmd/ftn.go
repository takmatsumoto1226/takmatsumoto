package cmd

import (
	"errors"
	"fmt"
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
	ftnCmd.PersistentFlags().StringP("dir", "D", "", "data 路徑")
	ftnCmd.PersistentFlags().StringP("action", "a", "", "操作")
	ftnCmd.PersistentFlags().StringP("round", "r", "", "幾回")
	ftnCmd.PersistentFlags().StringP("valie", "v", "", "閾值")
	ftnCmd.PersistentFlags().StringP("sample", "s", "", "取樣倍數")
	ftnCmd.PersistentFlags().StringP("rtype", "t", "", "random type")
}

func ftnCommand(cmd *cobra.Command, args []string) error {

	act, _ := cmd.Flags().GetString("action")
	if len(act) == 0 {
		return errors.New("no action!!! Check Please")
	}

	dir, err := cmd.Flags().GetString("dir")
	if dir == "" || err != nil {
		dir = "/Users/tak 1/Documents/gitlab_project/takmatsumoto/model/ftn/gendata/"
	}

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
	config.LoadConfig("./config.yaml")

	var ar = ftn.FTNsManager{}
	ar.Prepare()

	logrus.Infof("action : %s", act)
	switch act {
	case "pre":
		start := 0
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
		fmt.Println(th.Presentation())

		ar.JSONGenerateTopPriceNumber(th)
		ar.SaveBTsWithDir(dir)
	case "pick":
		// files, _ := os.ReadDir(filepath.Join(dir))
		// filenames := []string{}
		// for _, f := range files {
		// 	if strings.Contains(f.Name(), ".json") {
		// 		filenames = append(filenames, filepath.Join(dir, f.Name()))
		// 	}
		// }
		// top := ar.List.GetNode(0)
		// group := ftn.NewGroup(100, ar.Combinations, ar.List)
		// p := ftn.PickParam{SortType: df.Descending, Interval: 20, Whichfront: df.Normal, Freq: 655}
		// infl1s := ar.List.FragmentRange([]int{})
		// exfl2s := ar.List.FragmentRange([]int{0})
		// filterPick := ar.FilterByGroupIndex(group, []int{0, 1}).FilterHighFreqNumber(ar.List, p).FilterPickBySpecConfition().FilterIncludes(infl1s, []int{}).FilterExcludes(exfl2s, []int{}).FilterExcludeNode(ar.List).FilterCol(&top, 0).FilterNeighber(&top, 0).FilterByTebGroup([]int{df.FeatureTenGroup2}, []int{3}).FilterFeatureExcludes(ar.List).findNumbers([]string{}, df.None).Distinct()
		// filterPick.ShowAll()
		// fmt.Println(len(filterPick))
		// fmt.Println(filterPick.IntervalBallsCountStatic(p).AppearBalls.Presentation(true))
		// fmt.Println("got top")

		// for _, f := range filterPick {
		// 	if f.IsSame(&top) {
		// 		fmt.Println("Oooooohhhhh My God!!!  it's " + f.formRow())
		// 	}
		// }

		// fmt.Printf("\n\n\nGod Pick....\n")
		// ar.GodPick(filterPick, 1)

		// ar.List.WithRange(0, 20).Reverse().ShowAll()
	case "bak":
	default:
	}

	os.Exit(0)
	return nil
}
