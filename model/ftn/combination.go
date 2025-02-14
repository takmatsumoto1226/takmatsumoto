package ftn

import (
	"fmt"
	"lottery/config"
	"lottery/model/df"
	"os"
	"path/filepath"
)

type Combination struct {
	Balls df.IBalls
}

func loadCombination() error {
	fullPath := fmt.Sprintf(filepath.Join(config.Config.CombinationsInfo[0].Path, config.Config.CombinationsInfo[0].Template), 1)

	file, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	fmt.Printf("%s", file)

	return nil
}
