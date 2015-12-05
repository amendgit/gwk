#include "events.h"

const char *WindowEventName(int eventType) {
	switch (eventType) {
		case kWindowEventResize:
			return "RESIZE";
		default:
			return "UNKNOWN";
	}
}
