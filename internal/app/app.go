package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

const (
	defaultRestDuration = 20 * time.Second
	minRestDuration     = 5 * time.Second
	maxRestDuration     = 2 * time.Minute
	restDurationStep    = 5 * time.Second
)

var (
	workDuration   = 20 * time.Minute
	restDuration   = defaultRestDuration
	restDurationMu sync.RWMutex
	ticker         *time.Ticker
	stopChan       chan struct{}
	isRunning      bool
	appVersion     = DefaultVersion
)

func Run(version string) {
	if version != "" {
		appVersion = version
	}

	beeep.AppName = Name
	systray.Run(onReady, onExit)
}

func onReady() {
	loadConfig()

	systray.SetTitle("👁️ 20m")
	systray.SetTooltip(DisplayName(appVersion) + " - Pengingat Istirahat Mata 20-20-20")

	mStatus := systray.AddMenuItem("Status: Berjalan", "Status timer")
	mStatus.Disable()
	mVersion := systray.AddMenuItem("Version: "+appVersion, DisplayName(appVersion))
	mVersion.Disable()
	systray.AddSeparator()

	mToggle := systray.AddMenuItem("Pause Timer", "Hentikan/Jalankan timer")
	mRestNow := systray.AddMenuItem("", "Mulai istirahat manual")
	mRestDuration := systray.AddMenuItem("", "Atur durasi istirahat")
	mDecreaseRest := mRestDuration.AddSubMenuItem("-5s", "Kurangi durasi istirahat")
	mIncreaseRest := mRestDuration.AddSubMenuItem("+5s", "Tambah durasi istirahat")
	mResetRest := mRestDuration.AddSubMenuItem("Reset ke 20s", "Gunakan durasi istirahat default")
	updateRestMenuTitles(mRestNow, mRestDuration)
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Keluar", "Tutup aplikasi")

	stopChan = make(chan struct{})
	isRunning = true

	go start2020Timer()

	go func() {
		for {
			select {
			case <-mToggle.ClickedCh:
				if isRunning {
					stopChan <- struct{}{}
					isRunning = false
					mToggle.SetTitle("Resume Timer")
					mStatus.SetTitle("Status: Di-pause")
					systray.SetTitle("👁️ Paused")
				} else {
					go start2020Timer()
					isRunning = true
					mToggle.SetTitle("Pause Timer")
					mStatus.SetTitle("Status: Berjalan")
					systray.SetTitle("👁️ 20m")
				}

			case <-mRestNow.ClickedCh:
				go triggerEyeRest()

			case <-mDecreaseRest.ClickedCh:
				adjustRestDuration(-restDurationStep)
				saveConfig()
				updateRestMenuTitles(mRestNow, mRestDuration)

			case <-mIncreaseRest.ClickedCh:
				adjustRestDuration(restDurationStep)
				saveConfig()
				updateRestMenuTitles(mRestNow, mRestDuration)

			case <-mResetRest.ClickedCh:
				setRestDuration(defaultRestDuration)
				saveConfig()
				updateRestMenuTitles(mRestNow, mRestDuration)

			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func start2020Timer() {
	ticker = time.NewTicker(workDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			triggerEyeRest()
		case <-stopChan:
			return
		}
	}
}

func triggerEyeRest() {
	durationSec := restDurationSeconds()

	// 1. Notifikasi Suara & System Toast
	_ = beeep.Notify(
		"Waktunya Istirahat Mata!",
		fmt.Sprintf("Tatap objek sejauh 6 meter selama %s.", formatDuration(time.Duration(durationSec)*time.Second)),
		"",
	)
	_ = beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)

	// 2. Tampilkan prompt istirahat sesuai platform.
	cleanupPrompt := showRestPrompt(durationSec)

	// 3. Countdown di Menu Bar
	for i := durationSec; i > 0; i-- {
		systray.SetTitle(fmt.Sprintf("👀 %ds", i))
		time.Sleep(1 * time.Second)
	}

	if cleanupPrompt != nil {
		cleanupPrompt()
	}

	// 4. Reset Status Menu Bar & Notifikasi Selesai
	systray.SetTitle("👁️ 20m")
	_ = beeep.Notify(
		"Istirahat Selesai",
		"Mata sudah rileks. Selamat mengoding kembali!",
		"",
	)
}

func onExit() {
	// Cleanup saat aplikasi exit
}

type config struct {
	RestDurationSeconds int `json:"rest_duration_seconds"`
}

func adjustRestDuration(delta time.Duration) time.Duration {
	restDurationMu.Lock()
	defer restDurationMu.Unlock()

	restDuration = clampRestDuration(restDuration + delta)
	return restDuration
}

func setRestDuration(duration time.Duration) {
	restDurationMu.Lock()
	defer restDurationMu.Unlock()

	restDuration = clampRestDuration(duration)
}

func restDurationSeconds() int {
	restDurationMu.RLock()
	defer restDurationMu.RUnlock()

	return int(restDuration.Seconds())
}

func clampRestDuration(duration time.Duration) time.Duration {
	if duration < minRestDuration {
		return minRestDuration
	}
	if duration > maxRestDuration {
		return maxRestDuration
	}
	return duration
}

func updateRestMenuTitles(mRestNow, mRestDuration *systray.MenuItem) {
	label := formatDuration(time.Duration(restDurationSeconds()) * time.Second)
	mRestNow.SetTitle("Istirahat Sekarang (" + label + ")")
	mRestDuration.SetTitle("Durasi Istirahat: " + label)
}

func loadConfig() {
	path, err := configPath()
	if err != nil {
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var cfg config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return
	}
	if cfg.RestDurationSeconds > 0 {
		setRestDuration(time.Duration(cfg.RestDurationSeconds) * time.Second)
	}
}

func saveConfig() {
	path, err := configPath()
	if err != nil {
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}

	cfg := config{RestDurationSeconds: restDurationSeconds()}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return
	}

	_ = os.WriteFile(path, data, 0o600)
}

func configPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "EyeRest", "config.json"), nil
}

func formatDuration(duration time.Duration) string {
	totalSec := int(duration.Seconds())
	if totalSec < 60 {
		return fmt.Sprintf("%ds", totalSec)
	}

	minutes := totalSec / 60
	seconds := totalSec % 60
	if seconds == 0 {
		return fmt.Sprintf("%dm", minutes)
	}
	return fmt.Sprintf("%dm%ds", minutes, seconds)
}
