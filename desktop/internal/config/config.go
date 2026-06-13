package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Settings struct {
	WorkDuration   time.Duration `json:"work_duration_ns"`
	RestDuration   time.Duration `json:"rest_duration_ns"`
	OverlayStyle   string        `json:"overlay_style"`   // expressive, calm, minimal
	AlertMode      string        `json:"alert_mode"`      // both, overlay, notification
	Opacity        float64       `json:"opacity"`
	SoundEnabled   bool          `json:"sound_enabled"`
	AutoStart      bool          `json:"auto_start"`       // start timer on app launch
	AutoStartBoot  bool          `json:"auto_start_boot"`  // start app on Windows boot
	StartMinimized bool          `json:"start_minimized"`
}

func Default() *Settings {
	return &Settings{
		WorkDuration:   20 * time.Minute,
		RestDuration:   20 * time.Second,
		OverlayStyle:   "expressive",
		AlertMode:      "both",
		Opacity:        0.90,
		SoundEnabled:   true,
		AutoStart:      false,
		AutoStartBoot:  false,
		StartMinimized: false,
	}
}

func configPath() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "Blink", "settings.json")
}

func Load() *Settings {
	s := Default()
	data, err := os.ReadFile(configPath())
	if err != nil {
		return s
	}
	_ = json.Unmarshal(data, s)
	if s.WorkDuration < time.Minute {
		s.WorkDuration = time.Minute
	}
	if s.RestDuration < 5*time.Second {
		s.RestDuration = 5 * time.Second
	}
	if s.Opacity < 0.3 {
		s.Opacity = 0.3
	}
	if s.Opacity > 1.0 {
		s.Opacity = 1.0
	}
	if s.AlertMode == "" {
		s.AlertMode = "both"
	}
	return s
}

func (s *Settings) Save() error {
	path := configPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (s *Settings) WorkMinutes() int {
	return int(s.WorkDuration.Minutes())
}

func (s *Settings) RestSeconds() int {
	return int(s.RestDuration.Seconds())
}
