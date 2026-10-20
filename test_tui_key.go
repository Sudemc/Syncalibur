package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	sttyOrigCmd := exec.Command("stty", "-g")
	sttyOrigCmd.Stdin = os.Stdin
	sttyOrigOut, err := sttyOrigCmd.Output()
	if err != nil {
		fmt.Printf("stty -g error: %v\n", err)
		return
	}
	sttyOrig := strings.TrimSpace(string(sttyOrigOut))
	fmt.Printf("Captured sttyOrig: %s\n", sttyOrig)

	rawCmd := exec.Command("stty", "-icanon", "-echo", "min", "1", "time", "0")
	rawCmd.Stdin = os.Stdin
	if err := rawCmd.Run(); err != nil {
		fmt.Printf("rawCmd error: %v\n", err)
		return
	}
	defer func() {
		c := exec.Command("stty", sttyOrig)
		c.Stdin = os.Stdin
		_ = c.Run()
		fmt.Println("\nRestored terminal!")
	}()

	fmt.Println("Raw mode active! Press 'q' or any key within 3 seconds...")
	keyCh := make(chan byte, 1)
	go func() {
		buf := make([]byte, 1)
		n, err := os.Stdin.Read(buf)
		if err == nil && n > 0 {
			keyCh <- buf[0]
		}
	}()

	select {
	case k := <-keyCh:
		fmt.Printf("Key received: %q (%d)\n", k, k)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("Timed out (no key pressed, as expected in non-interactive background test)")
	}
}
