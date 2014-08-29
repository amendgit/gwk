// go build -ldflags="-linkmode=external" -o ./a
#include "cocoa_gui.h"

int c_main(int argc, char *argv[]) {
	NSAutoreleasePool *pool = [[NSAutoreleasePool alloc] init];
	[NSApplication sharedApplication];
	int style = NSClosableWindowMask | NSResizableWindowMask | NSTexturedBackgroundWindowMask | NSTitledWindowMask | NSMiniaturizableWindowMask;
	NSWindow *win = [[NSWindow alloc] initWithContentRect:NSMakeRect(50, 50, 600, 400)
	styleMask:style
	backing:NSBackingStoreBuffered
	defer:NO];
	[win makeKeyAndOrderFront:win];
	[NSApp run];

	[pool release];

	return 0;
}