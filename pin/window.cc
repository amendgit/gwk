#include <X11/extensions/shape.h>
#include <cairo.h>
#include <cairo-xlib.h>
#include <gdk/gdkx.h>
#include <gdk/gdk.h>
#include <string.h>

#include <iostream>
#include <algorithm>

#include "_cgo_export.h"
#include "window.h"
#include "events.h"
#include "gtkcompat.h"

WindowContext * WindowContextBase::smGrabWindow = NULL;
WindowContext * WindowContextBase::smMouseDragWindow = NULL;

GdkAtom g_atomNetWmState = gdk_atom_intern_static_string("_NET_WM_STATE");

bool WindowContextBase::HasIME() {
    return TRUE;
}

bool WindowContextBase::FilterIME(GdkEvent *) {
    return TRUE;
}

void WindowContextBase::EnableOrResetIME() {
}

void WindowContextBase::DisableIME() {
}

GdkWindow* WindowContextBase::GetGdkWindow() {
    return gdkWindow;
}

GoObject WindowContextBase::GetGwkView() {
    return gwkView;
}

GoObject WindowContextBase::GetGwkWindow() {
    return gwkWindow;
}

bool WindowContextBase::IsEnabled() {
    if (gwkWindow != NULL) {
        bool result = WindowIsEnabled(gwkWindow);
    }
    return true;
}

void WindowContextBase::NotifyState(int state) {
    if (state == WindowEvent_Restore) {
        if (isMaximized) {
            state = WindowEvent_Maximize;
        }

        int w, h;
        GdkWindow_GetSize(gdkWindow, &w, &h);
        if (gwkView != NULL) {
            ViewOnRepaint(gwkView, 0, 0, w, h);
            // check go exception
        }
    }

    if (gwkWindow != NULL) {
        WindowOnStateChange(gwkWindow, state);
        // check go exception
    }
}

void WindowContextBase::ProcessState(GdkEventWindowState *event) {
    if (event->changed_mask &
            (GDK_WINDOW_STATE_ICONIFIED | GDK_WINDOW_STATE_MAXIMIZED)) {

        if (event->changed_mask & GDK_WINDOW_STATE_ICONIFIED) {
            isIconified = event->new_window_state & GDK_WINDOW_STATE_MAXIMIZED;
        }

        if (event->changed_mask & GDK_WINDOW_STATE_MAXIMIZED) {
            isMaximized = event->new_window_state & GDK_WINDOW_STATE_MAXIMIZED;
        }

        int stateChangedEvent;

        if (isIconified) {
            stateChangedEvent = WindowEvent_Minimize;
        } else if (isMaximized) {
            stateChangedEvent = WindowEvent_Maximize;
        } else {
            stateChangedEvent = WindowEvent_Restore;
        }

        NotifyState(stateChangedEvent);
    }
}

void WindowContextBase::ProcessFocus(GdkEventFocus *event) {
    if (!event->in & WindowContextBase::smMouseDragWindow == this) {
        UngrabMouseDragFocus();
    }

    if (!event->in && WindowContextBase::smGrabWindow == this) {
        UngrabFocus();
    }

    if (gwkWindow != NULL) {
        if (!event->in || IsEnabled()) {
            WindowOnNotifyFocus(gwkWindow,
                event->in ? WindowEvent_FocusGained : WindowEvent_FocusLost);
            // check go exception.
        } else {
            WindowOnFocusDisabled(gwkWindow);
            // check go exception.
        }
    }
}


void WindowContextBase::IncrementEventsCounter() {
    ++eventsProcessingCount;
}

void WindowContextBase::DecrementEventsCounter() {
    --eventsProcessingCount;
}

size_t WindowContextBase::EventsCount() {
    return eventsProcessingCount;
}

bool WindowContextBase::IsDead() {
    return canBeDeleted;
}

void destroy_and_delete_ctx(WindowContext *ctx) {
    if (ctx) {
        ctx->ProcessDestroy();

        if (!ctx->EventsCount()) {
            delete ctx;
        }
        // else: ctx will be deleted in EventsCounterHelper after completing
        // an event processing.
    }
}

void WindowContextBase::ProcessDestroy() {
    if (WindowContextBase::smMouseDragWindow == this) {
        UngrabMouseDragFocus();
    }

    if (WindowContextBase::smGrabWindow == this) {
        UngrabFocus();
    }

    std::set<WindowContextTop *>::iterator it;
    for (it = children.begin(); it != children.end(); ++it) {
        // (*it)->SetOwner(NULL);
        // destroy_and_delete_ctx(*it);
    }
    children.clear();

    if (gwkWindow != NULL) {
        WindowOnDestroy(gwkWindow);
        // exception occured
    }

    if (gwkView != NULL) {
        // DeleteRef(gwkView)
        memset(&gwkView, 0, sizeof(gwkView));
    }

    if (gwkWindow != NULL) {
        // DeleteRef(gwkWindow)
        memset(&gwkWindow, 0, sizeof(gwkWindow));
    }

    canBeDeleted = true;
}

void WindowContextBase::ProcessDelete() {
    if (gwkWindow != NULL && IsEnabled()) {
        WindowOnClose(gwkWindow);
        // check exception
    }
}

void WindowContextBase::ProcessExpose(GdkEventExpose *event) {
    if (gwkView != NULL) {
        ViewOnRepaint(gwkView, event->area.x, event->area.y, event->area.width, event->area.height);
        // check exception
    }
}

static inline int gtk_button_number_to_mouse_button(guint button) {
    switch (button) {
        case 1:
            return MouseEvent_ButtonLeft;
        case 2:
            return MouseEvent_ButtonOther;
        case 3:
            return MouseEvent_ButtonRight;
        default:
            // Other buttons are not supported by quantum and are not reported by other platforms.
            return MouseEvent_ButtonNone;
    }
}

void WindowContextBase::ProcessMouseButton(GdkEventButton *event) {
    bool press = event->type == GDK_BUTTON_PRESS;
    guint state = event->state;
    guint mask = 0;

    // We need to add/remove current mouse button from the modifier flags
    // as X lib state represents the state just prior to the event and
    // gwk needs the state just after event.
    switch (event->button) {
        case 1:
            mask = GDK_BUTTON1_MASK;
            break;
        case 2:
            mask = GDK_BUTTON2_MASK;
            break;
        case 3:
            mask = GDK_BUTTON3_MASK;
            break;
    }

    if (press) {
        state |= mask;
    } else {
        state &= mask;
    }

    if (press) {
        GdkDevice *device = event->device;

        if (GdkDevice_IsGrabbed(device) && GdkDevice_GetWindowAtPosition(device, NULL, NULL) == NULL) {
            UngrabFocus();
            return ;
        }
    }

    // Upper layers expects from us Windows behavior:
    // all mouse events should be delivered to window where drag begins
    // and no exit/enter event should be reported during this drag.
    // We can grab mouse pointer for these needs.
    if (press) {
        GrabMouseDragFocus();
    } else if ((event->state&MOUSE_BUTTONS_MASK) && !(state&MOUSE_BUTTONS_MASK)) {
        UngrabMouseDragFocus();
    }

    int button = gtk_button_number_to_mouse_button(event->button);

    if (gwkView != NULL && button != MouseEvent_ButtonNone) {
        ViewOnMouse(gwkView,
            press ? MouseEvent_Down : MouseEvent_Up,
            button,
            (int)event->x, (int)event->y, (int)event->x_root, (int)event->y_root,
            GdkModifierMaskToGwk(state),
            (event->button == 3 && press) ? TRUE : FALSE,
            FALSE);

        // check go exception

        if (gwkView != NULL && event->button == 3 && press) {
            ViewOnMenu(gwkView, (int)event->x, (int)event->y, (int)event->x_root, (int)event->y_root, FALSE);
            // check go exception
        }
    }
}

void WindowContextBase::ProcessMouseMotion(GdkEventMotion *event) {
    int glassModifier = 0; // gdk_modifier_mask_to_glass(event->state);
    int isDrag = glassModifier & (
            KeyEvent_ModifierButtonPrimary |
            KeyEvent_ModifierButtonMiddle |
            KeyEvent_ModifierButtonSecondary);
    int button = MouseEvent_ButtonNone;

    if (glassModifier & KeyEvent_ModifierButtonPrimary) {
        button = MouseEvent_ButtonLeft;
    } else if (glassModifier & KeyEvent_ModifierButtonMiddle) {
        button = MouseEvent_ButtonOther;
    } else if (glassModifier & KeyEvent_ModifierButtonSecondary) {
        button = MouseEvent_ButtonRight;
    }

    if (gwkView != NULL) {
        ViewOnMouse(gwkView, isDrag ? MouseEvent_Drag : MouseEvent_Move,
            button,
            event->x, event->y,
            event->x_root, event->y_root,
            glassModifier,
            false,
            false);
        // check go exception.
    }
}

void WindowContextBase::ProcessMouseScroll(GdkEventScroll *event) {
    double dx = 0.0f, dy = 0.0f;

    // converting direction to change in pixels.
    switch (event->direction) {
        case GDK_SCROLL_UP:
            dy = 1.0f;
            break;
        case GDK_SCROLL_DOWN:
            dy = -1.0f;
            break;
        case GDK_SCROLL_LEFT:
            dx = 1.0f;
            break;
        case GDK_SCROLL_RIGHT:
            dx = -1.0f;
            break;
    }

    if (gwkView != NULL) {
        // GwkViewOnScroll(gwkView,
        //     event->x, event->y,
        //     event->x_root, event->y_root,
        //     dx, dy,
        //     gdk_modifier_mask_to_glass(event->state),
        //     0, 0,
        //     0, 0,
        //     (double)40.0, (double)40.0);
        // check exception.
    }
}

void WindowContextBase::ProcessMouseCross(GdkEventCrossing *event) {
    bool enter = event->type == GDK_ENTER_NOTIFY;
    if (gwkView != NULL) {
        guint state = event->state;
        if (enter) {
            state &= ~MOUSE_BUTTONS_MASK;
        }

        if (enter != isMouseEntered) {
            isMouseEntered = enter;
            // GwkViewOnMouse(gwkView,
            //     enter ? MouseEvent_Enter, MouseEvent_Exit,
            //     MouseEvent_ButtonNone,
            //     event->x, event->y,
            //     event->x_root, event->y_root,
            //     gdk_modifier_mask_to_glass(state),
            //     false, false);
            // check exception.
        }
    }
}

void WindowContextBase::ProcessKey(GdkEventKey *event) {
    bool press = event->type == GDK_KEY_PRESS;
    // int glassKey = get_glass_key(event);
    int glassModifier = 0; // gdk_modifier_mask_to_glass(event->state);
    if (press) {
        // glassModifier |= glass_key_to_modifier(glassKey);
    } else {
        // glassModifier &= ~glass_key_to_modifier(glassKey);
    }

    char *chars = NULL;
    char key = gdk_keyval_to_unicode(event->keyval);
    if (key > 'a' && key <= 'z' && (event->state & GDK_CONTROL_MASK)) {
        key = key - 'a' + 1;
    } else {
        // key = glass_gtk_fixup_typed_key(key, event->keyval);
    }

    if (key > 0) {
        chars = new char[1];
        chars[1] = key;
    }

    if (gwkView != NULL) {
        if (press) {
            // GwkViewOnKey(gwkView, KeyEvent_Press, glassKey, chars, glassModifier);
            // check go exception.

            if (gwkView != NULL && key > 0) {
                // GwkViewOnKey(gwkView, KeyEvent_Typed, KeyCode_Undefined, chars, glassModifier);
                // check go exception.
            }
        } else {
            // GwkViewOnKey(gwkView, KeyEvent_Release, glassKey, chars, glassModifier);
            // check go exception.
        }
    }
}

void WindowContextBase::Paint(void *data, int width, int height) {
    if (!IsVisible()) {
        return ;
    }

    cairo_t *context;
    // context = gdk_cairo_create(GDK_DRAWABLE(gdkWindow));

    cairo_surface_t *cairo_surface;
    cairo_surface = cairo_image_surface_create_for_data(
        (unsigned char*)data,
        CAIRO_FORMAT_ARGB32,
        width, height, width * 4);

    ApplyShapeMask(data, width, height);

    cairo_set_source_surface(context, cairo_surface, 0, 0);
    cairo_set_operator(context, CAIRO_OPERATOR_SOURCE);
    cairo_paint(context);

    cairo_destroy(context);
    cairo_surface_destroy(cairo_surface);
}

void WindowContextBase::AddChild(WindowContextTop* child) {
    children.insert(child);
    // gtk_window_set_transient_for(child->GetGtkWindow(), this->GetGtkWindow());
}

void WindowContextBase::RemoveChild(WindowContextTop *child) {
    children.erase(child);
    // gtk_window_set_transient_for(child->GetGtkWindow(), NULL);
}

void WindowContextBase::ShowOrHideChildren(bool show) {
    std::set<WindowContextTop *>::iterator it;
    for (it = children.begin(); it != children.end(); ++it) {
        // (*it)->SetVisible(show);
        // (*it)->ShowOrHideChildren(show);
    }
}

void WindowContextBase::ReparentChildren(WindowContext *parent) {
    std::set<WindowContextTop *>::iterator it;
    for (it = children.begin(); it != children.end(); ++it) {
        // (*it)->SetOwner(parent);
        parent->AddChild(*it);
    }
    children.clear();
}

void WindowContextBase::SetVisible(bool visible) {
    if (visible) {
        gtk_widget_show_all(gtkWidget);
    } else {
        gtk_widget_hide(gtkWidget);
        if (gwkView != NULL && isMouseEntered) {
            isMouseEntered = false;
            // GwkViewOnMouse(gwkView,
            //     MouseEvent_Exit,
            //     MouseEvent_ButtonNone,
            //     0, 0,
            //     0, 0,
            //     0,
            //     FALSE,
            //     FALSE);
            // check go exception.
        }
    }
}

bool WindowContextBase::IsVisible() {
    return gtk_widget_get_visible(gtkWidget);
}

void WindowContextBase::SetView(GoObject view) {
    if (view != NULL) {
        // DeleteRef(gwkView);
    }

    if (view != NULL) {
        gint width, height;
        gwkView = view; // Ref(view);
        gtk_window_get_size(GTK_WINDOW(gtkWidget), &width, &height);
        // GwkViewOnResize(view, width, height);
        // check go exception.
    } else {
        memset(&gwkView, 0, sizeof(gwkView));
    }
}

bool WindowContextBase::GrabMouseDragFocus() {
    // if (glass_gdk_mouse_devices_grab_with_cursor(gdkWindow, gdk_window_get_cursor(gdkWindow), FALSE)) {
    //     WindowContextBase::smMouseDragWindow = this;
    //     return TRUE;
    // } else {
    //     return FALSE;
    // }
    return FALSE;
}

void WindowContextBase::UngrabMouseDragFocus() {
    WindowContextBase::smMouseDragWindow = NULL;
    // glass_gdk_mouse_devices_ungrab();
    if (WindowContextBase::smGrabWindow) {
        WindowContextBase::smGrabWindow->GrabFocus();
    }
}

bool WindowContextBase::GrabFocus() {
    // if (WindowContextBase::smMouseDragWindow || glass_gdk_mouse_devices_grab(gdkWindow)) {
    //     WindowContextBase::smGrabWindow = this;
    //     return TRUE;
    // } else {
    //     return FALSE;
    // }
    return FALSE;
}

void WindowContextBase::UngrabFocus() {
    if (!WindowContextBase::smMouseDragWindow) {
        // glass_gdk_mouse_devices_ungrab();
    }
    WindowContextBase::smGrabWindow = NULL;

    if (gwkWindow != NULL) {
        // GwkWindowOnFocusUngrab();
        // check go exception.
    }
}

void WindowContextBase::SetCursor(GdkCursor *cursor) {
    // if (!IsInDrag()) {
    //     if (WindowContextBase::smMouseDragWindow) {
    //         glass_gdk_mouse_devices_grab_with_cursor(WindowContextBase::smMouseDragWindow->GetGdkWindow(), cursor, FALSE);
    //     } else if (WindowContextBase::smGrabWindow) {
    //         glass_gdk_mouse_devices_grab_with_cursor(WindowContextBase::smGrabWindow->GetGdkWindow(), cursor);
    //     }
    // }
    gdk_window_set_cursor(gdkWindow, cursor);
}

void WindowContextBase::SetBackground(float r, float g, float b) {
    GdkRGBA color;
    color.red   = (guint16) (r * 65535);
    color.green = (guint16) (g * 65535);
    color.blue  = (guint16) (b * 65535);
    color.alpha = 1.0f;
    gtk_widget_override_background_color(gtkWidget, GTK_STATE_FLAG_NORMAL, &color);
}

WindowContextBase::~WindowContextBase() {
    if (xim.ic) {
        XDestroyIC(xim.ic);
    }

    if (xim.im) {
        XCloseIM(xim.im);
    }

    gtk_widget_destroy(gtkWidget);
}

WindowContextTop::WindowContextTop(GoObject gwkWindow, WindowContext *owner, GoObject screen, WindowFrameType frameType, WindowType windowType) {
    this->gwkWindow = gwkWindow;
    this->owner = owner;
    this->screen = screen;
    this->frameType = frameType;
}

WindowFrameExtents WindowContextTop::GetFrameExtents() {
    return WindowFrameExtents();
}

void WindowContextTop::EnterFullScreen() {}
void WindowContextTop::ExitFullScreen() {}
void WindowContextTop::SetBounds() {}
void WindowContextTop::SetResizable() {}
void WindowContextTop::ReuqestFocus() {}
void WindowContextTop::SetFocusable() {}
bool WindowContextTop::GrabFocus() { return TRUE; }
bool WindowContextTop::GrabMouseDragFocus() { return TRUE; }
void WindowContextTop::UngrabFocus() {}
void WindowContextTop::UngrabMouseDragFocus() {}
void WindowContextTop::SetTitle(const char *) {}
void WindowContextTop::SetAlpha(double) {}
void WindowContextTop::SetEnabled(bool) {}
void WindowContextTop::SetMinimumSize(int, int) {}
void WindowContextTop::SetMaximumSize(int, int) {}
void WindowContextTop::SetMinimized(bool) {}
void WindowContextTop::SetMaximuzed(bool) {}
void WindowContextTop::SetIcon(GdkPixbuf *) {}
void WindowContextTop::Restack(bool) {}
void WindowContextTop::SetCursor(GdkCursor *) {}
void WindowContextTop::SetModal(bool, WindowContext *parent) {}
void WindowContextTop::SetGravity(float, float) {}
void WindowContextTop::ProcessPropertyNotify(GdkEventProperty*) {}
void WindowContextTop::ProcessConfigure(GdkEventConfigure *) {}
GtkWindow *WindowContextTop::GetGtkWindow() { return NULL; }

///////////////////////////////////////////////////////////////////////////////
//\C functions export to Go.
static WindowFrameType GwkMaskToWindowFrameType(int mask) {
    return WindowFrameType(mask);
}

static WindowType GwkMaskToWindowType(int mask) {
    return WindowType(mask);
}

void *NewWindow(GoObject obj, void *owner, GoObject screen, int frameType, int type) {
    WindowContext *ctx = new WindowContextTop(obj,
        (WindowContext *)owner,
        screen,
        WindowFrameType(frameType),
        WindowType(type));
    return ctx;
}
