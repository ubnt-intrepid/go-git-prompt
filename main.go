package main

import (
	"flag"
	"fmt"

	"github.com/ubnt-intrepid/go-git-prompt/color"
	"github.com/ubnt-intrepid/go-git-prompt/prompt"
)

func main() {
	status, err := prompt.GetCurrentStatus()
	if err != nil {
		return
	}

	var fColored = flag.String("colored", "default", "colored library (default, zsh)")
	flag.Parse()

	var colored color.Colored
	if *fColored == "zsh" {
		colored = color.NewZshColoredOutput()
	} else {
		colored = color.NewDefaultColoredOutput()
	}
	fmt.Print(status.Format(colored))
}
