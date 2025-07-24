package utils

import (
	"fmt"
	"strings"
	"time"
)

func ParseHumanTime(input string) (int64, error) {
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return 0, nil
	}
	if input == "now" {
		return time.Now().Unix(), nil
	}

	now := time.Now()
	currentYear := now.Year()

	if t, err := time.Parse("15:04", input); err == nil {
		full := time.Date(
			currentYear,
			now.Month(),
			now.Day(),
			t.Hour(),
			t.Minute(),
			0,
			0,
			now.Location(),
		)
		return full.Unix(), nil
	}

	if t, err := time.Parse("02.01", input); err == nil {
		full := time.Date(
			currentYear,
			t.Month(),
			t.Day(),
			0,
			0,
			0,
			0,
			now.Location(),
		)
		return full.Unix(), nil
	}

	if t, err := time.Parse("02.01 15:04", input); err == nil {
		full := time.Date(
			currentYear,
			t.Month(),
			t.Day(),
			t.Hour(),
			t.Minute(),
			0,
			0,
			now.Location(),
		)
		return full.Unix(), nil
	}

	return 0, fmt.Errorf("Invalid time format: %q", input)
}
