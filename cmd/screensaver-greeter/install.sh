#!/bin/bash

set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

SCREENSAVER_BINARY="mate-screensaver-go"
GREETER_BINARY="mate-screensaver-greeter"

INSTALL_DIR="/usr/libexec/mate-screensaver"
WRAPPER="/usr/libexec/mate-screensaver-greeter-wrapper"

LIGHTDM_CONFIG="/etc/lightdm/lightdm.conf"
LIGHTDM_BACKUP_DIR="/etc/lightdm"

DESKTOP_FILE="/usr/share/applications/mate-screensaver-go.desktop"

echo
echo "=============================================="
echo " MATE Screensaver Go - Instalador"
echo "=============================================="
echo
echo "Proyecto:"
echo "  $PROJECT_ROOT"
echo
echo "Instalación:"
echo "  $INSTALL_DIR"
echo
echo "LightDM:"
echo "  $LIGHTDM_CONFIG"
echo

# ------------------------------------------------------------
# Comprobar root
# ------------------------------------------------------------

if [[ "${EUID}" -ne 0 ]]; then
    echo "ERROR: este instalador debe ejecutarse como root."
    echo
    echo "Usá:"
    echo "  sudo $0"
    exit 1
fi

# ------------------------------------------------------------
# Comprobar herramientas
# ------------------------------------------------------------

for command in go install cp mkdir chmod mv mktemp; do
    if ! command -v "$command" >/dev/null 2>&1; then
        echo "ERROR: no se encontró el comando: $command"
        exit 1
    fi
done

if [[ ! -f "$PROJECT_ROOT/go.mod" ]]; then
    echo "ERROR: no se encontró go.mod."
    echo "¿El script está dentro del proyecto mate-screensaver-go?"
    exit 1
fi

if [[ ! -f "$LIGHTDM_CONFIG" ]]; then
    echo "ERROR: no existe:"
    echo "  $LIGHTDM_CONFIG"
    exit 1
fi

# ------------------------------------------------------------
# Compilar
# ------------------------------------------------------------

echo "[1/7] Compilando mate-screensaver-go..."

cd "$PROJECT_ROOT"

go build \
    -o "$PROJECT_ROOT/$SCREENSAVER_BINARY" \
    ./cmd/screensaver

echo "      OK"

echo "[2/7] Compilando mate-screensaver-greeter..."

go build \
    -o "$PROJECT_ROOT/$GREETER_BINARY" \
    ./cmd/screensaver-greeter

echo "      OK"

# ------------------------------------------------------------
# Crear directorio de instalación
# ------------------------------------------------------------

echo "[3/7] Instalando binarios..."

mkdir -p "$INSTALL_DIR"

install -m 0755 \
    "$PROJECT_ROOT/$SCREENSAVER_BINARY" \
    "$INSTALL_DIR/$SCREENSAVER_BINARY"

install -m 0755 \
    "$PROJECT_ROOT/$GREETER_BINARY" \
    "$INSTALL_DIR/$GREETER_BINARY"

echo "      $INSTALL_DIR/$SCREENSAVER_BINARY"
echo "      $INSTALL_DIR/$GREETER_BINARY"

# ------------------------------------------------------------
# Instalar wrapper
# ------------------------------------------------------------

echo "[4/7] Instalando wrapper de LightDM..."

cat > "$WRAPPER" <<'EOF'
#!/bin/sh

SCREENSAVER_GREETER="/usr/libexec/mate-screensaver/mate-screensaver-greeter"

"$SCREENSAVER_GREETER" &
SCREENSAVER_PID=$!

"$@" &
GREETER_PID=$!

cleanup()
{
    trap - INT TERM HUP EXIT

    if kill -0 "$GREETER_PID" 2>/dev/null; then
        kill "$GREETER_PID" 2>/dev/null
    fi

    if kill -0 "$SCREENSAVER_PID" 2>/dev/null; then
        kill "$SCREENSAVER_PID" 2>/dev/null
    fi

    wait "$GREETER_PID" 2>/dev/null
    wait "$SCREENSAVER_PID" 2>/dev/null
}

trap cleanup INT TERM HUP EXIT

wait "$GREETER_PID"
STATUS=$?

if kill -0 "$SCREENSAVER_PID" 2>/dev/null; then
    kill "$SCREENSAVER_PID" 2>/dev/null
fi

wait "$SCREENSAVER_PID" 2>/dev/null

trap - INT TERM HUP EXIT

exit "$STATUS"
EOF

chmod 0755 "$WRAPPER"

echo "      $WRAPPER"

# ------------------------------------------------------------
# Instalar desktop entry
# ------------------------------------------------------------

echo "[5/7] Instalando entrada de escritorio..."

cat > "$DESKTOP_FILE" <<'EOF'
[Desktop Entry]
Name=MATE Screensaver Go
Comment=Reloj, ubicación y clima
Exec=/usr/libexec/mate-screensaver/mate-screensaver-go
TryExec=/usr/libexec/mate-screensaver/mate-screensaver-go
StartupNotify=false
Terminal=false
Type=Application
Categories=Screensaver;
Keywords=MATE;screensaver;clock;weather;location;
OnlyShowIn=MATE;
EOF

chmod 0644 "$DESKTOP_FILE"

echo "      $DESKTOP_FILE"

# ------------------------------------------------------------
# Backup LightDM
# ------------------------------------------------------------

echo "[6/7] Configurando LightDM..."

TIMESTAMP="$(date +%Y%m%d-%H%M%S)"
BACKUP="$LIGHTDM_BACKUP_DIR/lightdm.conf.backup-$TIMESTAMP"

cp "$LIGHTDM_CONFIG" "$BACKUP"

echo "      Backup:"
echo "      $BACKUP"

# ------------------------------------------------------------
# Modificar [Seat:*]
# ------------------------------------------------------------

TMP_CONFIG="$(mktemp)"

awk '
BEGIN {
    in_seat = 0
    inserted = 0
}

function print_option() {
    print "greeter-wrapper=/usr/libexec/mate-screensaver-greeter-wrapper"
    inserted = 1
}

# Nueva sección
/^[[:space:]]*\[/ {
    if (in_seat && !inserted) {
        print_option()
    }

    in_seat = ($0 ~ /^\[Seat:\*\][[:space:]]*$/)

    print
    next
}

# Dentro de [Seat:*], reemplazar cualquier greeter-wrapper
in_seat && /^[[:space:]]*greeter-wrapper[[:space:]]*=/ {
    if (!inserted) {
        print_option()
    }

    next
}

{
    print
}

END {
    if (in_seat && !inserted) {
        print_option()
    }
}
' "$LIGHTDM_CONFIG" > "$TMP_CONFIG"

mv "$TMP_CONFIG" "$LIGHTDM_CONFIG"

# ------------------------------------------------------------
# Verificar configuración
# ------------------------------------------------------------

echo
echo "      Configuración resultante:"
echo

lightdm --show-config

echo

if ! lightdm --show-config | grep -q \
    '^B  greeter-wrapper=/usr/libexec/mate-screensaver-greeter-wrapper$'
then
    echo "ERROR: LightDM no está tomando greeter-wrapper."
    echo
    echo "Restaurando backup:"
    echo "  $BACKUP"

    cp "$BACKUP" "$LIGHTDM_CONFIG"

    exit 1
fi

# ------------------------------------------------------------
# Verificar que NO se agregó autologin-user
# ------------------------------------------------------------

if grep -Eq '^[[:space:]]*autologin-user[[:space:]]*=' "$LIGHTDM_CONFIG"; then
    echo
    echo "AVISO:"
    echo "La configuración contiene autologin-user."
    echo
    echo "El instalador no agregó esta opción."
    echo "Revisá la configuración si no querés autologin."
fi

# ------------------------------------------------------------
# Final
# ------------------------------------------------------------

echo
echo "=============================================="
echo " Instalación completada correctamente"
echo "=============================================="
echo
echo "Binarios:"
echo "  $INSTALL_DIR/$SCREENSAVER_BINARY"
echo "  $INSTALL_DIR/$GREETER_BINARY"
echo
echo "Wrapper:"
echo "  $WRAPPER"
echo
echo "Desktop entry:"
echo "  $DESKTOP_FILE"
echo
echo "Backup de LightDM:"
echo "  $BACKUP"
echo

# ------------------------------------------------------------
# Reinicio opcional
# ------------------------------------------------------------

if [[ "${1:-}" == "--restart" ]]; then
    echo "Reiniciando LightDM..."
    echo
    systemctl restart lightdm
else
    echo "LightDM NO fue reiniciado automáticamente."
    echo
    echo "Para aplicar la integración:"
    echo
    echo "  sudo systemctl restart lightdm"
    echo
    echo "O ejecutá directamente:"
    echo
    echo "  sudo $0 --restart"
fi

echo