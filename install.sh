```bash
#!/usr/bin/env bash

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_NAME="mate-screensaver-go"
BINARY_PATH="$PROJECT_DIR/$BINARY_NAME"

INSTALL_DIR="/usr/libexec/mate-screensaver"
DESKTOP_DIR="/usr/share/applications/screensavers"

echo "==> Proyecto: $PROJECT_DIR"

echo "==> Compilando $BINARY_NAME..."
cd "$PROJECT_DIR"

go build -o "$BINARY_PATH" ./cmd/screensaver

echo "==> Instalando binario..."
sudo install -m 755 \
    "$BINARY_PATH" \
    "$INSTALL_DIR/$BINARY_NAME"

echo "==> Instalando archivo .desktop..."
sudo install -m 644 \
    "$PROJECT_DIR/mate-screensaver-go.desktop" \
    "$DESKTOP_DIR/mate-screensaver-go.desktop"

echo "==> Actualizando caché de aplicaciones..."
if command -v update-desktop-database >/dev/null 2>&1; then
    sudo update-desktop-database "$DESKTOP_DIR"
fi

echo
echo "Instalación completada."
echo
echo "Binario:"
echo "  $INSTALL_DIR/$BINARY_NAME"
echo
echo "Desktop:"
echo "  $DESKTOP_DIR/mate-screensaver-go.desktop"
```