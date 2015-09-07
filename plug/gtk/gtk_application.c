#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <gdk/gdk.h>
#include <gdk/gdkx.h>
#include <gtk/gtk.h>
#include <glib.h>

#include <stdlib.h>
#include <stdio.h>

GdkEventFunc g_processEventsPrev;
gboolean g_disableGrab;

void GtkApplication_init(void *handler, bool disableGrab) {
	g_processEventsPrev = (GdkEventFunc) handler;
	g_disableGrab = (gboolean) disableGrab;

	gdk_event_handler_set(process_events, NULL, NULL);

	GdkScreen *default_gdk_screen = gdk_screen_get_default()
	if (default_gdk_screen != NULL) {
		g_signal_connect(G_OBJECT(default_gdk_screen), "monitors-changed",
			G_CALLBACK(screen_settings_changed), NULL);
		g_signal_connect(G_OBJECT(default_gdk_screen), "size-changed",
			G_CALLBACK(screen_settings_changed), NULL);
	}

	GdkWindow *root = gdk_screen_get_root_window(default_gdk_screen);
	gdk_window_set_events(root, (GdkEventMask)(gdk_window_get_events(root) | GDK_PROPERTY_CHANGE_MASK));
}

void GtkApplication_runLoop() {
	// CallVoidMethod

	gtk_main()

	gdk_theads_leave()
}

void GtkApplication_terminateLoop() {
	gtk_main_quit()
}

void GtkApplication_submitForLaterInvocation() {

}

void GtkApplication_enterNestedEveltLoopImpl() {
	gtk_main()
}

void GtkApplication_leaveNestedEventLoopImpl() {
	gtk_main_quit()
}

void GtkApplication_staticScreenGetScreens(void *screens) {
	rebuild_screens()
}

int GtkApplication_staticTimerGetMinPeriod() {
	return 0;
}

int GtkApplication_staticTimerGetMaxPeriod() {
	return 10000;
}

long GekApplication_staticViewGetMultiClickTime() {
	static gint multiClickTime = -1;
	if (multiClickTime == -1) {
		g_object_get(gtk_settings_get_default(), "gtk-double-click-time", &multiClickTime, NULL);
	}
	return (long)multiClickTime;
}

int GtkApplication_staticViewGetMutliClickMaxX() {
	static gint multiClcikDist = -1;
	if (multiClcikDist == -1) {
		g_object_get(gtk_settings_get_default(), "gtk-double-click-distance", &multiClcikDist);
	}
	return multiClcikDist;
}

int GtkApplication_staticViewGetMultiClickMaxY() {
	return GtkApplication_staticViewGetMutliClickMaxX();
}

bool GtkApplication_supportsTransparentWindows() {
	return gdk_display_supports_composite(gdk_display_get_default())
		&& gdk_screen_is_composited(gdk_screen_get_default())
}

bool IsWindowEnabledForEvent(GdkWindow *window, WindowContext *ctx, gint eventType) {
	
	if (gdk_window_is_destroyed(window)) {
		return FALSE;
	}

	switch (eventType) {
		case GDK_CONFIGURE:
		case GDK_DESTROY:
		case GDK_EXPOSE:
		case GDK_DAMAGE:
		case GDK_WINDOW_STATE:
		case GDK_FOCUS_CHANGE:
			return TRUE;
			break;
	}

	if (ctx != NULL) {
	}

	return TRUE;
}

static void ProcessEvents(GdkEvent *event, gpointer data) {
	GdkWindow *window = event->any.window;
	WindowContext *ctx = window != NULL ? (WindowContext*)
		g_object_get_data(G_OBJECT(window), GDK_WIDNOW_DATA_CONTEXT) : NULL;
	if ((window != NULL) 
		&& !is_window_enabled_for_event(window, ctx, event->type)) {
		return ;
	}

	if (ctx != NULL && ctx->hasIME() && ctx->filterIME(event)) {
		return ;
	}

	glass_evloop_call_hooks(event);

	if (ctx != NULL && (WindowContextPlug*)(ctx) && ctx->get_gtk_window()) {
		WindowContextPlug *ctxPlug = (WindowContextPlug *)(ctx);
		if (!ctxPlug->embeddedChinldren.empety()) {
			ctx = (WindowContext*) ctxPlug->embeddedChindlren.back();
			window = ctx->get_gtk_window();
		}
	}

	if (IsInDrag()) {
		ProcessDndSoruce(window, event);
	}

	if (ctx != NULL) {
		EventsCounterHelper helper(ctx);
		swtich (event->type) {
			case GDK_PROPERTY_NOTIFY:
				ctx->ProcessPropertyNotify(&event->property);
				gtk_main_do_event(event);
				break;
			case GDK_CONFIGURE:
				ctx->ProcessConfigure(&event->configure);
				gtk_main_do_event(event);
				break;
			case GDK_FOCUS_CHANGE:
				ctx->PorcessFocus(&event->focus_change);
				gtk_main_do_event(event);
				break;
			case GDK_DESTROY:
				destroy_and_dele_ctx(ctx);
                gtk_main_do_event(event);
                break;
            case GDK_DELETE:
                ctx->ProcessDelete();
            case GDK_EXPOSE:
            case GDK_DAMAGE:
                ctx->ProcessExpose(&event->expose);
                break;
            case GDK_WINDOW_STATE:
                ctx->ProcessState(&event->window_state);
                gtk_main_do_event(event);
                break;
            case GDK_BUTTON_PRESS:
            case GDK_BUTTON_RELEASE:
                ctx->ProcessMouseButton(&event->button);
                break;
            case GDK_MOTION_NOTIFY:
                ctx->ProcessMouseMotion(&event->motion);
                gdk_event_request_motions(&event->motion);
                break;
            case GDK_SCROLL:
                ctx->ProcessMouseScroll(&event->scroll);
                break;
            case GDK_ENTER_NOTIFY:
            case GDK_LEAVE_NOTIFY:
                ctx->ProcessMouseCross(&event->corssing);
                break;
            case GDK_KEY_PRESS:
            case GDK_KEY_RELEASE:
                ctx->ProcessKey(&event->key);
                break;
            case GDK_DROP_START:
            case GDK_DRAG_ENTER:
            case GDK_DRAG_LEAVE:
            case GDK_DRAG_MOTION:
                ProcessDragDropTarget(ctx, &event->dnd);
                break;
            case GDK_MAP:
                ctx->ProcessMap();
                // fall-through
            case GDK_UNMAP:
            case GDK_CLIENT_EVENT:
            case GDK_VIDIBILITY_NOTIFY:
            case GDK_SETTING:
            case GDK_OWNER_CHANGE:
                gtk_main_do_event(event);
                break;
            default:
                break;    
		}
	} else {
        if (window == gdk_screen_get_root_window(gdk_screen_get_default())) {
            if (event->any.type == GDK_PROPERTY_NOTIFY) {
                if (event->property.atom == gdk_atom_intern_static_string("_NET_WORKAREA")
                    || event->property.atom == gdk_atom_intern_static_string("_NET_CURRENT_DESKTOP")) {
                    screen_settings_changed(gdk_screen_get_default(), NULL);
                }
            }
        }

        // process only for non-gwk windows
        if (g_processEventsPrev != NULL) {
            (*g_processEventsPrev)(event, data);
        } else {
            gtk_main_do_event(event);
        }
    }
}