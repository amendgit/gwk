#ifndef WINDOW_H
#define WINDOW_H

#include <stdbool.h>
#include <gtk/gtk.h>
#include <X11/Xlib.h>

#include <set>
#include <vector>

#include "_cgo_export.h"

enum WindowFrameType {
	WindowFrameType_Titled,
	WindowFrameType_Untitled,
	WindowFrameType_Transparent,
};

enum WindowType {
	WindowType_Normal,
	WindowType_Utility,
	WindowType_Popup,
};

enum RequestType {
	RequestType_None,
	RequestType_Resizable,
	RequestType_NotResizable,
};

struct WindowFrameExtents {
	int top;
	int left;
	int width;
	int height;
};

static const guint MOUSE_BUTTONS_MASK = (guint)(GDK_BUTTON1_MASK |GDK_BUTTON2_MASK | GDK_BUTTON3_MASK);

enum BoundsType {
	BOUNDSTYPE_CONTENT,
	BOUNDSTYPE_WINDOW,
};

struct WindowGeometry {
	WindowGeometry(): finalWidth(), finalHeight(), refX(), refY(), gravityX(), 
	gravityY(), currentWidth(), currentHeight(), extents() {}

	// estimate of the final width the window will get after all pending 
	// configure requests are processed by the window manager.
	struct {
		int value;
		BoundsType type;
	} finalWidth, finalHeight;

	float refX;
	float refY;
	float gravityX;
	float gravityY;

	// the last width which was configured or obtained from configure
	// notification
	int currentWidth;

	// the last height which was configured or obtained from configure
	// notification
	int currentHeight;

	WindowFrameExtents extents;
};

class WindowContextChild;
class WindowContextTop;

class WindowContext {
public:
	virtual bool IsEnabled() = 0;
	virtual bool HasIME() = 0;
	virtual bool FilterIME(GdkEvent *) = 0;
	virtual void EnableOrResetIME() = 0;
	virtual void DisableIME() = 0;
	virtual void Paint(void *data, int width, int height) = 0;
	virtual WindowFrameExtents GetFrameExtents() = 0;

	virtual void EnterFullScreen() = 0;
	virtual void ExitFullScreen() = 0;
	virtual void SetVisible(bool) = 0;
	virtual bool IsVisible() = 0;
	virtual void SetBounds() = 0;
	virtual void SetResizable() = 0;
	virtual void ReuqestFocus() = 0;
	virtual void SetFocusable() = 0;
	virtual bool GrabFocus() = 0;
	virtual bool GrabMouseDragFocus() = 0;
	virtual void UngrabFocus() = 0;
	virtual void UngrabMouseDragFocus() = 0;
	virtual void SetTitle(const char *) = 0;
	virtual void SetAlpha(double) = 0;
	virtual void SetEnabled(bool) = 0;
	virtual void SetMinimumSize(int, int) = 0;
	virtual void SetMaximumSize(int, int) = 0;
	virtual void SetMinimized(bool) = 0;
	virtual void SetMaximuzed(bool) = 0;
	virtual void SetIcon(GdkPixbuf *) = 0;
	virtual void Restack(bool) = 0;
	virtual void SetCursor(GdkCursor *) = 0;
	virtual void SetModal(bool, WindowContext *parent = NULL) = 0;
	virtual void SetGravity(float, float) = 0;
	virtual void SetLevel(int) = 0;
	virtual void SetBackground(float, float, float) = 0;

	virtual void ProcessPropertyNotify(GdkEventProperty*) = 0;
	virtual void ProcessConfigure(GdkEventConfigure *) = 0;
	virtual void ProcessMap() = 0;
	virtual void ProcessFocus(GdkEventFocus *) = 0;
	virtual void ProcessDestroy() = 0;
	virtual void ProcessDelete() = 0;
	virtual void ProcessExpose(GdkEventExpose *) = 0;
	virtual void ProcessMouseButton(GdkEventButton *) = 0;
	virtual void ProcessMouseMotion(GdkEventMotion *) = 0;
	virtual void ProcessMouseScroll(GdkEventScroll *) = 0;
	virtual void ProcessMouseCross(GdkEventCrossing *) = 0;
	virtual void ProcessKey(GdkEventKey *) = 0;
	virtual void ProcessState(GdkEventWindowState *) = 0;

	virtual void NotifyState(int) = 0;

	virtual void AddChild(WindowContextTop *child) = 0;
	virtual void RemoveChild(WindowContextTop *child) = 0;
	virtual void SetView(GoObject) = 0;

	virtual GdkWindow *GetGdkWindow() = 0;
	virtual GtkWindow *GetGtkWindow() = 0;
	virtual GoObject GetGwkView() = 0;
	virtual GoObject GetGwkWindow() = 0;

	virtual int GetEmbeddedX() = 0;
	virtual int GetEmbeddedY() = 0;

	virtual void IncrementEventsCounter() = 0;
	virtual void DecrementEventsCounter() = 0;
	virtual size_t EventsCount() = 0;
	virtual bool IsDead() = 0;
	virtual ~WindowContext() {}
};

class WindowContextBase: public WindowContext {
	std::set<WindowContextTop *> children;

	struct _XIM {
		XIM im;
		XIC ic;
		bool enabled;
	} xim;

	size_t eventsProcessingCount;
	bool canBeDeleted;

protected:
	GoObject gwkWindow;
	GoObject gwkView;
	GtkWidget *gtkWidget;
	GdkWindow *gdkWindow;

	bool isIconified;
	bool isMaximized;
	bool isMouseEntered;

	// smGrabWindow points to a WindowContext holding a mouse grab.
	// It is mostly used for popup windows.
	static WindowContext *smGrabWindow;

	// smMouseDragWindow points to a WindowContext from which a mouse drag 
	// strated. This WindowContext holding a mouse grab during this drag. After
	// releasing all mouse buttons smMouseDragWindow becomes NULL and 
	// smGrabWindow's mouse grab should be restored if present.
	//
	// This is done in order to mimic Windows behavior:
	// All mouse events should be delivered to a window from which mouse drag
	// started, until all mouse buttons released. No mouse ENTER/EXIT events
	// should be reported during this drag.
	static WindowContext *smMouseDragWindow;

public:
	bool IsEnabled();
	bool HasIME();
	bool FilterIME(GdkEvent *);
	void EnableOrResetIME();
	void DisableIME();
	void Paint(void *, int, int);
	GdkWindow  *GetGdkWindow();
	GoObject GetGwkWindow();
	GoObject GetGwkView();

	void AddChild(WindowContextTop *);
	void RemoveChild(WindowContextTop *);
	void ShowOrHideChildren(bool);
	void ReparentChildren(WindowContext *parent);
	void SetVisible(bool);
	bool IsVisible();
	void SetView(GoObject);
	bool GrabFocus();
	bool GrabMouseDragFocus();
	void UngrabFocus();
	void UngrabMouseDragFocus();
	void SetCursor(GdkCursor *);
	void SetLevel(int) {};
	void SetBackground(float, float, float);

	void ProcessMap() {}
	void ProcessFocus(GdkEventFocus *);
	void ProcessDestroy();
	void ProcessDelete();
	void ProcessExpose(GdkEventExpose *);
	void ProcessMouseButton(GdkEventButton *);
	void ProcessMouseMotion(GdkEventMotion *);
	void ProcessMouseScroll(GdkEventScroll *);
	void ProcessMouseCross(GdkEventCrossing *);
	void ProcessKey(GdkEventKey *event);
	void ProcessState(GdkEventWindowState *);

	void NotifyState(int);

	int GetEmbeddedX() { return 0; }
	int GetEmbeddedY() { return 0; }

	void IncrementEventsCounter();
	void DecrementEventsCounter();
	size_t EventsCount();
	bool IsDead();

	~WindowContextBase();

protected:
	virtual void ApplyShapeMask(void *, uint width, uint height) = 0;

private:
	bool IMFilterKeypress(GdkEventKey *);
};

class WindowContextTop : public WindowContextBase {
	GoObject screen;
	WindowFrameType frameType;
	WindowContext *owner;
	// WindowGeometry geometry;
	int stateConfigNotifications;
	struct Resizable {
		Resizable(): requestType(RequestType_None), value(true), prev(false), 
			minW(-1), minH(-1), maxW(-1), maxH(-1) {
			// empty		
		}
		RequestType requestType;
		bool value;
		bool prev;
		int minW, minH, maxW, maxH;
	} resizable;
	bool frameExtentsInitialized;
	bool mapReceived;
	bool locationAssigned;
	bool sizeAssigned;

public:
	WindowContextTop(GoObject gwkWindow, WindowContext *owner, GoObject screen, WindowFrameType frameType, WindowType windowType);
	
	WindowFrameExtents GetFrameExtents();
	void EnterFullScreen();
	void ExitFullScreen();
	void SetBounds();
	void SetResizable();
	void ReuqestFocus();
	void SetFocusable();
	bool GrabFocus();
	bool GrabMouseDragFocus();
	void UngrabFocus();
	void UngrabMouseDragFocus();
	void SetTitle(const char *);
	void SetAlpha(double);
	void SetEnabled(bool);
	void SetMinimumSize(int, int);
	void SetMaximumSize(int, int);
	void SetMinimized(bool);
	void SetMaximuzed(bool);
	void SetIcon(GdkPixbuf *);
	void Restack(bool);
	void SetCursor(GdkCursor *);
	void SetModal(bool, WindowContext *parent = NULL);
	void SetGravity(float, float);
	void ProcessPropertyNotify(GdkEventProperty*);
	void ProcessConfigure(GdkEventConfigure *);
	GtkWindow *GetGtkWindow(); // TODO: get window from parent.
	void ApplyShapeMask(void *, uint width, uint height) {}
};

#endif