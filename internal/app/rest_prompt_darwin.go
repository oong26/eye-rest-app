package app

import (
	"fmt"
	"os/exec"
	"time"
)

// showRestPrompt uses macOS JSA so the break prompt can close itself.
func showRestPrompt(durationSec int) func() {
	jsScript := fmt.Sprintf(`
		var app = Application.currentApplication();
		app.includeStandardAdditions = true;

		try {
			app.displayAlert("👀 WAKTUNYA ISTIRAHAT MATA (%s)", {
				message: "Pandanglah objek yang jauh (minimal 6 meter / 20 kaki).\nLayar ini akan menutup otomatis.",
				as: "critical",
				givingUpAfter: %d
			});
		} catch(e) {}
	`, formatDuration(time.Duration(durationSec)*time.Second), durationSec)

	cmd := exec.Command("osascript", "-l", "JavaScript", "-e", jsScript)
	if err := cmd.Start(); err != nil {
		return nil
	}

	return func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	}
}
