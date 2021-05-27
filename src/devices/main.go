package main

import "fmt"

func main() {
	comm := make(chan error)

	go startProbe(comm)
	go startTelescope(comm)

	fmt.Println("[*] Press CTRL+C to interrupt.")
	for err := range comm {
		fmt.Print(err)
	}
}
