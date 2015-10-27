#ifndef GTK_APPLICTION_H
#define GTK_APPLICTION_H

#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>

int GtkShowWindow();

void GtkApplicationInit(void *handler, bool disableGrab);

#ifdef __cplusplus
}
#endif

#endif
