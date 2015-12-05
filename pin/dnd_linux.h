#ifndef DND_LINUX_H
#define DND_LINUX_H

#include "general_linux.h"
#include "window_linux.h"

#include <gtk/gtk.h>

void
ProcessDNDTarget(WindowContext *, GdkEventDND *);
int
ProcessDNDSource(WindowContext *, GdkEvent *);

#endif /* DND_LINUX_H */
