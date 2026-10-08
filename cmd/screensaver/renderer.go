package main

/*
#cgo pkg-config: cairo

#include <cairo.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
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

	drawCenteredText(
		cr,
		data.Clock.Time,
		float64(width)/2,
		float64(height)/2,
		150,
	)

	drawCenteredText(
		cr,
		data.Clock.Date,
		float64(width)/2,
		float64(height)/2+90,
		38,
	)
	drawCenteredText(
		cr,
		data.Location.City+", "+data.Location.Region+", "+data.Location.Country,
		float64(width)/2,
		float64(height)/2+145,
		30,
	)

	drawCenteredText(
		cr,
		fmt.Sprintf(
			"%.1f °C",
			data.Weather.Temperature,
		),
		float64(width)/2,
		float64(height)/2+195,
		32,
	)

	drawCenteredText(
		cr,
		data.Weather.Description,
		float64(width)/2,
		float64(height)/2+235,
		26,
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

	C.cairo_set_source_rgb(
		cr,
		1,
		1,
		1,
	)

	C.cairo_select_font_face(
		cr,
		cFont,
		C.CAIRO_FONT_SLANT_NORMAL,
		C.CAIRO_FONT_WEIGHT_NORMAL,
	)

	C.cairo_set_font_size(
		cr,
		C.double(fontSize),
	)

	var extents C.cairo_text_extents_t

	C.cairo_text_extents(
		cr,
		cText,
		&extents,
	)

	x := centerX -
		float64(extents.width)/2

	C.cairo_move_to(
		cr,
		C.double(x),
		C.double(y),
	)

	C.cairo_show_text(
		cr,
		cText,
	)

	// La memoria reservada por C.CString queda pendiente
	// de liberar. Lo resolveremos con un helper C en el
	// siguiente paso.
	_ = unsafe.Pointer(cText)
	_ = unsafe.Pointer(cFont)
}
