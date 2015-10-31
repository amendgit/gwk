#ifndef KEY_LINUX_H
#define KEY_LINUX_H

#include <gtk/gtk.h>

int GdkKeyvalToGwk(guint keyval);
int GetGwkKey(GdkEventKey *e);
int GwkKeyToModifier(int glassKey);
int GdkModifierMaskToGwk(guint mask);
gint FindGdkKeyvalForGwkKeycode(int code);

#endif /* KEY_LINUX_H */
