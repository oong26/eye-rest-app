# EyeRest

> A lightweight, battery-friendly tray app written in Go to prevent eye strain using the 20-20-20 rule.

## Overview

EyeRest is a minimal status bar/system tray application for people who spend long hours looking at screens. It automates the 20-20-20 rule: every 20 minutes, it prompts you to look at an object 20 feet (6 meters) away for 20 seconds.

Built with Go, EyeRest runs quietly in the background without heavy frameworks like Electron.

Current version: `1.0.0`

## Features

* Native system tray/menu bar app on macOS, Windows, and Linux.
* Dynamic tray indicators: `20m`, `20s`, and paused states.
* 20-second break prompt:
  * macOS uses an auto-closing native `osascript` alert.
  * Windows and Linux use the cross-platform `beeep` alert/notification fallback.
* System notifications and audio cues.
* Manual pause/resume and "rest now" controls from the tray menu.

## Tech Stack

* Language: Go
* System tray: [`github.com/getlantern/systray`](https://github.com/getlantern/systray)
* Notifications and audio cues: [`github.com/gen2brain/beeep`](https://github.com/gen2brain/beeep)
* macOS integration: AppleScript / JavaScript for Automation (`osascript`)

## Prerequisites

* Go 1.20 or later.
* macOS: no extra native packages are required for local builds.
* Windows: local builds work with the Go toolchain.
* Linux: install AppIndicator and GTK development packages before building.

Debian/Ubuntu:

```bash
sudo apt install gcc pkg-config libayatana-appindicator3-dev libgtk-3-dev
```

Fedora:

```bash
sudo dnf install gcc pkgconf-pkg-config libayatana-appindicator-gtk3-devel gtk3-devel
```

## Run Locally

```bash
go run .
```

## Build

Build for your current platform:

```bash
go build -o EyeRest .
```

Build with an explicit semantic version:

```bash
go build -ldflags "-X main.appVersion=1.0.0" -o EyeRest .
```

Build a Windows binary from macOS or Linux:

```bash
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.appVersion=1.0.0" -o EyeRest.exe .
```

Build a Linux binary on Linux:

```bash
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.appVersion=1.0.0" -o EyeRest .
```

Linux tray support uses CGO and native desktop libraries, so cross-compiling Linux binaries from macOS usually requires a Linux C toolchain and the AppIndicator/GTK headers for the target platform. The most reliable path is to build Linux releases on Linux or in a Linux container.
