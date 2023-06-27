package engine

import (
	"math"

	"github.com/jakecoffman/cp"
)

func Cpvnormalize(v cp.Vector) cp.Vector {
	length := Cpvlength(v)
	if length == 0 {
		return cp.Vector{}
	}
	return cp.Vector{X: v.X / length, Y: v.Y / length}
}

func Cpvlength(v cp.Vector) float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
}

// cpvmult multiplies a cpVect by a scalar value and returns the result as a new cpVect.
func Cpvmult(v cp.Vector, s float64) cp.Vector {
	return cp.Vector{X: v.X * s, Y: v.Y * s}
}

// cpvadd adds two cpVects together and returns the result as a new cpVect.
func Cpvadd(v1, v2 cp.Vector) cp.Vector {
	return cp.Vector{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

// cpvsub subtracts two cpVects and returns the result as a new cpVect.
func Cpvsub(v1, v2 cp.Vector) cp.Vector {
	return cp.Vector{X: v1.X - v2.X, Y: v1.Y - v2.Y}
}

func Cpvneg(v cp.Vector) cp.Vector {
	return cp.Vector{X: -v.X, Y: -v.Y}
}

func DistanceAndOverlap(rect1, rect2 cp.BB) (float64, float64) {
	center1 := rect1.Center()
	center2 := rect2.Center()
	distance := center1.Distance(center2)

	overlapWidth := math.Min(rect1.R, rect2.R) - math.Max(rect1.L, rect2.L)
	overlapHeight := math.Min(rect1.T, rect2.T) - math.Max(rect1.B, rect2.B)

	overlapArea := 0.0
	if overlapWidth > 0 && overlapHeight > 0 {
		overlapArea = overlapWidth * overlapHeight
	}

	rect1Area := rect1.Area()
	rect2Area := rect2.Area()
	totalArea := rect1Area + rect2Area - overlapArea

	overlapPercentage := 0.0
	if totalArea > 0 {
		overlapPercentage = (overlapArea / totalArea) * 100
	}

	return distance, overlapPercentage
}

func BBIntersects(bb1, bb2 cp.BB) bool {
	return bb1.L <= bb2.R && bb1.R >= bb2.L && bb1.B <= bb2.T && bb1.T >= bb2.B
}

func OverlapSide(rect1, rect2 cp.BB) cp.Vector {
	overlapWidth := math.Min(rect1.R, rect2.R) - math.Max(rect1.L, rect2.L)
	overlapHeight := math.Min(rect1.T, rect2.T) - math.Max(rect1.B, rect2.B)

	if overlapWidth <= 0 || overlapHeight <= 0 {
		return cp.Vector{}
	}

	if overlapWidth < overlapHeight {
		if rect1.R < rect2.R {
			return cp.Vector{X: -1, Y: 0}
		} else {
			return cp.Vector{X: 1, Y: 0}
		}
	} else {
		if rect1.T < rect2.T {
			return cp.Vector{X: 0, Y: -1}
		} else {
			return cp.Vector{X: 0, Y: 1}
		}
	}
}

func Lerp(f1, f2, t float64) float64 {
	return f1*(1-t) + f2*t
}

func VecLerp(v1, v2 cp.Vector, t float64) cp.Vector {
	return cp.Vector{
		X: v1.X + (v2.X-v1.X)*t,
		Y: v1.Y + (v2.Y-v1.Y)*t,
	}
}
