package utils

import (
	
	"runtime"
	"os/exec"
	"fmt"
)

// openURL opens the provided URL in the default browser
func OpenURL(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url) // Linux command to open the default browser
	case "darwin":
		cmd = exec.Command("open", url) // macOS command to open the default browser
	case "windows":
		cmd = exec.Command("start", url) // Windows command to open the default browser
	default:
		fmt.Println("Unsupported OS")
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error opening URL: %s\n", err)
	}
}

