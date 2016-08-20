package ggk

// DrawLooper
// Subclasses of DrawLooper can be attached to a SkPaint. Where they are,  and
// something is drawn to a canvas with that paint, the looper subclass will be
// called, allowing it to modify the canvas and/or paint for that draw call.
// More than that, via the next() method, the looper can modify the draw to be
// invoked multiple times (hence the name loop-er), allow it to perform effects
// like shadows or frame/fills, that require more than one pass.
type DrawLooper struct {

}

// Called right before something is being drawn. Returns a Context
// whose next() method should be called until it returns false.
// The caller has to ensure that the storage pointer provides enough
// memory for the Context. The required size can be queried by calling
// contextSize(). It is also the caller's responsibility to destroy the
// object after use.
func (l *DrawLooper) CreateContext(canvas Canvas, storage interface{}) *DrawLooperContext {

}

// Returns the number of bytes needed to store subclasses of Context (belonging to the
// corresponding DrawLooper subclass).
func (l *DrawLooper) ContextSize() int {

}

// The fast bounds functions are used to enable the paint to be culled early
// in the drawing pipeline. If a subclass can support this feature it must
// return true for the canComputeFastBounds() function.  If that function
// returns false then computeFastBounds behavior is undefined otherwise it
// is expected to have the following behavior. Given the parent paint and
// the parent's bounding rect the subclass must fill in and return the
// storage rect, where the storage rect is with the union of the src rect
// and the looper's bounding rect.
func (l *DrawLooper) CanComputeFastBounds(paint *Paint) bool {

}

func (l *DrawLooper) ComputeFastBounds(paint *Paint, src, dst Rect) {

}

// If this looper can be interpreted as having two layers, such that
//     1. The first layer (bottom most) just has a blur and translate
//     2. The second layer has no modifications to either paint or canvas
//     3. No other layers.
// then return true, and if not null, fill out the BlurShadowRec).
//
// If any of the above are not met, return false and ignore the BlurShadowRec parameter.
func (l *DrawLooper) AsABlurShadow(rec *DrawLooperBlurShadowRec) bool {

}