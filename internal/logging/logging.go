package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func New() (*slog.Logger, io.Writer, io.Closer, error) {
	stateDir, err := os.UserConfigDir()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get user config directory: %w", err)
	}

	logDir := filepath.Join(stateDir, "o7k")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, nil, nil, fmt.Errorf("create log directory: %w", err)
	}

	logPath := filepath.Join(logDir, "o7k.log")
	file, err := os.OpenFile(
		logPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
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
