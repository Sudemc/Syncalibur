package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()
	fmt.Printf("With Stdin: err=%v, out=%q\n", err, string(out))
}
