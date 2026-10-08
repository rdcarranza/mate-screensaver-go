#include <gtk/gtk.h>
#include <cairo.h>
#include <gdk/gdkx.h>
#include <gdk/x11/gdkx11window.h>
#include <X11/Xlib.h>

#include <stdlib.h>
#include <stdint.h>

extern void draw_screen(
    cairo_t *cr,
    int width,
    int height
);

void screen_init(void)
{
    int argc = 0;
    char **argv = NULL;

    gtk_init(&argc, &argv);
}

void *screen_window_create(void)
{
    const char *window_id_string;

    window_id_string = getenv("XSCREENSAVER_WINDOW");

    if (window_id_string == NULL) {
        g_printerr(
            "XSCREENSAVER_WINDOW no está definido\n"
        );

        return NULL;
    }

    Window window_id;

    window_id = (Window)strtoull(
        window_id_string,
        NULL,
        0
    );

    GdkDisplay *display;

    display = gdk_display_get_default();

    if (display == NULL) {
        g_printerr(
            "No se pudo obtener el GdkDisplay\n"
        );

        return NULL;
    }

    GdkWindow *window;

    window = gdk_x11_window_foreign_new_for_display(
        display,
        window_id
    );

    if (window == NULL) {
        g_printerr(
            "No se pudo obtener GdkWindow para XID %s\n",
            window_id_string
        );

        return NULL;
    }

    g_print(
        "XSCREENSAVER_WINDOW=%s\n",
        window_id_string
    );

    g_print(
        "GdkWindow obtenido correctamente\n"
    );

    return window;
}

void screen_window_show(void *window)
{
    GdkWindow *gdk_window;

    gdk_window = (GdkWindow *)window;

    gdk_window_show(gdk_window);
}

/*
static gboolean draw_screen_idle(gpointer data)
{
    GdkWindow *gdk_window;
    cairo_t *cr;

    gdk_window = (GdkWindow *)data;

    cr = gdk_cairo_create(gdk_window);
    
    
    draw_screen(
        cr,
        gdk_window_get_width(gdk_window) ,
        gdk_window_get_height(gdk_window)
    );

    cairo_destroy(cr);

    return G_SOURCE_REMOVE;
}
*/

static gboolean draw_screen_idle(gpointer data)
{
    GdkWindow *gdk_window;
    GdkDrawingContext *drawing_context;
    cairo_t *cr;
    cairo_region_t *region;

    gdk_window = (GdkWindow *)data;

    region = cairo_region_create_rectangle(
        &(GdkRectangle){
            0,
            0,
            gdk_window_get_width(gdk_window),
            gdk_window_get_height(gdk_window)
        }
    );

    drawing_context = gdk_window_begin_draw_frame(
        gdk_window,
        region
    );

    if (drawing_context == NULL) {
        cairo_region_destroy(region);

        return G_SOURCE_REMOVE;
    }

    cr = gdk_drawing_context_get_cairo_context(
        drawing_context
    );

    draw_screen(
        cr,
        gdk_window_get_width(gdk_window),
        gdk_window_get_height(gdk_window)
    );

    gdk_window_end_draw_frame(
        gdk_window,
        drawing_context
    );

    cairo_region_destroy(region);

    return G_SOURCE_REMOVE;
}


void screen_queue_draw(void *window)
{
    g_main_context_invoke(
        NULL,
        draw_screen_idle,
        window
    );
}

void screen_main(void)
{
    gtk_main();
}

