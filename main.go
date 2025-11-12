package main

import "github.com/ColonyPM/cpm-cli/cmd"

func add(a int, b int) int {
	return a + b
}

func main() {
	cmd.Execute()
}
