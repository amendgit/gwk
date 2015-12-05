package ggk

import (
	"log"
	"math"
)

type Scalar float32

const (
	KScalar1          Scalar = 1.0
	KScalarHalf       Scalar = 0.5
	KScalarSqrt2      Scalar = 1.41421356
	KScalarPI         Scalar = 3.14159265
	KScalarTanPIOver8 Scalar = 0.414213562
	KScalarRoot2Over2 Scalar = 0.707106781
	KScalarMax        Scalar = 3.402823466e+38
	KScalarMin        Scalar = -KScalarMax
	KScalarNearlyZero Scalar = KScalar1 / (1 << 12)
)

func ScalarInfinity() Scalar {
	return Scalar(math.Float32frombits(0x7F800000)) // IEEE infinity
}

func ScalarNegativeInfinity() Scalar {
	return Scalar(math.Float32frombits(0xFF800000)) // IEEE negative infinity
}

func ScalarNaN() Scalar {
	return Scalar(math.Float32frombits(0x7FFFFFFF)) // IEEE not a number
}

func ScalarFloor(x Scalar) Scalar {
	return Scalar(math.Floor(float64(x)))
}

func ScalarCeil(x Scalar) Scalar {
	return Scalar(math.Ceil(float64(x)))
}

func ScalarRound(x Scalar) Scalar {
	if x < 0 {
		return Scalar(math.Ceil(float64(x - 0.5)))
	}
	return Scalar(math.Floor(float64(x + 0.5)))
}

func ScalarFloorToInt(x Scalar) int {
	return int(ScalarFloor(x))
}

func ScalarCeilToInt(x Scalar) int {
	return int(ScalarCeil(x))
}

func ScalarRoundToInt(x Scalar) int {
	return int(ScalarRound(x))
}

func ScalarAbs(x Scalar) Scalar {
	return Scalar(math.Abs(float64(x)))
}

func ScalarCopysign(x, y Scalar) Scalar {
	return Scalar(math.Copysign(float64(x), float64(y)))
}

func ScalarMod(x, y Scalar) Scalar {
	return Scalar(math.Mod(float64(x), float64(y)))
}

func ScalarFraction(x Scalar) Scalar {
	return ScalarMod(x, 1.0)
}

func ScalarSqrt(x Scalar) Scalar {
	return Scalar(math.Sqrt(float64(x)))
}

func ScalarPow(b, e Scalar) Scalar {
	return Scalar(math.Pow(float64(b), float64(e)))
}

func ScalarSin(x Scalar) Scalar {
	return Scalar(math.Sin(float64(x)))
}

func ScalarCos(x Scalar) Scalar {
	return Scalar(math.Cos(float64(x)))
}

func ScalarTan(x Scalar) Scalar {
	return Scalar(math.Tan(float64(x)))
}

func ScalarAsin(x Scalar) Scalar {
	return Scalar(math.Asin(float64(x)))
}

func ScalarAcos(x Scalar) Scalar {
	return Scalar(math.Acos(float64(x)))
}

func ScalarAtan2(y, x Scalar) Scalar {
	return Scalar(math.Atan2(float64(y), float64(x)))
}

func ScalarExp(x Scalar) Scalar {
	return Scalar(math.Exp(float64(x)))
}

func ScalarLog(x Scalar) Scalar {
	return Scalar(math.Log(float64(x)))
}

func ScalarLog2(x Scalar) Scalar {
	return Scalar(math.Log2(float64(x)))
}

func ScalarFromInt(i int) Scalar {
	return Scalar(i)
}

func ScalarTruncToInt(x Scalar) int {
	return int(x)
}

func ScalarFromFloat32(f float32) Scalar {
	return Scalar(f)
}

func ScalarToFloat32(x Scalar) float32 {
	return float32(x)
}

func ScalarFromFloat64(f float64) Scalar {
	return Scalar(f)
}

func ScalarToFloat64(x Scalar) float64 {
	return float64(x)
}

func ScalarIsNaN(x Scalar) bool {
	return x != x
}

// Returns true if x is not NaN and not infinite
func ScalarIsFinite(x Scalar) bool {
	// We rely on the following behavior of infinites and nans
	// 0 * finite --> 0
	// 0 * infinity --> NaN
	// 0 * NaN --> NaN
	var prod Scalar = x * 0
	// At this point, prod will either be NaN or 0
	return !ScalarIsNaN(prod)
}

func ScalarsAreFinite(array []Scalar) bool {
	var prod Scalar = 0
	for i := 0; i < len(array); i++ {
		prod = prod * array[i]
	}
	// At this point, prod will either be NaN or 0
	return !ScalarIsNaN(prod)
}

func ScalarMul(a, b Scalar) Scalar {
	return a * b
}

func ScalarMulAdd(a, b, c Scalar) Scalar {
	return a*b + c
}

func ScalarMulDiv(a, b, c Scalar) Scalar {
	return a * b / c
}

func ScalarInvert(x Scalar) Scalar {
	return KScalar1 / x
}

func ScalarAverage(a, b Scalar) Scalar {
	return (a + b) * KScalarHalf
}

func ScalarHalf(x Scalar) Scalar {
	return x * KScalarHalf
}

func DegreesToRadians(degrees float32) float32 {
	return degrees * float32(KScalarPI) / 180
}

func RadiansToDegrees(radians float32) float32 {
	return radians * 180 / float32(KScalarPI)
}

func ScalarMax(a, b Scalar) Scalar {
	if a > b {
		return a
	}
	return b
}

func ScalarMin(a, b Scalar) Scalar {
	if a > b {
		return b
	}
	return a
}

func ScalarIsInteger(x Scalar) bool {
	return x == Scalar(int(x))
}

// Returns -1, 0 or 1 depending on the sign of the value:
// -1 if x < 0
// 0 if x == 0
// 1 if x > 0
func ScalarSignAsInt(x Scalar) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}

// Scalar result version of above
func ScalarSignAsScalar(x Scalar) Scalar {
	if x > 0 {
		return KScalar1
	} else if x < 0 {
		return -KScalar1
	}
	return Scalar(0)
}

func ScalarNearlyZero(x, tolerance Scalar) bool {
	return ScalarAbs(x) <= tolerance
}

func ScalarNearlyEqual(x, y, tolerance Scalar) bool {
	return ScalarAbs(x-y) <= tolerance
}

// Linearly interpolate between A and B, based on t.
// If t is 0, return A
// If t is 1, return B
// else interpolate.
// t must be [0..Scalar_1]
func ScalarInterpolate(x, y, t Scalar) Scalar {
	return x + (y-x)*t
}

// Interpolate along the function described by (keys[length], values[length])
// for the passed searchKey. SearchKeys outside the range keys[0]-keys[length]
// clamp to the min or max value. This function was inspired by a desire to
// change the multiplier for thickness in fakeBold; therefore it assumes the
// number of pairs (length) will be small, and a linear search is used. Repeated
// keys are allowed for dicountinuous functions (so long as keys is monotonically
// increasing), and if key is the value of a repeated scalar in keys, the first
// one will be used. However, that may change if a binary search is used.
func ScalarInterpolateFunc(searchKey Scalar, keys, values []Scalar) Scalar {
	length := len(keys)

	for i := 1; i < length; i++ {
		if keys[i] < keys[i-1] {
			log.Printf("ScalarInterpolateFunc: keys should be increasing.")
			return 0
		}
	}

	var right int = 0
	for right < length && searchKey > keys[right] {
		right++
	}

	// Could use sentinel values to eliminate conditionals, but since the
	// tables are taken as input, a simpler format is better.
	if length == right {
		return values[length-1]
	}

	if right == 0 {
		return values[0]
	}

	// Otherwise. interpolate between right - q and right.
	var (
		rightKey Scalar = keys[right]
		leftKey  Scalar = keys[right-1]
		fract    Scalar = (searchKey - leftKey) / (rightKey - leftKey)
	)

	return ScalarInterpolate(values[right-1], values[right], fract)
}

// Helper to compare an array of scalars.
func ScalarsEqual(a, b []Scalar, n int) bool {
	for i := 0; i < n; i++ {
		if a[i] != b[1] {
			return false
		}
	}
	return true
}
