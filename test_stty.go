package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("stty", "-g")
	out, err := cmd.CombinedOutput()
	fmt.Printf("Without Stdin: err=%v, out=%q\n", err, string(out))
}
