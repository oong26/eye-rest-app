package main

import (
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

var (
	workDuration = 20 * time.Minute
	restDuration = 20 * time.Second
	ticker       *time.Ticker
	stopChan     chan struct{}
	isRunning    bool
)

func main() {
	beeep.AppName = appName
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("👁️ 20m")
	systray.SetTooltip(appDisplayName() + " - Pengingat Istirahat Mata 20-20-20")

	mStatus := systray.AddMenuItem("Status: Berjalan", "Status timer")
	mStatus.Disable()
	mVersion := systray.AddMenuItem("Version: "+appVersion, appDisplayName())
	mVersion.Disable()
	systray.AddSeparator()

	mToggle := systray.AddMenuItem("Pause Timer", "Hentikan/Jalankan timer")
	mRestNow := systray.AddMenuItem("Istirahat Sekarang (20s)", "Mulai istirahat manual")
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
	// 1. Notifikasi Suara & System Toast
	_ = beeep.Notify(
		"Waktunya Istirahat Mata!",
		"Tatap objek sejauh 6 meter selama 20 detik.",
		"",
	)
	_ = beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)

	// 2. Tampilkan prompt istirahat sesuai platform.
	cleanupPrompt := showRestPrompt(int(restDuration.Seconds()))

	// 3. Countdown 20 Detik di Menu Bar
	for i := int(restDuration.Seconds()); i > 0; i-- {
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
