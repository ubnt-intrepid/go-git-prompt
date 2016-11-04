package main

import "fmt"

func main() {
	status, err := GetCurrentStatus()
	if err != nil {
		return
	}
	fmt.Print(status.Format())
}
