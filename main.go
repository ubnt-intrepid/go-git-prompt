package main

import (
	"fmt"

	"github.com/ubnt-intrepid/go-git-prompt/src"
)

func main() {
	status, err := prompt.GetCurrentStatus()
	if err != nil {
		return
	}
	fmt.Println(status.Format())
}
