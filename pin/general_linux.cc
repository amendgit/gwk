#include "events.h"
#include "general_linux.h"

int GdkModifierMaskToGwk(guint mask) {
	int gwkMask = 0;
	if (mask & GDK_SHIFT_MASK) {
		gwkMask |= KeyEvent_ModifierShift;
	}
	if (mask & GDK_CONTROL_MASK) {
		gwkMask |= KeyEvent_ModifierControl;
	}
	if (mask & GDK_MOD1_MASK) {
		gwkMask |= KeyEvent_ModifierAlt;
	}
	if (mask & GDK_META_MASK) {
		gwkMask |= KeyEvent_ModifierAlt;
	}
	if (mask & GDK_BUTTON1_MASK) {
		gwkMask |= KeyEvent_ModifierButtonPrimary;
	}
	if (mask & GDK_BUTTON2_MASK) {
		gwkMask |= KeyEvent_ModifierButtonMiddle;
	}
	if (mask & GDK_BUTTON3_MASK) {
		gwkMask |= KeyEvent_ModifierButtonSecondary;
	}
	if (mask & GDK_SUPER_MASK) {
		gwkMask |= KeyEvent_ModifierWindows;
	}
	return gwkMask;
}
