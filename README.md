# MATE Screensaver Go

<p align="center">
  <strong>Un screensaver moderno para MATE, escrito en Go.</strong>
</p>

<p align="center">
  Reloj · Fecha · Ubicación · Clima
</p>

---

## 📖 Descripción

**MATE Screensaver Go** es un screensaver para el entorno de escritorio **MATE**, desarrollado en **Go**, utilizando **GTK 3**, **GDK** y **Cairo** para la integración gráfica.

El proyecto está pensado como un screensaver sencillo, liviano y extensible, capaz de mostrar información útil directamente sobre la superficie que proporciona `mate-screensaver`.

Actualmente muestra:

- 🕐 Hora actual
- 📅 Fecha
- 📍 Ubicación
- 🌡️ Temperatura
- ☁️ Condiciones meteorológicas

La información meteorológica se obtiene mediante **Open-Meteo**, mientras que la ubicación puede determinarse automáticamente mediante la dirección IP o mediante coordenadas específicas.

---

## ✨ Características

- 🐹 Desarrollado en Go
- 🖥️ Integración nativa con MATE Screensaver
- 🎨 Renderizado mediante Cairo
- 🪟 Utiliza la ventana X11 proporcionada por MATE
- 🕐 Actualización del reloj en tiempo real
- 🌎 Detección automática de ubicación
- 📍 Soporte para coordenadas manuales
- 🌤️ Información meteorológica mediante Open-Meteo
- 🔑 No requiere API key para consultar el clima
- 🧩 Arquitectura modular
- 🐧 Diseñado para Linux/MATE

---

## 🖥️ Captura

> Próximamente.

---

## 🏗️ Arquitectura

El proyecto separa la obtención de información, la preparación de datos y el renderizado gráfico.

```text
                    ┌──────────────────────┐
                    │    MATE Screensaver  │
                    └──────────┬───────────┘
                               │
                    XSCREENSAVER_WINDOW
                               │
                               ▼
                    ┌──────────────────────┐
                    │   GTK / GDK / Cairo  │
                    └──────────┬───────────┘
                               │
                               ▼
                    ┌──────────────────────┐
                    │       Renderer       │
                    └──────────┬───────────┘
                               │
                         display.Data
                               │
             ┌─────────────────┼─────────────────┐
             ▼                 ▼                 ▼
        ┌─────────┐       ┌──────────┐      ┌─────────┐
        │  Clock  │       │ Location │      │ Weather │
        └─────────┘       └──────────┘      └─────────┘
```

Una decisión importante del proyecto es que **MATE es quien crea y administra la superficie del screensaver**.

El programa obtiene el identificador X11 mediante:

```text
XSCREENSAVER_WINDOW
```

y dibuja directamente sobre esa ventana.

Esto evita crear una ventana fullscreen independiente que compita con `mate-screensaver`.

---

## 📂 Estructura del proyecto

```text
mate-screensaver-go/
│
├── cmd/
│   └── screensaver/
│       ├── main.go
│       ├── renderer.go
│       └── gtk.c
│
├── internal/
│   ├── clock/
│   │   └── clock.go
│   │
│   ├── display/
│   │   └── display.go
│   │
│   ├── location/
│   │   └── location.go
│   │
│   └── weather/
│       └── weather.go
│
├── mate-screensaver-go.desktop
├── install.sh
├── go.mod
├── go.sum
└── README.md
```

### `cmd/screensaver`

Contiene el ejecutable principal y la integración con GTK/GDK/Cairo.

### `internal/clock`

Obtiene y formatea la fecha y hora actual.

### `internal/location`

Obtiene la ubicación automáticamente o realiza geocodificación inversa a partir de coordenadas.

### `internal/weather`

Obtiene las condiciones meteorológicas actuales.

### `internal/display`

Agrupa la información que necesita el renderizador.

---

## 🔧 Requisitos

El proyecto requiere:

- Linux
- MATE Desktop
- MATE Screensaver
- Go
- GTK 3
- GDK X11
- Cairo
- X11

En Debian/Ubuntu, los paquetes de desarrollo pueden instalarse con:

```bash
sudo apt install \
    golang \
    libgtk-3-dev \
    libcairo2-dev \
    mate-screensaver
```

El proyecto utiliza CGO para acceder a GTK, GDK y Cairo.

---

## 🚀 Instalación

Cloná el repositorio:

```bash
git clone https://github.com/rdcarranza/mate-screensaver-go.git
cd mate-screensaver-go
```

Ejecutá el instalador:

```bash
chmod +x install.sh
./install.sh
```

El script:

1. Compila el proyecto.
2. Instala el ejecutable.
3. Instala el archivo `.desktop`.
4. Actualiza la base de datos de aplicaciones de escritorio cuando está disponible.

El ejecutable queda instalado en:

```text
/usr/libexec/mate-screensaver/mate-screensaver-go
```

y el descriptor en:

```text
/usr/share/applications/screensavers/mate-screensaver-go.desktop
```

---

## ⚙️ Configuración de MATE

Una vez instalado, el screensaver puede seleccionarse desde:

**Sistema → Preferencias → Apariencia y comportamiento → Salvapantallas**

También puede seleccionarse mediante `gsettings`.

Para comprobar el modo actual:

```bash
gsettings get org.mate.screensaver mode
```

El proyecto utiliza el modo:

```text
'single'
```

Para seleccionar MATE Screensaver Go:

```bash
gsettings set org.mate.screensaver themes \
    "['screensavers-mate-screensaver-go']"
```

---

## 🧪 Ejecución manual

Durante el desarrollo resulta útil ejecutar el programa directamente.

Por ejemplo:

```bash
DISPLAY=:10.0 \
XSCREENSAVER_WINDOW=0x26002BB \
/usr/libexec/mate-screensaver/mate-screensaver-go
```

El valor de `XSCREENSAVER_WINDOW` debe corresponder a una ventana válida creada por MATE.

Normalmente no es necesario establecerlo manualmente: **MATE lo proporciona automáticamente cuando ejecuta el screensaver**.

---

## 🕐 Actualización de información

El diseño previsto separa las distintas frecuencias de actualización:

| Información | Actualización |
|---|---:|
| Hora | 1 segundo |
| Fecha | 1 segundo |
| Ubicación | Al iniciar |
| Clima | Cada 15 minutos |
| Renderizado | Según necesidad |

La idea es evitar realizar solicitudes HTTP durante cada ciclo de renderizado.

---

## 🌎 Ubicación

La ubicación automática se obtiene mediante:

```text
https://ipwho.is/
```

La respuesta proporciona información como:

```text
Ciudad
Provincia
País
Latitud
Longitud
Zona horaria
```

También se permite utilizar coordenadas específicas.

Para obtener el nombre de la ubicación a partir de coordenadas se utiliza:

```text
OpenStreetMap Nominatim
```

---

## 🌤️ Clima

Los datos meteorológicos se obtienen mediante:

```text
https://open-meteo.com/
```

Actualmente se utilizan:

- Temperatura
- Temperatura aparente
- Humedad relativa
- Velocidad del viento
- Código meteorológico
- Indicador día/noche

Los códigos meteorológicos WMO son convertidos a descripciones en español.

Ejemplo:

```text
31.3 °C
Principalmente despejado
```

No se requiere API key.

---

## 🎨 Renderizado

El renderizado se realiza utilizando **Cairo**.

El renderer recibe los datos preparados por la capa de display:

```go
renderScreen(
    cr,
    width,
    height,
    data,
)
```

Esto permite mantener separadas:

```text
Obtención de datos
        ↓
Modelo de display
        ↓
Renderizado
```

El renderer no realiza solicitudes HTTP ni conoce cómo se obtiene la información.

---

## 🪟 Integración con MATE

MATE Screensaver proporciona una ventana X11 para cada ejecución del screensaver.

El programa obtiene su identificador mediante:

```bash
echo "$XSCREENSAVER_WINDOW"
```

Por ejemplo:

```text
0x26002BB
```

El programa convierte ese identificador en un `GdkWindow`:

```c
gdk_x11_window_foreign_new_for_display(
    display,
    window_id
);
```

De esta forma, MATE continúa siendo responsable de:

- crear la superficie;
- posicionarla;
- mostrarla;
- ocultarla;
- controlar el ciclo de vida del screensaver.

Mientras que MATE Screensaver Go se ocupa exclusivamente del contenido visual.

---

## 🛠️ Desarrollo

Para compilar durante el desarrollo:

```bash
go build -o mate-screensaver-go ./cmd/screensaver
```

Para instalar la versión recién compilada:

```bash
./install.sh
```

Para comprobar que el proceso está ejecutándose:

```bash
pgrep -af '/usr/libexec/mate-screensaver/mate-screensaver-go'
```

Para consultar la ventana X11 utilizada:

```bash
PID=$(pgrep -f '/usr/libexec/mate-screensaver/mate-screensaver-go' | head -1)

tr '\0' '\n' < /proc/$PID/environ |
    grep '^XSCREENSAVER_WINDOW='
```

---

## 🔍 Depuración

Durante el desarrollo pueden utilizarse las herramientas estándar de X11:

```bash
xwininfo
```

```bash
xprop
```

```bash
xwininfo -root -tree
```

y herramientas del sistema:

```bash
strace
```

```bash
gdb
```

Por ejemplo:

```bash
sudo strace -f \
    -e trace=write \
    -p "$(pgrep -f '/usr/libexec/mate-screensaver/mate-screensaver-go' | head -1)"
```

Esto permite comprobar que el callback de renderizado está siendo ejecutado.

---



## Integración con LightDM

`mate-screensaver-go` puede integrarse con **LightDM** para mostrar el screensaver sobre la pantalla de login sin utilizar autologin.

### Arquitectura

```text
LightDM
   │
   ▼
greeter-wrapper
   │
   ├── lightdm-gtk-greeter
   │
   └── mate-screensaver-greeter
          │
          ├── detección de inactividad X11
          ├── ventana fullscreen
          ├── gestión del cursor
          ├── control de DPMS
          │
          └── mate-screensaver-go
```

El proyecto se divide en dos componentes:

```text
cmd/
├── screensaver/
│   ├── main.go
│   ├── renderer.go
│   └── gtk.c
│
└── screensaver-greeter/
    ├── main.go
    ├── x11.c
    └── install.sh
```

`mate-screensaver-go` es el renderer y muestra reloj, fecha, ubicación y clima utilizando GTK3, Cairo y X11.

`mate-screensaver-greeter` se encarga de la integración con LightDM: detecta la inactividad, crea la ventana fullscreen, inicia el renderer, oculta el cursor y desactiva el blanking/DPMS.

### Funcionamiento

1. LightDM inicia `lightdm-gtk-greeter`.
2. El wrapper inicia `mate-screensaver-greeter`.
3. Después de **15 segundos de inactividad**, se crea una ventana fullscreen.
4. Se oculta el cursor y se inicia `mate-screensaver-go`.
5. Al detectar actividad de teclado, mouse o pantalla táctil, el renderer termina, se destruye la ventana y se restaura el cursor.
6. Cuando LightDM finaliza el greeter, el wrapper también finaliza el proceso del screensaver.

El display permanece encendido mediante el control directo de X11/DPMS.

### Instalación

El proyecto incluye un instalador automático:

```bash
chmod +x cmd/screensaver-greeter/install.sh
sudo ./cmd/screensaver-greeter/install.sh
```

El instalador:

- compila `mate-screensaver-go`;
- compila `mate-screensaver-greeter`;
- instala ambos binarios en `/usr/libexec/mate-screensaver/`;
- instala el wrapper de LightDM;
- instala la entrada `.desktop`;
- configura `greeter-wrapper` en `/etc/lightdm/lightdm.conf`;
- crea un backup de la configuración de LightDM;
- verifica la configuración resultante.

Para reiniciar LightDM automáticamente:

```bash
sudo ./cmd/screensaver-greeter/install.sh --restart
```

También puede reiniciarse manualmente:

```bash
sudo systemctl restart lightdm
```

> **Importante:** reiniciar LightDM finaliza la sesión gráfica actual. Se recomienda ejecutar este comando desde una TTY (`Ctrl + Alt + F3`) o mediante SSH.

### Archivos instalados

```text
/usr/libexec/mate-screensaver/
├── mate-screensaver-go
└── mate-screensaver-greeter

/usr/libexec/
└── mate-screensaver-greeter-wrapper

/usr/share/applications/
└── mate-screensaver-go.desktop
```

La configuración de LightDM se mantiene en:

```text
/etc/lightdm/lightdm.conf
```

Antes de modificarla, el instalador crea automáticamente un backup:

```text
/etc/lightdm/lightdm.conf.backup-YYYYMMDD-HHMMSS
```

### Configuración

El tiempo de inactividad está definido actualmente en `mate-screensaver-greeter`:

```go
const (
    idleTimeout  = 15 * time.Second
    pollInterval = 250 * time.Millisecond
)
```

Por defecto:

- **15 segundos** de inactividad para activar el screensaver.
- **250 ms** entre comprobaciones de actividad.

### Diagnóstico

Verificar los procesos:

```bash
ps aux | grep -E '[l]ightdm|[m]ate-screensaver'
```

Verificar la configuración de LightDM:

```bash
sudo lightdm --show-config
```

Debe aparecer:

```text
greeter-wrapper=/usr/libexec/mate-screensaver-greeter-wrapper
```

Consultar el log de LightDM:

```bash
sudo journalctl -u lightdm -f
```

Comprobar el acceso X11 del usuario `lightdm`:

```bash
sudo -u lightdm env \
    DISPLAY=:0 \
    XAUTHORITY=/var/lib/lightdm/.Xauthority \
    xdpyinfo | head
```

Comprobar el estado de DPMS:

```bash
sudo -u lightdm env \
    DISPLAY=:0 \
    XAUTHORITY=/var/lib/lightdm/.Xauthority \
    xset q
```

## 🧭 Estado del proyecto

### Implementado

- [x] Aplicación Go
- [x] Integración CGO
- [x] GTK 3
- [x] GDK X11
- [x] Cairo
- [x] Integración con MATE Screensaver
- [x] Ventana `XSCREENSAVER_WINDOW`
- [x] Reloj
- [x] Fecha
- [x] Detección de ubicación
- [x] Geocodificación inversa
- [x] Información meteorológica
- [x] Instalador
- [x] Archivo .desktop

### En desarrollo

- [ ] Actualización robusta del renderizado

- [ ] Gestión completa del ciclo de vida de la ventana

- [ ] Actualización periódica del clima

- [ ] Indicadores gráficos para el clima

- [ ] Mejoras visuales

- [ ] Configuración del usuario

- [ ] Animaciones


### 🤝 Contribuciones

Las contribuciones, ideas y reportes de errores son bienvenidos.

Si encontrás un problema, abrí un issue indicando:

distribución;

versión de MATE;

versión de Go;

versión de GTK;

comportamiento observado;

mensajes de error relevantes.

Los pull requests también son bienvenidos.



## 📄 Licencia

MATE Screensaver Go está distribuido bajo los términos de la:

**GNU Affero General Public License v3.0 (AGPL-3.0)**

Copyright © 2026 Dario Carranza.

El texto completo de la licencia se encuentra en el archivo [`LICENSE`](LICENSE).

La licencia AGPL-3.0 permite utilizar, estudiar, modificar y redistribuir el software, siempre que se respeten las condiciones establecidas por la licencia.