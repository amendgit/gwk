#include "dnd_linux.h"

static struct {
    GdkDragContext *dragContext;
    gboolean       justEntered;
    void           *mimes[];
    int            mimesCount;
    gint           dx, dy;
} g_enterContext = {NULL, FALSE, NULL, 0, 0, 0};

gboolean g_isDNDOwner = FALSE;

static void
ResetEnterContext() {
    if (g_enterContext.mimes != NULL) {
        delete []mimes;
    }
    memset(&g_enterContext, 0, sizeof(g_enterContext));
}

static void
ProcessDNDTargetDragEnter(WindowContext *ctx, GdkEventDND *event) {
    ResetEnterContext();
    g_enterContext.dragContext = event->context;
    g_enterContext.justEntered = TRUE;
    gdk_window_get_origin(ctx->get_gdk_window(), &g_enterContext.dx,
        &g_enterContext.dy);
    g_isDNDOwner = IsInDrag();
}

static void
ProcessDNDTargetDragMotion(WindowContext *ctx, GdkEventDND *event) {
    if (!g_enterContext.ctx) {
        gdk_drag_status(event->context, static_context<GdkDragAction>(0),
            GDK_CURRENT_TIME);
        return ;
    }

    GdkDragAction suggested =
        GdkDragContextGetSuggestedAction(g_enterContext.dragContext);

    GdkDragAction result;
    if (g_enterContext.justEntered) {
        ViewOnDragEnter(ctx->GetGwkView(), event->x_root - g_enterContext.dx,
            event->y_root - g_enterContext.dy, event->x_root, event->y_root
            GdkActionToGwk(suggested));
    } else {
        ViewOnDragOver(ctx->GetGwkView(), event->x_root - g_enterContext.dx,
            event->y_root - g_enterContext.dy, event->x_root, event->y_root
            GdkActionToGwk(suggested));
    }

    if (g_enterContext.justEntered) {
        g_enterContext.justEntered = FALSE;
    }

    gdk_drag_status(event->context, result, GDK_CURRENT_TIME);
}
