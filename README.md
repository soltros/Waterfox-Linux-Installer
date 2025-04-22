# 🦊 Waterfox GUI Installer for Linux

This is a simple **Go application** with a graphical interface (built using [Fyne](https://fyne.io)) that lets you:

- 📂 Select a pre-downloaded Waterfox archive
- 📦 Install Waterfox to `/opt/waterfox`
- 🖇️ Create a symlink (`/usr/local/bin/waterfox`)
- 🖼️ Add a `.desktop` launcher for your app menu
- ❌ Uninstall Waterfox in one click

> 🔐 The install/uninstall steps run with `pkexec`, with full GUI support thanks to environment passthrough for `zenity`.

---

## ✨ Features

- Native GUI built in Go with [Fyne](https://fyne.io)
- Cross-distro compatible (Debian, Arch, Fedora, etc.)
- No need to install Waterfox manually on Linux.
- Uses your system’s GTK dialogs for archive selection
- **No embedded files** — everything is clean and external
- Fully portable binary (after `go build`)

---

## 📦 Requirements

- Go 1.16+
- `zenity`
- `tar`
- `pkexec`

> 💡 On Debian-based systems:
> ```bash
> sudo apt install zenity policykit-1 libgl1-mesa-dev xorg-dev
> ```

---

## 🚀 Usage

### 🛠️ Build

```bash
git clone https://github.com/soltros/Waterfox-Linux-Installer
cd Waterfox-Linux-Installer

go mod tidy
go build -o waterfox-installer
