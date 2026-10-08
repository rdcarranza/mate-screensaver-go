package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

const DefaultPath = "/etc/mate-screensaver-go.conf"

type Config struct {
	Latitude  string
	Longitude string
}

// Load carga la configuración desde path.
//
// El archivo de configuración es opcional.
// Si no existe, se devuelve una configuración vacía y ningún error.
//
// Las variables soportadas son:
//
//	LATITUDE=-27.8746117
//	LONGITUDE=-63.9869298
func Load(path string) (Config, error) {
	var cfg Config

	file, err := os.Open(path)

	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}

	if err != nil {
		return cfg, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignorar líneas vacías.
		if line == "" {
			continue
		}

		// Ignorar comentarios.
		if strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")

		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch key {
		case "LATITUDE":
			cfg.Latitude = value

		case "LONGITUDE":
			cfg.Longitude = value
		}
	}

	if err := scanner.Err(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// LoadDefault carga la configuración desde la ubicación estándar:
//
//	/etc/mate-screensaver-go.conf
//
// Si el archivo no existe, devuelve una configuración vacía
// y permite continuar normalmente.
func LoadDefault() (Config, error) {
	return Load(DefaultPath)
}
