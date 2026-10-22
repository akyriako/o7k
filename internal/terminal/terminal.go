package terminal

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func Open(command string, args []string, preferred *string) error {
	executable, err := exec.LookPath(command)
	if err != nil {
		return fmt.Errorf("executable %q not found: %w", command, err)
	}

	switch runtime.GOOS {
	case "linux":
		candidates := [][]string{
			{"gnome-terminal", "--", executable},
			{"konsole", "-e", executable},
			{"kitty", "--", executable},
			{"alacritty", "-e", executable},
			{"xterm", "-e", executable},
			{"ghostty", "-e", executable},
		}

		if preferred != nil {
			for _, candidate := range candidates {
				if candidate[0] == *preferred {
					if _, err := exec.LookPath(candidate[0]); err != nil {
						return fmt.Errorf("configured terminal %q not found: %w", *preferred, err)
					}
					return exec.Command(candidate[0], append(candidate[1:], args...)...).Start()
				}
			}
			return fmt.Errorf("unsupported terminal %q on Linux", *preferred)
		}

		for _, candidate := range candidates {
			if _, err := exec.LookPath(candidate[0]); err == nil {
				return exec.Command(candidate[0], append(candidate[1:], args...)...).Start()
			}
		}

		return fmt.Errorf("no supported terminal found")

	case "darwin":
		candidates := []string{"Terminal", "iTerm"}

		if preferred != nil {
			candidates = []string{*preferred}
		}

		for _, candidate := range candidates {
			if candidate != "Terminal" && candidate != "iTerm" {
				return fmt.Errorf("unsupported terminal %q on macOS", candidate)
			}

			if err := exec.Command("open", "-Ra", candidate).Run(); err != nil {
				if preferred != nil {
					return fmt.Errorf("configured terminal %q not found: %w", candidate, err)
				}
				continue
			}

			parts := append([]string{executable}, args...)
			for i, part := range parts {
				parts[i] = shellQuote(part)
			}
			command := strings.Join(parts, " ")

			var script string

			switch candidate {
			case "Terminal":
				script = `on run argv
    tell application "Terminal"
        activate
        do script (item 1 of argv)
    end tell
end run`

			case "iTerm":
				script = `on run argv
    tell application "iTerm"
        activate
        create window with default profile
        tell current session of current window
            write text (item 1 of argv)
        end tell
    end tell
end run`
			}

			out, err := exec.Command("osascript", "-e", script, command).CombinedOutput()
			if err != nil {
				return fmt.Errorf("opening %s failed: %w: %s", candidate, err, strings.TrimSpace(string(out)))
			}

			return nil
		}

		return fmt.Errorf("no supported terminal found on macOS")

	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
