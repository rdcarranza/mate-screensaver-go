package main

/*
#cgo LDFLAGS: -lX11 -lXss -lXext

#include <X11/Xlib.h>

Display *screen_open(void);
void screen_close(Display *display);

unsigned long screen_idle(Display *display);

void screen_disable_power_management(Display *display);

void screen_hide_cursor(Display *display);
void screen_show_cursor(Display *display);

Window screen_window_create(Display *display);
void screen_window_show(Display *display, Window window);
void screen_window_hide(Display *display, Window window);
void screen_window_destroy(Display *display, Window window);

unsigned long screen_window_width(Display *display);
unsigned long screen_window_height(Display *display);
*/
import "C"

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

const (
	idleTimeout  = 15 * time.Second
	pollInterval = 250 * time.Millisecond

	screensaverPath = "/usr/libexec/mate-screensaver/mate-screensaver-go"
)

type greeter struct {
	display *C.Display
	window  C.Window
	process *exec.Cmd
}

func main() {
	fmt.Println("mate-screensaver-greeter")

	display := C.screen_open()

	if display == nil {
		fmt.Fprintln(os.Stderr, "no se pudo abrir DISPLAY")
		os.Exit(1)
	}

	C.screen_disable_power_management(display)

	defer C.screen_close(display)

	width := C.screen_window_width(display)
	height := C.screen_window_height(display)

	fmt.Printf(
		"Display: %dx%d\n",
		uint64(width),
		uint64(height),
	)

	fmt.Println("Esperando 15 segundos de inactividad...")

	g := &greeter{
		display: display,
	}

	setupSignals(g)

	for {
		if g.process == nil {
			g.waitForIdle()
			g.activate()
		} else {
			g.waitForActivity()
			g.deactivate()
		}
	}
}

func (g *greeter) waitForIdle() {
	for {
		idle := time.Duration(C.screen_idle(g.display)) * time.Millisecond

		if idle >= idleTimeout {
			return
		}

		time.Sleep(pollInterval)
	}
}

func (g *greeter) waitForActivity() {
	for {
		idle := time.Duration(C.screen_idle(g.display)) * time.Millisecond

		if idle < idleTimeout {
			return
		}

		time.Sleep(pollInterval)
	}
}

func (g *greeter) activate() {
	fmt.Println("Activando screensaver...")

	g.window = C.screen_window_create(g.display)

	if g.window == 0 {
		fmt.Fprintln(os.Stderr, "no se pudo crear la ventana")
		return
	}

	C.screen_window_show(
		g.display,
		g.window,
	)

	C.screen_hide_cursor(g.display)

	fmt.Println("Ventana creada y mostrada.")

	windowID := uint64(g.window)

	cmd := exec.Command(screensaverPath)

	cmd.Env = append(
		os.Environ(),
		"HOME=/var/lib/lightdm",
		"USER=lightdm",
		"LOGNAME=lightdm",
		"DISPLAY="+os.Getenv("DISPLAY"),
		"XAUTHORITY="+os.Getenv("XAUTHORITY"),
		"MATE_SCREENSAVER_GREETER=1",
		"XSCREENSAVER_WINDOW=0x"+strconv.FormatUint(windowID, 16),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf(
		"XSCREENSAVER_WINDOW=0x%x\n",
		windowID,
	)

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"no se pudo iniciar screensaver: %v\n",
			err,
		)

		C.screen_window_destroy(
			g.display,
			g.window,
		)

		g.window = 0

		return
	}

	g.process = cmd

	fmt.Printf(
		"mate-screensaver-go iniciado (PID %d)\n",
		cmd.Process.Pid,
	)

	_ = unsafe.Pointer(nil)
}

func (g *greeter) deactivate() {
	fmt.Println("Actividad detectada: desactivando screensaver...")

	if g.process != nil && g.process.Process != nil {
		_ = g.process.Process.Signal(syscall.SIGTERM)

		done := make(chan error, 1)

		go func() {
			done <- g.process.Wait()
		}()

		select {
		case err := <-done:
			if err != nil {
				fmt.Printf(
					"mate-screensaver-go finalizó: %v\n",
					err,
				)
			}

		case <-time.After(2 * time.Second):
			fmt.Println(
				"mate-screensaver-go no terminó; enviando SIGKILL",
			)

			_ = g.process.Process.Kill()
			<-done
		}

		g.process = nil
	}

	if g.window != 0 {
		C.screen_window_hide(
			g.display,
			g.window,
		)

		C.screen_window_destroy(
			g.display,
			g.window,
		)

		g.window = 0
	}

	C.screen_show_cursor(g.display)

	fmt.Println("Screensaver desactivado.")
}

func setupSignals(g *greeter) {
	signals := make(chan os.Signal, 1)

	signal.Notify(
		signals,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	go func() {
		<-signals

		fmt.Println("Finalizando...")

		g.deactivate()

		os.Exit(0)
	}()
}
