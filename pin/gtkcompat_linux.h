#ifndef GTKCOMPAT_LINUX_H
#define GTKCOMPAT_LINUX_H

#include <gtk/gtk.h>
#include <gdk/gdkx.h>
#include <gdk/gdkkeysyms.h>

// GdkDragContext versions compatible.
#if GTK_CHECK_VERSION(2, 22, 0)

#define GdkDragContextGetSelectedAction(context) gdk_drag_context_get_selected_action(context)
#define GdkDragContextGetActions(context) gdk_drag_context_get_actions(context)
#define GdkDragContextListTargets(context) gdk_drag_context_list_targets(context)
#define GdkDragContextGetSuggestedAction(context) gdk_drag_context_get_suggested_action(context)

#else /* GTK_CHECK_VERSION(2, 22, 0) */

#define GdkDragContextGetSselectedAction(context) (context->action)
#define GdkDragContextGetActions(context) (context->actions)
#define GdkDragContextListTargets(context) (context->targets)
#define GdkDragContextGetSuggestedAction(context) (context->suggested_action)

#endif /* GTK_CHECK_VERSION(2, 22, 0) */


// GdkWindow versions compatible.
#if GTK_CHECK_VERSION(2, 24, 0)

#define GDK_KEY_CONSTANT(key) (GDK_KEY_ ## key)
#define GdkWindowForeignNewForDisplay(display, anid) gdk_x11_window_foreign_new_for_display(display, anid)
#define GdkWindowLookupForDisplay(display, anid) gdk_x11_window_lookup_for_display(display, anid)

#else /* GTK_CHECK_VERSION(2, 24, 0) */

#define GDK_KEY_CONSTANT(key) (GDK_ ## key)
#define GdkWindowForeignNewForDisplay(display, anid) gdk_window_foreign_new_for_display(display, anid)
#define GdkWindowLookupForDisplay(display, anid) gdk_window_lookup_for_display(display, anid)

#endif /* GTK_CHECK_VERSION(2, 24, 0) */


#if GTK_CHECK_VERSION(3, 0, 0)

#define GtkWindowSetHasResizeGrip(window, value) gtk_window_set_has_resize_grip(window, TRUE)
#define GdkSelectionEventGetRequestor(event) (event->requestor)
#define GdkDragContextGetDestWindow(context) gdk_drag_context_get_dest_window(context)

#else /* GTK_CHECK_VERSION(3, 0, 0) */

#define GtkWindowSetHasResizeGrip(window, value) \
    (void) window; \
    (void) value;
#define GdkSelectionEventGetRequestor(event) GdkWindow_ForeignNewForDisplay(gdk_display_get_default(), event->requestor)
#define GdkDragContextGetDestWindow(context) ((context != NULL) ? context->dest_window : NULL)

#endif /* GTK_CHECK_VERSION(3, 0, 0) */


GdkScreen *
GdkWindowGetScreen(GdkWindow * gdkWindow);
GdkDisplay *
GdkWindowGetDisplay(GdkWindow * gdkWindow);

gboolean
GdkMouseDevicesGrab(GdkWindow * gdkWindow);
void
GdkMouseDevicesUngrab();
gboolean
GdkMouseDevicesGrabWithCursor(GdkWindow * gdkWindow, GdkCursor *cursor);
gboolean
GdkMouseDevicesGrabWithCursor(GdkWindow * gdkWindow, GdkCursor *cursor, gboolean owner_events);

void
GdkMasterPointerGrab(GdkWindow *window, GdkCursor *cursor);
void
GdkMasterPointerUngrab();
void
GdkMasterPointerGetPosition(gint *x, gint *y);

gboolean
GdkDeviceIsGrabbed(GdkDevice *device);
void
GdkDeviceUngrab(GdkDevice *device);
GdkWindow *
GdkDeviceGetWindowAtPosition(GdkDevice *device, gint *x, gint *y);

void
GtkWidgetConfigureTransparencyAndRealize(GtkWidget *widget, gboolean transparent);

const guchar *
GtkSelectionDataGetDataWithLength(GtkSelectionData * selectionData,
    gint * length);

void
GtkWidgetConfigureFromVisual(GtkWidget *widget, GdkVisual *visual);
int
GtkFixupTypedKey(int key, int keyval);

void
GdkWindowGetSize(GdkWindow *window, gint *w, gint *h);

void
GdkDisplayGetPointer(GdkDisplay* display, gint* x, gint *y);

#endif /* GTKCOMPAT_LINUX_H */
