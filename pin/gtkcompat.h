#ifndef GTKCOMPAT_H
#define GTKCOMPAT_H

#include <gtk/gtk.h>
#include <gdk/gdkx.h>
#include <gdk/gdkkeysyms.h>

// GdkDragContext versions compatible.
#if GTK_CHECK_VERSION(2, 22, 0)

#define GdkDragContext_GetSselectedAction(context) gdk_drag_context_get_selected_action(context)
#define GdkDragContext_GetActions(context) gdk_drag_context_get_actions(context)
#define GdkDragContext_ListTargets(context) gdk_drag_context_list_targets(context)
#define GdkDragContext_GetSuggestedAction(context) gdk_drag_context_get_suggested_action(context)

#else /* GTK_CHECK_VERSION(2, 22, 0) */

#define GdkDragContext_GetSselectedAction(context) (context->action)
#define GdkDragContext_GetActions(context) (context->actions)
#define GdkDragContext_ListTargets(context) (context->targets)
#define GdkDragContext_GetSuggestedAction(context) (context->suggested_action)

#endif /* GTK_CHECK_VERSION(2, 22, 0) */


// GdkWindow versions compatible.
#if GTK_CHECK_VERSION(2, 24, 0)

#define GDK_KEY_CONSTANT(key) (GDK_KEY_ ## key)
#define GdkWindow_ForeignNewForDisplay(display, anid) gdk_x11_window_foreign_new_for_display(display, anid)
#define GdkWindow_LookupForDisplay(display, anid) gdk_x11_window_lookup_for_display(display, anid)

#else /* GTK_CHECK_VERSION(2, 24, 0) */

#define GDK_KEY_CONSTANT(key) (GDK_ ## key)
#define GdkWindow_ForeignNewForDisplay(display, anid) gdk_window_foreign_new_for_display(display, anid)
#define GdkWindow_LookupForDisplay(display, anid) gdk_window_lookup_for_display(display, anid)

#endif /* GTK_CHECK_VERSION(2, 24, 0) */


#if GTK_CHECK_VERSION(3, 0, 0)

#define GtkWindow_SetHasResizeGrip(window, value) gtk_window_set_has_resize_grip(window, TRUE)
#define GdkSelectionEvent_GetRequestor(event) (event->requestor)
#define GdkDragContext_GetDestWindow(context) gdk_drag_context_get_dest_window(context)

#else /* GTK_CHECK_VERSION(3, 0, 0) */

#define GtkWindow_SetHasResizeGrip(window, value) \
    (void) window; \
    (void) value;
#define GdkSelectionEvent_GetRequestor(event) GdkWindow_ForeignNewForDisplay(gdk_display_get_default(), event->requestor)
#define GdkDragContext_GetDestWindow(context) ((context != NULL) ? context->dest_window : NULL)

#endif /* GTK_CHECK_VERSION(3, 0, 0) */


GdkScreen *GdkWindow_GetScreen(GdkWindow * gdkWindow);
GdkDisplay *GdkWindow_GetDisplay(GdkWindow * gdkWindow);

gboolean GdkMouseDevices_Grab(GdkWindow * gdkWindow);
void GdkMouseDevices_Ungrab();
gboolean GdkMouseDevices_GrabWithCursor(GdkWindow * gdkWindow, GdkCursor *cursor);
gboolean GdkMouseDevices_GrabWithCursor(GdkWindow * gdkWindow, GdkCursor *cursor, gboolean owner_events);

void GdkMasterPointer_Grab(GdkWindow *window, GdkCursor *cursor);
void GdkMasterPointer_Ungrab();
void GdkMasterPointer_GetPosition(gint *x, gint *y);

gboolean GdkDevice_IsGrabbed(GdkDevice *device);
void GdkDevice_Ungrab(GdkDevice *device);
GdkWindow *GdkDevice_GetWindowAtPosition(GdkDevice *device, gint *x, gint *y);

void GtkWidget_ConfigureTransparencyAndRealize(GtkWidget *widget, gboolean transparent);

const guchar * GtkSelectionData_GetDataWithLength(GtkSelectionData * selectionData,
    gint * length);

void GtkWidget_ConfigureFromVisual(GtkWidget *widget, GdkVisual *visual);
int GtkFixupTypedKey(int key, int keyval);

void GdkWindow_GetSize(GdkWindow *window, gint *w, gint *h);

void GdkDisplay_GetPointer(GdkDisplay* display, gint* x, gint *y);

#endif /* GTKCOMPAT_H */
