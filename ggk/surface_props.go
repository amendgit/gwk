package ggk

// SurfaceProps describe of how the LCD strips are arranged for each pixel. If
// this is unknown, or the pixels are meant to be "protable" and/or transformed
// before showing (e.g. rotated, scaled) then use KPixelsGeometryUnknown.
type PixelGeometry int

const (
	KPixelGeometryUnknown = iota
	KPixelGeometryRGBH
	KPixelGeometryBGRH
	KPixelGeometryRGBV
	KPixelGeometryBGRV
)

type SurfacePropsFlags int

const (
	KSurfacePropsFlagNone              = 0
	KSurfacePropsFlagDisallowAntiAlias = 1 << iota
	KSurfacePropsFlagDisallowDither
	KSurfacePropsFlagUseDistanceFieldFonts
)

type SurfacePropsInitType int

const (
	KSurfacePropsInitTypeNone = iota
	KSurfacePropsInitTypeLegacyFontHost
)

type SurfaceProps struct {
	flags         SurfacePropsFlags
	pixelGeometry PixelGeometry
}

func MakeSurfaceProps(flags SurfacePropsFlags, initType SurfacePropsInitType) SurfaceProps {
	var props SurfaceProps
	props.flags = flags
	if initType == KSurfacePropsInitTypeLegacyFontHost {
		props.pixelGeometry = props.computeDefaultGeometry()
	}
	return props
}

func (props *SurfaceProps) computeDefaultGeometry() PixelGeometry {
	toimpl()
	return KPixelGeometryBGRH
}
