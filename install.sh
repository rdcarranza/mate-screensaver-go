#!/usr/bin/env bash

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

BINARY_NAME="mate-screensaver-go"
BINARY_PATH="$PROJECT_DIR/$BINARY_NAME"

INSTALL_DIR="/usr/libexec/mate-screensaver"
DESKTOP_DIR="/usr/share/applications/screensavers"

CONFIG_FILE="/etc/mate-screensaver-go.conf"

echo "==> Proyecto: $PROJECT_DIR"
echo

echo "==> Compilando $BINARY_NAME..."
cd "$PROJECT_DIR"

go build -o "$BINARY_PATH" ./cmd/screensaver

echo "==> Instalando binario..."
sudo install -m 0755 \
    "$BINARY_PATH" \
    "$INSTALL_DIR/$BINARY_NAME"

echo "==> Instalando archivo .desktop..."
sudo install -m 0644 \
    "$PROJECT_DIR/mate-screensaver-go.desktop" \
    "$DESKTOP_DIR/mate-screensaver-go.desktop"

echo "==> Configurando ubicación..."

if [[ -f "$CONFIG_FILE" ]]; then
    echo "    Configuración existente preservada:"
    echo "    $CONFIG_FILE"
else
    echo "    Creando configuración:"
    echo "    $CONFIG_FILE"

    sudo tee "$CONFIG_FILE" >/dev/null <<'EOF'
# Configuración de mate-screensaver-go

# Ubicación del dispositivo.
# Dejar vacío para utilizar la ubicación automática.
LATITUDE=
LONGITUDE=
EOF

    sudo chmod 0644 "$CONFIG_FILE"
fi

echo "==> Actualizando caché de aplicaciones..."

if command -v update-desktop-database >/dev/null 2>&1; then
    sudo update-desktop-database "$DESKTOP_DIR"
fi

echo
echo "========================================"
echo " Instalación completada"
echo "========================================"
echo
echo "Binario:"
echo "  $INSTALL_DIR/$BINARY_NAME"
echo
echo "Desktop:"
echo "  $DESKTOP_DIR/mate-screensaver-go.desktop"
echo
echo "Configuración:"
echo "  $CONFIG_FILE"
echo

if [[ -f "$CONFIG_FILE" ]]; then
    echo "Contenido de la configuración:"
    echo
    sudo cat "$CONFIG_FILE"
    echo
fi