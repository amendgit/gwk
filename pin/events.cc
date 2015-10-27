#include "events.h"

const char *WindowEventName(int eventType) {
	switch (eventType) {
		case WindowEvent_Resize:
			return "RESIZE";
		default:
			return "UNKNOWN";
	}
}
