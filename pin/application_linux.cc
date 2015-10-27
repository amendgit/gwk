#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <gdk/gdk.h>
#include <gdk/gdkx.h>
#include <gtk/gtk.h>
#include <glib.h>

#include "_cgo_export.h"

char* HelloStringFromGo();
void CallVoidFunc(GoInterface p0);

static void
activate (GtkApplication* app, gpointer user_data) {
  GtkWidget *window;

  window = gtk_application_window_new (app);
  gtk_window_set_title (GTK_WINDOW (window), "Window");
  gtk_window_set_default_size (GTK_WINDOW (window), 200, 200);
  gtk_widget_show_all (window);
  return ;
}

int
GtkShowWindow () {
  GtkApplication *app;
  int status;

  app = gtk_application_new ("org.gtk.example", G_APPLICATION_FLAGS_NONE);
  g_signal_connect (app, "activate", G_CALLBACK (activate), NULL);
  status = g_application_run (G_APPLICATION (app), 0, NULL);
  g_object_unref (app);

  return status;
}

void
ScreenSettingsChanged() {
	return ;
}

static void
ProcessEvents(GdkEvent *event, gpointer data) {
	return ;
}

GdkEventFunc g_processEventsPrev;
extern gboolean     g_disableGrab;

void
GtkApplicationInit(void *handler, bool disableGrab) {
	g_processEventsPrev = (GdkEventFunc) handler;
	g_disableGrab = (gboolean) disableGrab;

	gdk_event_handler_set(ProcessEvents, NULL, NULL);

	GdkScreen *default_gdk_screen = gdk_screen_get_default();
	if (default_gdk_screen != NULL) {
		g_signal_connect(G_OBJECT(default_gdk_screen), "monitors-changed", G_CALLBACK(ScreenSettingsChanged), NULL);
		g_signal_connect(G_OBJECT(default_gdk_screen), "size-changed", G_CALLBACK(ScreenSettingsChanged), NULL);
	}

	GdkWindow *root = gdk_screen_get_root_window(default_gdk_screen);
	gdk_window_set_events(root, (GdkEventMask)(gdk_window_get_events(root) | GDK_PROPERTY_CHANGE_MASK));
}

void
GtkApplicationRunLoop() {
	// CallVoidFunc

	gtk_main();

	// gdk_theads_leave();
}

void
GtkApplicationTerminateLoop() {
	gtk_main_quit();
}

void
GtkApplicationEnterNestedEventLoopImpl() {
	gtk_main();
}

void
GtkApplicationLeaveNestedEventLoopImpl() {
	gtk_main_quit();
}

int
GtkApplicationStaticTimerGetMinPeriod() {
	return 0;
}

int
GtkApplicationStaticTimerGetMaxPeriod() {
	return 10000;
}

long
GtkApplicationStaticViewGetMultiClickTime() {
	static gint multiClickTime = -1;
	if (multiClickTime == -1) {
		g_object_get(gtk_settings_get_default(), "gtk-double-click-time", &multiClickTime, NULL);
	}
	return (long)multiClickTime;
}

int
GtkApplicationStaticViewGetMutliClickMaxX() {
	static gint multiClcikDist = -1;
	if (multiClcikDist == -1) {
		g_object_get(gtk_settings_get_default(), "gtk-double-click-distance", &multiClcikDist, NULL);
	}
	return multiClcikDist;
}

int
GtkApplicationSaticViewGetMultiClickMaxY() {
	return GtkApplicationStaticViewGetMutliClickMaxX();
}

int
GtkApplicationSupportsTransparentWindows() {
	return gdk_display_supports_composite(gdk_display_get_default())
		&& gdk_screen_is_composited(gdk_screen_get_default());
}

int
IsWindowEnabledForEvent(GdkWindow *window, gint eventType) {

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

	return TRUE;
}
