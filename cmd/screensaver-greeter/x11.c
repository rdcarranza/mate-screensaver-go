#include <stdio.h>
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/extensions/scrnsaver.h>
#include <X11/extensions/dpms.h>
#include <stdlib.h>

Display *screen_open(void)
{
    return XOpenDisplay(NULL);
}

void screen_close(Display *display)
{
    if (display != NULL) {
        XCloseDisplay(display);
    }
}

unsigned long screen_idle(Display *display)
{
    XScreenSaverInfo *info;
    unsigned long idle;

    if (display == NULL) {
        return 0;
    }

    info = XScreenSaverAllocInfo();

    if (info == NULL) {
        return 0;
    }

    XScreenSaverQueryInfo(
        display,
        DefaultRootWindow(display),
        info
    );

    idle = info->idle;

    XFree(info);

    return idle;
}

Window screen_window_create(Display *display)
{
    Window root;
    Window window;
    XSetWindowAttributes attributes;
    unsigned long mask;

    if (display == NULL) {
        return 0;
    }

    root = DefaultRootWindow(display);

    attributes.override_redirect = True;
    attributes.background_pixel = BlackPixel(display, DefaultScreen(display));
    attributes.border_pixel = BlackPixel(display, DefaultScreen(display));

    mask =
        CWOverrideRedirect |
        CWBackPixel |
        CWBorderPixel;

    window = XCreateWindow(
        display,
        root,
        0,
        0,
        DisplayWidth(display, DefaultScreen(display)),
        DisplayHeight(display, DefaultScreen(display)),
        0,
        DefaultDepth(display, DefaultScreen(display)),
        InputOutput,
        DefaultVisual(display, DefaultScreen(display)),
        mask,
        &attributes
    );

    return window;
}

void screen_window_show(Display *display, Window window)
{
    if (display == NULL || window == 0) {
        return;
    }

    XMapRaised(display, window);
    XFlush(display);
}

void screen_window_hide(Display *display, Window window)
{
    if (display == NULL || window == 0) {
        return;
    }

    XUnmapWindow(display, window);
    XFlush(display);
}

void screen_window_destroy(Display *display, Window window)
{
    if (display == NULL || window == 0) {
        return;
    }

    XDestroyWindow(display, window);
    XFlush(display);
}

unsigned long screen_window_width(Display *display)
{
    if (display == NULL) {
        return 0;
    }

    return DisplayWidth(
        display,
        DefaultScreen(display)
    );
}

unsigned long screen_window_height(Display *display)
{
    if (display == NULL) {
        return 0;
    }

    return DisplayHeight(
        display,
        DefaultScreen(display)
    );
}

void screen_disable_power_management(Display *display)
{
    int event_base;
    int error_base;

    if (display == NULL) {
        return;
    }

    fprintf(
        stderr,
        "DPMS: display válido: %p\n",
        (void *)display
    );

    fprintf(
        stderr,
        "DPMS: desactivando X Screen Saver...\n"
    );

    XSetScreenSaver(
        display,
        0,
        0,
        DontPreferBlanking,
        DefaultExposures
    );

    XSync(display, False);

    fprintf(
        stderr,
        "DPMS: X Screen Saver desactivado.\n"
    );

    if (!DPMSQueryExtension(
        display,
        &event_base,
        &error_base
    )) {
        fprintf(
            stderr,
            "DPMS: extensión DPMS no disponible.\n"
        );

        XFlush(display);

        return;
    }

    fprintf(
        stderr,
        "DPMS: extensión disponible (event=%d error=%d).\n",
        event_base,
        error_base
    );

    fprintf(
        stderr,
        "DPMS: llamando DPMSDisable()...\n"
    );

    DPMSDisable(display);

    XSync(display, False);

    fprintf(
        stderr,
        "DPMS: DPMSDisable() finalizado.\n"
    );

    XFlush(display);
}

void screen_hide_cursor(Display *display)
{
    Cursor invisible_cursor;
    Pixmap bitmap;
    XColor color;
    static char empty[] = { 0 };

    if (display == NULL) {
        return;
    }

    bitmap = XCreateBitmapFromData(
        display,
        DefaultRootWindow(display),
        empty,
        1,
        1
    );

    color.red = 0;
    color.green = 0;
    color.blue = 0;
    color.flags = DoRed | DoGreen | DoBlue;
    color.pixel = 0;

    invisible_cursor = XCreatePixmapCursor(
        display,
        bitmap,
        bitmap,
        &color,
        &color,
        0,
        0
    );

    XDefineCursor(
        display,
        DefaultRootWindow(display),
        invisible_cursor
    );

    XFreeCursor(
        display,
        invisible_cursor
    );

    XFreePixmap(
        display,
        bitmap
    );

    XFlush(display);
}

void screen_show_cursor(Display *display)
{
    if (display == NULL) {
        return;
    }

    XUndefineCursor(
        display,
        DefaultRootWindow(display)
    );

    XFlush(display);
}