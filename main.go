package main

import "fmt"

func main() {
	status, err := GetCurrentStatus()
	if err != nil {
		panic(err)
	}
	fmt.Print(status.Format())
}
