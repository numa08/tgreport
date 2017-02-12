package main

import (
	"github.com/numa08/tgreport/lib/config"
	"fmt"
	"os"
	"github.com/numa08/tgreport/lib/toggl"
	"github.com/numa08/tgreport/lib/report"
)

func Run(args []string) int {
 	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("設定ファイルの読み込みに失敗しました.\n %s", err))
		return 1
	}
	r, err := toggl.LastWeekReport(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("レポートの読み込みに失敗しました.\n %s", err))
	}
	f := getFormatter(args)
	fmt.Fprintln(os.Stdout, f.Format(r))
	return 0
}

func getFormatter(args []string) report.Formatter {
	return report.NewMarkdownFormatter()
}