
#include "gtkcompat_linux.h"
#include "general_linux.h"
#include <X11/Xlib.h>
#include <gdk/gdk.h>
#include <gtk/gtk.h>

gboolean g_disableGrab = FALSE;

static gboolean ConfigureTransparentWindow(GtkWidget *widget);
static void ConfigureOpaqueWindow(GtkWidget *window);
static gboolean ConfigureWindowTransparency(GtkWidget *window, gboolean transparent);

// -----------------------------------------------------------------------------
#if GTK_CHECK_VERSION(3, 0, 0)

typedef struct {
    GdkWindow * window;
    gboolean grabbed;
} DeviceGrabContext;

static void GrabMouseDevice(GdkDevice *device, DeviceGrabContext *context);
static void UngrabMouseDevice(GdkDevice *device);

GdkScreen *
GdkWindowGetScreen(GdkWindow * gdkWindow) {
    GdkVisual * gdkVisual = gdk_window_get_visual(gdkWindow);
    return gdk_visual_get_screen(gdkVisual);
}

GdkDisplay *
GdkWindowGetDisplay(GdkWindow * gdkWindow) {
    return gdk_window_get_display(gdkWindow);
}


gboolean
GdkMouseDevicesGrab(GdkWindow *gdkWindow) {
    if (g_disableGrab) {
        return TRUE;
    }

    DeviceGrabContext context;
    GList *devices = gdk_device_manager_list_devices(
                         gdk_display_get_device_manager(
                             gdk_display_get_default()),
                             GDK_DEVICE_TYPE_MASTER);

    context.window = gdkWindow;
    context.grabbed = FALSE;
    g_list_foreach(devices, (GFunc) GrabMouseDevice, &context);

    return context.grabbed;
}

gboolean
GdkMouseDevicesGrabWithCursor(GdkWindow *gdkWindow, GdkCursor *cursor) {
    return GdkMouseDevicesGrabWithCursor(gdkWindow, cursor, TRUE);
}

gboolean
GdkMouseDevicesGrabWithCursor(GdkWindow *gdkWindow, GdkCursor *cursor, gboolean owner_events) {
    // if (g_disableGrab) {
    //     return TRUE;
    // }
    // GdkGrabStatus status = gdk_device_grab(gdkWindow, owner_events, (GdkEventMask)
    //                                         (GDK_POINTER_MOTION_MASK
    //                                             | GDK_POINTER_MOTION_HINT_MASK
    //                                             | GDK_BUTTON_MOTION_MASK
    //                                             | GDK_BUTTON1_MOTION_MASK
    //                                             | GDK_BUTTON2_MOTION_MASK
    //                                             | GDK_BUTTON3_MOTION_MASK
    //                                             | GDK_BUTTON_PRESS_MASK
    //                                             | GDK_BUTTON_RELEASE_MASK),
    //                                         NULL, cursor, GDK_CURRENT_TIME);
    //
    // return (status == GDK_GRAB_SUCCESS) ? TRUE : FALSE;
    // TOIMPL
}

void
GdkMouseDevicesUngrab() {
    GList *devices = gdk_device_manager_list_devices(
                         gdk_display_get_device_manager(
                             gdk_display_get_default()),
                             GDK_DEVICE_TYPE_MASTER);
    g_list_foreach(devices, (GFunc) UngrabMouseDevice, NULL);
}

void
GdkMasterPointerGrab(GdkWindow *window, GdkCursor *cursor) {
    if (g_disableGrab) {
        gdk_window_set_cursor(window, cursor);
        return;
    }
    gdk_device_grab(gdk_device_manager_get_client_pointer(
                        gdk_display_get_device_manager(
                            gdk_display_get_default())),
                    window, GDK_OWNERSHIP_NONE, FALSE, GDK_ALL_EVENTS_MASK,
                    cursor, GDK_CURRENT_TIME);
}

void
GdkMasterPointerUngrab() {
    gdk_device_ungrab(gdk_device_manager_get_client_pointer(
                          gdk_display_get_device_manager(
                              gdk_display_get_default())),
                      GDK_CURRENT_TIME);
}

void
GdkMasterPointerGetPosition(gint *x, gint *y) {
    gdk_device_get_position(gdk_device_manager_get_client_pointer(
                                gdk_display_get_device_manager(
                                    gdk_display_get_default())),
                            NULL, x, y);
}

gboolean
GdkDeviceIsGrabbed(GdkDevice *device) {
    return gdk_display_device_is_grabbed(gdk_display_get_default(), device);
}

void
GdkDeviceUngrab(GdkDevice *device) {
    gdk_device_ungrab(device, GDK_CURRENT_TIME);
}

GdkWindow *
GdkDeviceGetWindowAtPosition(GdkDevice *device, gint *x, gint *y) {
    return gdk_device_get_window_at_position(device, x, y);
}

void
GtkConfigureTransparencyAndRealize(GtkWidget *window,
                                             gboolean transparent) {
    gboolean isTransparent = ConfigureWindowTransparency(window, transparent);
    gtk_widget_realize(window);
    if (isTransparent) {
        GdkRGBA rgba = { 1.0, 1.0, 1.0, 0.0 };
        gdk_window_set_background_rgba(gtk_widget_get_window(window), &rgba);
    }
}

void
GtkWindowConfigureFromVisual(GtkWidget *widget, GdkVisual *visual) {
    gtk_widget_set_visual(widget, visual);
}

static gboolean
ConfigureTransparentWindow(GtkWidget *window) {
    GdkScreen *default_screen = gdk_screen_get_default();
    GdkDisplay *default_display = gdk_display_get_default();
    GdkVisual *visual = gdk_screen_get_rgba_visual(default_screen);
    if (visual
            && gdk_display_supports_composite(default_display)
            && gdk_screen_is_composited(default_screen)) {
        gtk_widget_set_visual(window, visual);
        return TRUE;
    }

    return FALSE;
}

static void
GrabMouseDevice(GdkDevice *device, DeviceGrabContext *context) {
    GdkInputSource source = gdk_device_get_source(device);
    if (source == GDK_SOURCE_MOUSE) {
        GdkGrabStatus status = gdk_device_grab(device,
                                               context->window,
                                               GDK_OWNERSHIP_NONE,
                                               TRUE,
                                               GDK_ALL_EVENTS_MASK,
                                               NULL,
                                               GDK_CURRENT_TIME);
        if (status == GDK_GRAB_SUCCESS) {
            context->grabbed = TRUE;
        }
    }
}

static void
UngrabMouseDevice(GdkDevice *device) {
    GdkInputSource source = gdk_device_get_source(device);
    if (source == GDK_SOURCE_MOUSE) {
        gdk_device_ungrab(device, GDK_CURRENT_TIME);
    }
}

int
GtkFixupTypedKey(int key, int keyval) {
    return key;
}

void
GdkWindowGetSize(GdkWindow *window, gint *w, gint *h) {
    *w = gdk_window_get_width(window);
    *h = gdk_window_get_height(window);
}

void
GdkDisplayGetPointer(GdkDisplay* display, gint* x, gint *y) {
    gdk_device_get_position(gdk_device_manager_get_client_pointer(gdk_display_get_device_manager(display)),
        NULL , x, y);
}

// -----------------------------------------------------------------------------
#else /* GTK_CHECK_VERSION(3, 0, 0) */

GdkScreen *
GdkWindowGetScreen(GdkWindow * gdkWindow) {
    return gdk_drawable_get_screen(GDK_DRAWABLE(gdkWindow));
}

GdkDisplay *
GdkWindowGetDisplay(GdkWindow * gdkWindow) {
    return gdk_drawable_get_display(GDK_DRAWABLE(gdkWindow));
}

gboolean
GdkMouseDevicesGrab(GdkWindow *gdkWindow) {
    return GdkMouseDevicesGrabWithCursor(gdkWindow, NULL, TRUE);
}

gboolean
GdkMouseDevicesGrab(GdkWindow *gdkWindow) {
    return GdkMouseDevicesGrabWithCursor(gdkWindow, NULL, TRUE);
}

gboolean
GdkMouseDevicesGrabWithCursor(GdkWindow *gdkWindow, GdkCursor *cursor) {
    return GdkMouseDevicesGrabWithCursor(gdkWindow, cursor, TRUE);
}

gboolean
GdkMouseDevicesGrabWithCursor(GdkWindow *gdkWindow, GdkCursor *cursor, gboolean owner_events) {
    if (g_disableGrab) {
        return TRUE;
    }
    GdkGrabStatus status = gdk_pointer_grab(gdkWindow, owner_events, (GdkEventMask)
                                            (GDK_POINTER_MOTION_MASK
                                                | GDK_POINTER_MOTION_HINT_MASK
                                                | GDK_BUTTON_MOTION_MASK
                                                | GDK_BUTTON1_MOTION_MASK
                                                | GDK_BUTTON2_MOTION_MASK
                                                | GDK_BUTTON3_MOTION_MASK
                                                | GDK_BUTTON_PRESS_MASK
                                                | GDK_BUTTON_RELEASE_MASK),
                                            NULL, cursor, GDK_CURRENT_TIME);

    return (status == GDK_GRAB_SUCCESS) ? TRUE : FALSE;
}

void
GdkMouseDevicesUngrab() {
    gdk_pointer_ungrab(GDK_CURRENT_TIME);
}

void
GdkMasterPointerGrab(GdkWindow *window, GdkCursor *cursor) {
    if (disableGrab) {
        gdk_window_set_cursor(window, cursor);
        return;
    }
    gdk_pointer_grab(window, FALSE, (GdkEventMask)
                     (GDK_POINTER_MOTION_MASK
                         | GDK_BUTTON_MOTION_MASK
                         | GDK_BUTTON1_MOTION_MASK
                         | GDK_BUTTON2_MOTION_MASK
                         | GDK_BUTTON3_MOTION_MASK
                         | GDK_BUTTON_RELEASE_MASK),
                     NULL, cursor, GDK_CURRENT_TIME);
}

void
GdkMasterPointerUngrab() {
    gdk_pointer_ungrab(GDK_CURRENT_TIME);
}

void
GdkMasterPointerGetPosition(gint *x, gint *y) {
    gdk_display_get_pointer(gdk_display_get_default(), NULL, x, y, NULL);
}

gboolean
GdkDeviceIsGrabbed(GdkDevice *device) {
    (void) device;
    return gdk_display_pointer_is_grabbed(gdk_display_get_default());
}

void
GdkDeviceUngrab(GdkDevice *device) {
    (void) device;
    gdk_pointer_ungrab(GDK_CURRENT_TIME);
}

GdkWindow *
GdkDeviceGetWindowAtPosition(GdkDevice *device, gint *x, gint *y) {
    (void) device;
    return gdk_display_get_window_at_pointer(gdk_display_get_default(), x, y);
}

void
GtkConfigureTransparencyAndRealize(GtkWidget *window, gboolean transparent) {
    configure_window_transparency(window, transparent);
    gtk_widget_realize(window);
}

void
GtkWindowConfigureFromVisual(GtkWidget *widget, GdkVisual *visual) {
    GdkColormap *colormap = gdk_colormap_new(visual, TRUE);
    gtk_widget_set_colormap(widget, colormap);
}

static gboolean
ConfigureTransparentWindow(GtkWidget *window) {
    GdkScreen *default_screen = gdk_screen_get_default();
    GdkDisplay *default_display = gdk_display_get_default();
    GdkColormap *colormap = gdk_screen_get_rgba_colormap(default_screen);
    if (colormap
            && gdk_display_supports_composite(default_display)
            && gdk_screen_is_composited(default_screen)) {
        gtk_widget_set_colormap(window, colormap);
        return TRUE;
    }

    return FALSE;
}

int
GtkFixupTypedKey(int key, int keyval) {
    if (key == 0) {
        // Work around "bug" fixed in gtk-3.0:
        // http://mail.gnome.org/archives/commits-list/2011-March/msg06832.html
        switch (keyval) {
        case 0xFF08 /* Backspace */: return '\b';
        case 0xFF09 /* Tab       */: return '\t';
        case 0xFF0A /* Linefeed  */: return '\n';
        case 0xFF0B /* Vert. Tab */: return '\v';
        case 0xFF0D /* Return    */: return '\r';
        case 0xFF1B /* Escape    */: return '\033';
        case 0xFFFF /* Delete    */: return '\177';
        }
    }
    return key;
}

void
GdkWindowGetSize(GdkWindow *window, gint *w, gint *h) {
    gdk_drawable_get_size(GDK_DRAWABLE(window), w, h);
}

void
GdkDisplayGetPointer(GdkDisplay* display, gint* x, gint *y) {
    gdk_display_get_pointer(display, NULL, x, y, NULL);
}

// -----------------------------------------------------------------------------
#endif /* GTK_CHECK_VERSION(3, 0, 0) */

const guchar*
GtkSelectionDataGetDataWithLength(
        GtkSelectionData * selectionData,
        gint * length) {
    if (selectionData == NULL) {
        return NULL;
    }

    *length = gtk_selection_data_get_length(selectionData);
    return gtk_selection_data_get_data(selectionData);
}

static void
ConfigureOpaqueWindow(GtkWidget *window) {
    gtk_widget_set_visual(window,
                          gdk_screen_get_system_visual(
                              gdk_screen_get_default()));
}

static gboolean
ConfigureWindowTransparency(GtkWidget *window, gboolean transparent) {
    if (transparent) {
        if (ConfigureTransparentWindow(window)) {
            return TRUE;
        }

        ERROR0("Can't create transparent stage, because your screen doesn't"
               " support alpha channel."
               " You need to enable XComposite extension.\n");
    }

    ConfigureOpaqueWindow(window);
    return FALSE;
}
