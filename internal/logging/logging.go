package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func New() (*slog.Logger, io.Writer, io.WriteCloser, error) {
	stateDir, err := userStateDir()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get user state directory: %w", err)
	}

	logDir := filepath.Join(stateDir, "o7k")
	if err := os.MkdirAll(logDir, 0o700); err != nil {
		return nil, nil, nil, fmt.Errorf("create log directory: %w", err)
	}

	logPath := filepath.Join(logDir, "o7k.log")
	file, err := os.OpenFile(
		logPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o600,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("open log file: %w", err)
	}

	handler := slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)

	stderr := io.MultiWriter(os.Stderr, file)

	return logger, stderr, file, nil
}

func userStateDir() (string, error) {
	if dir := os.Getenv("XDG_STATE_HOME"); filepath.IsAbs(dir) {
		return dir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".local", "state"), nil
}
