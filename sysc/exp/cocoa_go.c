#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

bool NSApp·isRunning() {
	return [NSApp isRunning];
}

void NSApp·sendEvent(id event) {
	[NSApp sendEvent:id];
}

NSEvent *NSApp·nextEvent(NSUInteger matchingMask, NSDate* untilDate, NSString* inMode, bool dequeue) {
	return [NSApp nextEventMachingMask:matchingMask untilDate:untilDate inMode:inMode dequeue:dequeue];
}