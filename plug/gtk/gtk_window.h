#include <gtk/gtk.h>
#include <X11/Xlib.h>

#include "glass_view.h"

typedef enum {
	TITLED,
	UNTITLED,
	TRANSPARENT
} WindowFrameType;

typedef enum {
	NORMAL,
	UTILITY,
	POPUP
} WindowType;

typedef enum {
	REQUEST_NONE,
	REQUEST_RESIZABLE,
	REQUEST_NOT_RESIZABLE
} request_type;

typedef struct {
	int top;
	int left;
	int bottom;
	int right;
} WindowFrameExtents;

static const guint MOUSE_BUTTONS_MASK = (guint) (GDK_BUTTON1_MASK | GDK_BUTTON2_MASK | GDK_BUTTON3_MASK);

typedef enum {
	BOUNDSTYPE_CONTENT,
	BOUNDSTYPE_WINDOW
} BoundsType;

typedef struct {
	struct {
		int        value;
		BoundsType type;
	} finalWidth;

	struct {
		int        value;
		BoundsType type;
	} finalHeight;

	float refx;
	float refy;
	float gravityX;
	float gravityY;

	// the last width which was configured or obtained from configure
	// notification
	int currentWidth;

	// the last height which was configured or obtainted from configure
	// notification
	int currentHeight;

	WindowFrameExtents extents;
} WindowGeometry;