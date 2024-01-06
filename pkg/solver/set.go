package solver

import (
	"github.com/RoaringBitmap/roaring"
)

// ParIntersect returns the intersection of all given units
func ParIntersect(units ...*roaring.Bitmap) *roaring.Bitmap {
	return roaring.ParAnd(0, units...)
}

// ParIntersectCells returns the intersection of all given cells
func ParIntersectCells(cells []*Cell) *roaring.Bitmap {
	marks := make([]*roaring.Bitmap, 0)
	for _, cell := range cells {
		marks = append(marks, roaring.BitmapOf(cell.Marks.ToArray()...))
	}
	return ParIntersect(marks...)
}

// ParUnion returns the union of all given units
func ParUnion(units ...*roaring.Bitmap) *roaring.Bitmap {
	return roaring.ParOr(0, units...)
}

// ParUnionCells returns the union of all given cells
func ParUnionCells(cells []*Cell) *roaring.Bitmap {
	marks := make([]*roaring.Bitmap, 0)
	for _, cell := range cells {
		marks = append(marks, roaring.BitmapOf(cell.Marks.ToArray()...))
	}
	return ParUnion(marks...)
}
