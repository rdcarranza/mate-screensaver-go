package main

/*
#cgo pkg-config: cairo

#include <cairo.h>
#include <stdlib.h>

static void draw_text(
    cairo_t *cr,
    const char *text,
    const char *font,
    double center_x,
    double y,
    double font_size
)
{
    cairo_set_source_rgb(
        cr,
        1,
        1,
        1
    );

    cairo_select_font_face(
        cr,
        font,
        CAIRO_FONT_SLANT_NORMAL,
        CAIRO_FONT_WEIGHT_NORMAL
    );

    cairo_set_font_size(
        cr,
        font_size
    );

    cairo_text_extents_t extents;

    cairo_text_extents(
        cr,
        text,
        &extents
    );

    double x =
        center_x -
        extents.width / 2.0;

    cairo_move_to(
        cr,
        x,
        y
    );

    cairo_show_text(
        cr,
        text
    );
}
*/
import "C"

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/rdcarranza/mate-screensaver-go/internal/display"
)

func renderScreen(
	cr *C.cairo_t,
	width int,
	height int,
	data display.Data,
) {
	drawBackground(cr)

	clockFontSize := 180.0
	temperatureFontSize := 100.0

	if os.Getenv("MATE_SCREENSAVER_GREETER") == "1" {
		clockFontSize = 220.0
		temperatureFontSize = 80.0
	}

	drawCenteredText(
		cr,
		data.Clock.Time,
		float64(width)/2,  //posición horizontal del dibujo.
		float64(height)/3, //posición vertical del dibujo.
		clockFontSize,
	)

	drawCenteredText(
		cr,
		data.Clock.Date,
		float64(width)/2,
		float64(height)/3+90,
		45,
	)

	drawCenteredText(
		cr,
		data.Location.City+", "+data.Location.Region+", "+data.Location.Country,
		float64(width)/2,
		float64(height)/2+125,
		35,
	)

	drawCenteredText(
		cr,
		fmt.Sprintf(
			"%.1f °C",
			data.Weather.Temperature,
		),
		float64(width)/2,
		float64(height)/2+235,
		temperatureFontSize,
	)

	drawCenteredText(
		cr,
		data.Weather.Description,
		float64(width)/2,
		float64(height)/2+315,
		50,
	)
}

func drawBackground(cr *C.cairo_t) {
	C.cairo_set_source_rgb(
		cr,
		0.02,
		0.02,
		0.02,
	)

	C.cairo_paint(cr)
}

func drawCenteredText(
	cr *C.cairo_t,
	text string,
	centerX float64,
	y float64,
	fontSize float64,
) {
	cText := C.CString(text)
	cFont := C.CString("Sans")

	defer C.free(unsafe.Pointer(cText))
	defer C.free(unsafe.Pointer(cFont))

	C.draw_text(
		cr,
		cText,
		cFont,
		C.double(centerX),
		C.double(y),
		C.double(fontSize),
	)
}
