package main

/*
#cgo pkg-config: gtk+-3.0 cairo

#include <gtk/gtk.h>
#include <cairo.h>

void screen_init(void);
void *screen_window_create(void);
void screen_window_show(void *window);
void screen_queue_draw(void *window);
void screen_main(void);
*/
import "C"

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rdcarranza/mate-screensaver-go/internal/config"
	"github.com/rdcarranza/mate-screensaver-go/internal/display"
	"github.com/rdcarranza/mate-screensaver-go/internal/location"
	"github.com/rdcarranza/mate-screensaver-go/internal/weather"
)

var currentLocation location.Info
var currentWeather weather.Info
var dataMutex sync.RWMutex

//export draw_screen
func draw_screen(
	cr *C.cairo_t,
	width C.int,
	height C.int,
) {
	println("draw_screen:", int(width), int(height))

	data := display.NewData()

	dataMutex.RLock()
	data.Location = currentLocation
	data.Weather = currentWeather
	dataMutex.RUnlock()

	renderScreen(
		cr,
		int(width),
		int(height),
		data,
	)
}

func main() {
	var err error

	/*
		Prioridad para determinar la ubicación:

		1. Coordenadas recibidas como argumento.
		2. /etc/mate-screensaver-go.conf
		3. Ubicación automática mediante location.Now().
	*/

	var coordinates string

	// Mantener compatibilidad con la ejecución mediante argumento:
	//
	// mate-screensaver-go "-27.8746117,-63.9869298"
	if len(os.Args) > 1 {
		coordinates = strings.TrimSpace(os.Args[1])
	}

	// Si no recibimos coordenadas por argumento, intentar cargar
	// la configuración común.
	if coordinates == "" {
		cfg, configErr := config.LoadDefault()

		if configErr != nil {
			fmt.Println(
				"Error leyendo configuración:",
				configErr,
			)

			// La configuración es opcional.
			// Continuamos con ubicación automática.
		} else if cfg.Latitude != "" && cfg.Longitude != "" {
			coordinates = cfg.Latitude + "," + cfg.Longitude

			fmt.Printf(
				"Coordenadas desde configuración: %s\n",
				coordinates,
			)
		}
	}

	// Si tenemos coordenadas, intentamos utilizarlas.
	if coordinates != "" {
		latitude, longitude, parseErr := parseCoordinates(coordinates)

		if parseErr == nil {
			currentLocation, err = location.FromCoordinates(
				latitude,
				longitude,
			)

			if err != nil {
				fmt.Println(
					"Error obteniendo ubicación desde coordenadas:",
					err,
				)

				// No hacemos fallar el screensaver.
				// Intentamos ubicación automática.
				currentLocation, err = location.Now()

				if err != nil {
					fmt.Println(
						"Error obteniendo ubicación automática:",
						err,
					)
					return
				}
			} else {
				fmt.Printf(
					"Ubicación por coordenadas: %s, %s, %s\n",
					currentLocation.City,
					currentLocation.Region,
					currentLocation.Country,
				)
			}
		} else {
			fmt.Printf(
				"Coordenadas inválidas (%s), usando ubicación automática\n",
				parseErr,
			)

			currentLocation, err = location.Now()

			if err != nil {
				fmt.Println(
					"Error obteniendo ubicación:",
					err,
				)
				return
			}
		}
	} else {
		// No hay coordenadas configuradas.
		// Usamos la ubicación automática existente.
		fmt.Println(
			"No hay coordenadas configuradas, usando ubicación automática",
		)

		currentLocation, err = location.Now()

		if err != nil {
			fmt.Println(
				"Error obteniendo ubicación:",
				err,
			)
			return
		}
	}

	currentWeather, err = weather.Now(
		currentLocation.Latitude,
		currentLocation.Longitude,
	)

	if err != nil {
		fmt.Println(
			"No se pudo obtener el clima:",
			err,
		)
	}

	C.screen_init()

	window := C.screen_window_create()

	C.screen_window_show(window)

	go refreshWeather()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			C.screen_queue_draw(window)
		}
	}()

	C.screen_main()
}

func refreshWeather() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		dataMutex.RLock()

		latitude := currentLocation.Latitude
		longitude := currentLocation.Longitude

		dataMutex.RUnlock()

		current, err := weather.Now(
			latitude,
			longitude,
		)

		if err != nil {
			fmt.Println(
				"Error actualizando clima:",
				err,
			)
			continue
		}

		dataMutex.Lock()

		currentWeather = current

		dataMutex.Unlock()

		fmt.Printf(
			"Clima actualizado: %.1f °C - %s\n",
			current.Temperature,
			current.Description,
		)
	}
}

func parseCoordinates(value string) (float64, float64, error) {
	parts := strings.Split(value, ",")

	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("formato inválido")
	}

	latitude, err := strconv.ParseFloat(
		strings.TrimSpace(parts[0]),
		64,
	)

	if err != nil {
		return 0, 0, fmt.Errorf("latitud inválida")
	}

	longitude, err := strconv.ParseFloat(
		strings.TrimSpace(parts[1]),
		64,
	)

	if err != nil {
		return 0, 0, fmt.Errorf("longitud inválida")
	}

	if latitude < -90 || latitude > 90 {
		return 0, 0, fmt.Errorf(
			"latitud fuera de rango",
		)
	}

	if longitude < -180 || longitude > 180 {
		return 0, 0, fmt.Errorf(
			"longitud fuera de rango",
		)
	}

	return latitude, longitude, nil
}
