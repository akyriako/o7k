package terminal

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func CheckK9s() error {
	if _, err := exec.LookPath("k9s"); err != nil {
		return fmt.Errorf("k9s executable not found in PATH: %w", err)
	}
	return nil
}

func OpenK9s(kubeconfig string) error {
	k9s, err := exec.LookPath("k9s")
	if err != nil {
		return fmt.Errorf("finding k9s failed: %w", err)
	}

	switch runtime.GOOS {
	case "linux":
		candidates := [][]string{
			{"gnome-terminal", "--", k9s, "--kubeconfig", kubeconfig},
			{"konsole", "-e", k9s, "--kubeconfig", kubeconfig},
			{"kitty", "--", k9s, "--kubeconfig", kubeconfig},
			{"alacritty", "-e", k9s, "--kubeconfig", kubeconfig},
			{"xterm", "-e", k9s, "--kubeconfig", kubeconfig},
		}

		for _, args := range candidates {
			if _, err := exec.LookPath(args[0]); err == nil {
				return exec.Command(args[0], args[1:]...).Start()
			}
		}

		return fmt.Errorf(
			"no supported terminal found",
		)

	case "darwin":
		command := shellQuote(k9s) + " --kubeconfig " + shellQuote(kubeconfig)
		script := `on run argv
    tell application "Terminal"
        activate
        do script (item 1 of argv)
    end tell
end run`

		out, err := exec.Command("osascript", "-e", script, command).CombinedOutput()
		if err != nil {
			return fmt.Errorf("opening Terminal.app: %w: %s", err, strings.TrimSpace(string(out)))
		}

		return nil

	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
