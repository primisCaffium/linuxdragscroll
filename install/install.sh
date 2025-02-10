#!/bin/bash

source .env

cd ../src/ || exit
go get
go build -o linuxdragscroll
chmod +x linuxdragscroll
mv linuxdragscroll "$INSTALL_PATH"

DESKTOP_FILE_BASE="[Desktop Entry]
Name=Linux Drag Scroll
Type=Application
Terminal=false
Comment=Launch the Linux Drag Scroll application
Icon=/usr/share/icons/gnome/48x48/categories/gnome-system.png
Categories=Utility;
X-GNOME-Autostart-enabled=true
"

AUTOSTART_PATH="$HOME/.config/autostart/linuxdragscroll.desktop"
DESKTOP_FILE_PATH="$HOME/.local/share/applications/linuxdragscroll.desktop"

# The sleep is to wait for the evtest device after login
echo "${DESKTOP_FILE_BASE}Exec=sh -c \"sleep 60; $INSTALL_PATH/linuxdragscroll\"" > "$AUTOSTART_PATH"

echo "${DESKTOP_FILE_BASE}Exec=$INSTALL_PATH/linuxdragscroll" > "$DESKTOP_FILE_PATH"

echo "Done"