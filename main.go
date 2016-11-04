package main

import "fmt"

func main() {
	status, err := GetCurrentStatus()
	if err != nil {
		return
	}
	fmt.Println(status)
}
