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
- 📍 Soporte para coordenadas configurables
- 🌤️ Información meteorológica mediante Open-Meteo
- 🔑 No requiere API key para consultar el clima
- 🧩 Arquitectura modular
- 🐧 Diseñado para Linux/MATE
- 🔐 Integración opcional con LightDM

---

## 🖥️ Captura

<img width="1600" height="900" alt="remmina_rp5-s3_192 168 1 253_20261008-220708" src="https://github.com/user-attachments/assets/91ae823c-4fa3-4acf-91ad-e67323775b0c" />


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
│   ├── screensaver/
│   │   ├── main.go
│   │   ├── renderer.go
│   │   └── gtk.c
│   │
│   └── screensaver-greeter/
│       ├── main.go
│       └── x11.c
│
├── internal/
│   ├── clock/
│   │   └── clock.go
│   │
│   ├── config/
│   │   └── config.go
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
├── LICENSE
└── README.md
```

### `cmd/screensaver`

Contiene el ejecutable principal y la integración con GTK/GDK/Cairo.

Es responsable del renderizado y de iniciar las actualizaciones de información.

### `cmd/screensaver-greeter`

Contiene la integración específica con LightDM.

Detecta la inactividad, crea la ventana fullscreen, controla el cursor y administra el proceso de `mate-screensaver-go`.

### `internal/clock`

Obtiene y formatea la fecha y hora actual.

### `internal/config`

Carga la configuración opcional del sistema desde:

```text
/etc/mate-screensaver-go.conf
```

La configuración es independiente del renderer y puede ser utilizada tanto por la ejecución normal como por la integración con LightDM.

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
4. Crea la configuración opcional si todavía no existe.
5. Actualiza la base de datos de aplicaciones de escritorio cuando está disponible.

El ejecutable queda instalado en:

```text
/usr/libexec/mate-screensaver/mate-screensaver-go
```

El descriptor queda instalado en:

```text
/usr/share/applications/screensavers/mate-screensaver-go.desktop
```

La configuración queda en:

```text
/etc/mate-screensaver-go.conf
```

### Preservación de la configuración

El instalador **no sobrescribe** un archivo de configuración existente.

Esto permite ejecutar nuevamente:

```bash
./install.sh
```

para actualizar el programa sin perder la configuración de ubicación.

---

## ⚙️ Configuración

MATE Screensaver Go utiliza un archivo de configuración opcional:

```text
/etc/mate-screensaver-go.conf
```

Un archivo de configuración básico tiene esta estructura:

```ini
# Configuración de mate-screensaver-go

# Ubicación del dispositivo.
# Dejar vacío para utilizar la ubicación automática.
LATITUDE=
LONGITUDE=
```

Para especificar una ubicación:

```ini
LATITUDE=-27.8746117
LONGITUDE=-63.9869298
```

### Comportamiento de la configuración

El archivo de configuración es **opcional**.

Si el archivo:

- no existe;
- existe pero `LATITUDE` está vacío;
- existe pero `LONGITUDE` está vacío;

el programa continúa normalmente y utiliza la detección automática de ubicación.

Por ejemplo, si el archivo no existe:

```bash
sudo rm /etc/mate-screensaver-go.conf
```

MATE Screensaver Go continuará utilizando:

```text
location.Now()
```

No es necesario crear manualmente el archivo para que el screensaver funcione.

### Coordenadas configuradas

Cuando ambas coordenadas están configuradas, el programa utiliza:

```text
location.FromCoordinates()
```

para determinar la ubicación correspondiente.

El archivo `.desktop` **no contiene las coordenadas**. Esto permite utilizar una única configuración independientemente de cómo se ejecute el screensaver.

---

## 🖥️ Configuración de MATE

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

Normalmente no es necesario establecerlo manualmente:

**MATE lo proporciona automáticamente cuando ejecuta el screensaver.**

---

## 🕐 Actualización de información

El diseño separa las distintas frecuencias de actualización:

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

También se permite utilizar coordenadas específicas mediante:

```text
/etc/mate-screensaver-go.conf
```

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

El renderer no realiza solicitudes HTTP directamente ni necesita conocer cómo se obtiene la información.

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

# 🔐 Integración con LightDM

`mate-screensaver-go` puede integrarse con **LightDM** para mostrar el screensaver sobre la pantalla de login sin utilizar autologin.

## Arquitectura

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
                │
                └── /etc/mate-screensaver-go.conf
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
    └── x11.c
```

`mate-screensaver-go` es el renderer y muestra reloj, fecha, ubicación y clima utilizando GTK3, Cairo y X11.

`mate-screensaver-greeter` se encarga de la integración con LightDM: detecta la inactividad, crea la ventana fullscreen, inicia el renderer, oculta el cursor y desactiva el blanking/DPMS.

### Funcionamiento

1. LightDM inicia `lightdm-gtk-greeter`.
2. El wrapper inicia `mate-screensaver-greeter`.
3. Después de **15 segundos de inactividad**, se crea una ventana fullscreen.
4. Se oculta el cursor y se inicia `mate-screensaver-go`.
5. `mate-screensaver-go` carga la misma configuración `/etc/mate-screensaver-go.conf`.
6. Al detectar actividad de teclado, mouse o pantalla táctil, el renderer termina.
7. Se destruye la ventana y se restaura el cursor.
8. Cuando LightDM finaliza el greeter, el wrapper también finaliza el proceso del screensaver.

El display permanece encendido mediante el control directo de X11/DPMS.

### Configuración de ubicación en LightDM

El greeter **no necesita conocer la ubicación**.

`mate-screensaver-greeter` simplemente inicia:

```text
mate-screensaver-go
```

El propio renderer carga:

```text
/etc/mate-screensaver-go.conf
```

Esto garantiza que la ejecución normal y la ejecución sobre el login de LightDM utilicen exactamente la misma configuración.

---

## 📦 Instalación de la integración con LightDM

La integración con LightDM se instala mediante:

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

---

## 📁 Archivos instalados

Los componentes principales quedan instalados en:

```text
/usr/libexec/mate-screensaver/
├── mate-screensaver-go
└── mate-screensaver-greeter
```

Wrapper:

```text
/usr/libexec/
└── mate-screensaver-greeter-wrapper
```

Descriptor:

```text
/usr/share/applications/screensavers/
└── mate-screensaver-go.desktop
```

Configuración:

```text
/etc/
└── mate-screensaver-go.conf
```

Configuración de LightDM:

```text
/etc/lightdm/lightdm.conf
```

Antes de modificarla, el instalador crea automáticamente un backup:

```text
/etc/lightdm/lightdm.conf.backup-YYYYMMDD-HHMMSS
```

---

## ⚙️ Configuración del greeter

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

El cursor se oculta mientras el screensaver está activo y se restaura al detectar actividad.

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

Y herramientas del sistema:

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

### Diagnóstico de LightDM

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

---

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
- [x] Detección automática de ubicación
- [x] Geocodificación inversa
- [x] Coordenadas configurables
- [x] Configuración común `/etc/mate-screensaver-go.conf`
- [x] Configuración opcional
- [x] Información meteorológica
- [x] Actualización periódica del clima
- [x] Instalador
- [x] Archivo `.desktop`
- [x] Integración con LightDM
- [x] Detección de inactividad
- [x] Ocultamiento/restauración del cursor
- [x] Control de DPMS
- [x] Limpieza del proceso al finalizar LightDM

### En desarrollo

- [ ] Mejoras visuales
- [ ] Animaciones
- [ ] Configuración adicional del usuario

---

## 🤝 Contribuciones

Las contribuciones, ideas y reportes de errores son bienvenidos.

Si encontrás un problema, abrí un issue indicando:

- distribución;
- versión de MATE;
- versión de Go;
- versión de GTK;
- comportamiento observado;
- mensajes de error relevantes.

Los pull requests también son bienvenidos.

---

## 📄 Licencia

MATE Screensaver Go está distribuido bajo los términos de la:

**GNU Affero General Public License v3.0 (AGPL-3.0)**

Copyright © 2026 Dario Carranza.

El texto completo de la licencia se encuentra en el archivo [`LICENSE`](LICENSE).

La licencia AGPL-3.0 permite utilizar, estudiar, modificar y redistribuir el software, siempre que se respeten las condiciones establecidas por la licencia.
