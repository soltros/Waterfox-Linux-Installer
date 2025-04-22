package main

import (
	"image/png"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const downloadURL = "https://www.waterfox.net/download/"

const bashScript = `#!/bin/bash

set -e

INSTALL_DIR="/opt/waterfox"
BIN_LINK="/usr/local/bin/waterfox"
DESKTOP_FILE="/usr/share/applications/waterfox.desktop"

if [ "$1" = "--uninstall" ]; then
    echo "Uninstalling Waterfox..."
    rm -rf "$INSTALL_DIR"
    rm -f "$BIN_LINK"
    rm -f "$DESKTOP_FILE"
    echo "✅ Waterfox removed."
    exit 0
fi

ARCHIVE=$(zenity --file-selection \
    --title="Select Waterfox Archive" \
    --filename="$HOME/Downloads/" \
    --file-filter="Waterfox tarball (waterfox*.tar.*) | waterfox*.tar.*")

if [ -z "$ARCHIVE" ]; then
    zenity --error --title="No File Selected" --text="You must select a Waterfox archive to install."
    exit 1
fi

(
    echo "10"
    echo "# Preparing install directory..."
    rm -rf "$INSTALL_DIR"
    mkdir -p "$INSTALL_DIR"
    sleep 0.5

    echo "40"
    echo "# Extracting Waterfox archive..."
    tar -xf "$ARCHIVE" -C "$INSTALL_DIR" --strip-components=1
    sleep 1

    echo "60"
    echo "# Creating symlink..."
    ln -sf "$INSTALL_DIR/waterfox" "$BIN_LINK"
    sleep 0.5

    echo "80"
    echo "# Creating desktop launcher..."
    cat > "$DESKTOP_FILE" <<EOF
[Desktop Entry]
Name=Waterfox
Exec=$BIN_LINK %u
Icon=$INSTALL_DIR/browser/chrome/icons/default/default128.png
Type=Application
Categories=Network;WebBrowser;
MimeType=text/html;text/xml;application/xhtml+xml;application/xml;x-scheme-handler/http;x-scheme-handler/https;
StartupNotify=true
EOF

    chmod +x "$DESKTOP_FILE"
    sleep 0.5

    echo "100"
    echo "# Installation complete!"
) | zenity --progress --title="Installing Waterfox" --percentage=0 --auto-close

zenity --info --title="Success" --text="Waterfox was installed successfully!"
`

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Waterfox Installer")

	// Load image from assets folder
	f, err := os.Open("assets/waterfox.png")
	if err != nil {
		log.Fatalf("Failed to open logo image: %v", err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		log.Fatalf("Failed to decode logo image: %v", err)
	}

	logo := canvas.NewImageFromImage(img)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(128, 128))

	info := widget.NewLabel("Waterfox is a privacy-focused browser based on Firefox.\nClick below to visit the official download page or begin installation.")
	link := widget.NewHyperlink("Go to Waterfox Download Page", parseURL(downloadURL))

	installBtn := widget.NewButton("Install Waterfox", func() {
		runBashScript(false)
	})
	uninstallBtn := widget.NewButton("Uninstall Waterfox", func() {
		runBashScript(true)
	})

	content := container.NewVBox(
		logo,
		info,
		link,
		installBtn,
		uninstallBtn,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(440, 300))
	myWindow.ShowAndRun()
}

func parseURL(s string) *url.URL {
	u, err := url.Parse(downloadURL)
	if err != nil {
		log.Fatal("Invalid URL:", err)
	}
	return u
}

func runBashScript(uninstall bool) {
	tmpFile := filepath.Join(os.TempDir(), "waterfox-install.sh")
	err := os.WriteFile(tmpFile, []byte(bashScript), 0755)
	if err != nil {
		log.Println("Failed to write bash script:", err)
		return
	}

	args := []string{tmpFile}
	if uninstall {
		args = append(args, "--uninstall")
	}

	// Use pkexec and preserve DISPLAY/XAUTHORITY for zenity dialogs to work
	cmd := exec.Command("pkexec", "env",
		"DISPLAY="+os.Getenv("DISPLAY"),
		"XAUTHORITY="+os.Getenv("XAUTHORITY"),
		args[0],
	)
	cmd.Args = append(cmd.Args, args[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Println("Script execution failed:", err)
	}
}
