// go build -ldflags="-linkmode=external" -o ./a
#import <AppKit/AppKit.h>

@interface GWKWindow : NSObject

- (void)run;
- (void)setDelegate:(id)delegate;

@end

@implementation GWKWindow {
	id delegate_; 
}

// - (NSSize)window:(NSWindow*)window willUseFullScreenContentSize:(NSSize)contentSize {
// 	// window_will_use_fullscreen_content_rect(host_window, window, &content_rect);
// }

- (void)run {
	NSAutoreleasePool *pool = [[NSAutoreleasePool alloc] init];
	[NSApplication sharedApplication];

	int style = NSClosableWindowMask | NSResizableWindowMask | NSTexturedBackgroundWindowMask | NSTitledWindowMask | NSMiniaturizableWindowMask;
	
	NSWindow *win = [[NSWindow alloc] initWithContentRect:NSMakeRect(50, 50, 600, 400) styleMask:style backing:NSBackingStoreBuffered defer:NO];
	// [win setDelegate:self];
	[win makeKeyAndOrderFront:win];
	
	[NSApp run];

	[pool release];
}

- (void)setDelegate:(id)delegate {
	delegate_ = delegate;
}

@end

void *new_gwk_window() {
	return [[GWKWindow alloc] init];
}

void gwk_window_set_delegate(void* slf, void* delegate) {
	[(id)slf setDelegate:(id)delegate];
}

void gwk_window_run(id slf) {
	[(id)slf run];
}

