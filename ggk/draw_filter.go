package ggk

type DrawFilterType int

const (
	KDrawFilterTypePaint DrawFilterType = iota
	KDrawFilterTypePoint
	KDrawFilterTypeLine
	KDrawFilterTypeBitmap
	KDrawFilterTypeRect
	KDrawFilterTypeRRect
	KDrawFilterTypeOval
	KDrawFilterTypePath
	KDrawFilterTypeText
)

const KDrawFilterTypeCount int = int(KDrawFilterTypeText) + 1

// Right before something is being draw, Filter() is called with the
// paint. The filter may modify the paint as it wishes, which will then be
// used for the actual drawing. Note: this modification only lasts for the
// current draw, as a temporary copy of the paint is used.
type DrawFilter interface {
	// Called with the paint that will be used to draw the specified type.
	// The implementation may modify the paint as they wish. If Filter()
	// returns false, the draw will be skipped.
	Filter(*Paint, DrawFilterType) bool
}
