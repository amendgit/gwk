
#include "gtkcompat.h"
#include <X11/Xlib.h>
#include <gdk/gdk.h>
#include <gtk/gtk.h>

gboolean g_disableGrab = FALSE;

static gboolean ConfigureTransparentWidget(GtkWidget *widget);
static bool ConfigureOpaqueWidget(GtkWidget *window);
static gboolean ConfigureWidgetTransparency(GtkWidget *window, gboolean transparent);

#if GTK_CHECK_VERSION(3, 0, 0)

struct DeviceGrabContext {
    GdkWindow *window;
    gboolean grabbed;
};

static void GrabMouseDevice(GdkDevice *device, DeviceGrabContext *context);
static void UngrabMouseDevice(GdkDevice *device);

GdkScreen *GdkWindow_GetScreen(GdkWindow *gdkWindow) {
    GdkVisual *gdkVisual = gdk_window_get_visual(gdkWindow);
    return gdk_visual_get_screen(gdkVisual);
}

GdkDisplay *GdkWindow_GetDisplay(GdkWindow *gdkWindow) {
    return gdk_window_get_display(gdkWindow);
}

void GdkWindow_GetSize(GdkWindow *window, gint *w, gint *h) {
    *w = gdk_window_get_width(window);
    *h = gdk_window_get_height(window);
}

gboolean GdkDevice_IsGrabbed(GdkDevice *device) {
    return gdk_display_device_is_grabbed(gdk_display_get_default(), device);
}

void GdkDevice_Ungrab(GdkDevice *device) {
    gdk_device_ungrab(device, GDK_CURRENT_TIME);
}

GdkWindow *GdkDevice_GetWindowAtPosition(GdkDevice *device, gint *x, gint *y) {
    return gdk_device_get_window_at_position(device, x, y);
}

#else /* GTK_CHECK_VERSION(3, 0, 0) */

void GdkWindow_GetSize(GdkWindow *window, gint *w, gint *h) {
    gdk_drawable_get_size(GDK_DRAWABLE(window), w, h);
}

gboolean GdkDevice_IsGrabbed(GdkDevice *device) {
    (void) device;
    return gdk_display_pointer_is_grabbed(gdk_display_get_default());
}

void GdkDevice_Ungrab(GdkDevice *device) {
    (void) device;
    gdk_pointer_ungrab(GDK_CURRENT_TIME);
}

GdkWindow *
GdkDevice_GetWindowAtPosition(GdkDevice *device, gint *x, gint *y) {
    (void) device;
    return gdk_display_get_window_at_pointer(gdk_display_get_default());
}

#endif /* GTK_CHECK_VERSION(3, 0, 0) */
