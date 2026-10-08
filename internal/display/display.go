package display

import (
	"github.com/rdcarranza/mate-screensaver-go/internal/clock"
	"github.com/rdcarranza/mate-screensaver-go/internal/location"
	"github.com/rdcarranza/mate-screensaver-go/internal/weather"
)

type Data struct {
	Clock    clock.Info
	Location location.Info
	Weather  weather.Info
}

func NewData() Data {
	return Data{
		Clock: clock.Now(),
	}
}
