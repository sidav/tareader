package middleware

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func FillPolygon(verts [][2]int32) {
	switch len(verts) {
	case 1, 2:
		return
	case 3:
		FillTriangleslope(verts[0][0], verts[0][1], verts[1][0], verts[1][1], verts[2][0], verts[2][1])
	case 4:
		FillTriangleslope(verts[0][0], verts[0][1], verts[1][0], verts[1][1], verts[2][0], verts[2][1])
		FillTriangleslope(verts[0][0], verts[0][1], verts[2][0], verts[2][1], verts[3][0], verts[3][1])
	default:
		for i := 1; i < len(verts)-1; i++ {
			FillTriangleslope(verts[0][0], verts[0][1], verts[i][0], verts[i][1], verts[i+1][0], verts[i+1][1])
		}
		// panic(fmt.Sprintf("Unimplemented number of vertices for polygon: %d\n", len(verts)))
	}
}

func hline(x1, x2, y int32) {
	rl.DrawLine(x1, y, x2, y, currColor)
}

func FillTriangleslope(x0, y0, x1, y1, x2, y2 int32) {
	var a, b, y, last int32
	// Sort coordinates by Y order (y2 >= y1 >= y0)
	if y0 > y1 {
		x0, y0, x1, y1 = x1, y1, x0, y0
	}
	if y1 > y2 {
		x2, y2, x1, y1 = x1, y1, x2, y2
	}
	if y0 > y1 {
		x0, y0, x1, y1 = x1, y1, x0, y0
	}

	if y0 == y2 { // All on same line case
		a = x0
		b = x0
		if x1 < a {
			a = x1
		} else if x1 > b {
			b = x1
		}
		if x2 < a {
			a = x2
		} else if x2 > b {
			b = x2
		}
		hline(a, b, y0)
		return
	}

	dx01 := x1 - x0
	dy01 := y1 - y0
	dx02 := x2 - x0
	dy02 := y2 - y0
	dx12 := x2 - x1
	dy12 := y2 - y1
	var sa, sb int32

	// For upper part of triangle, find scanline crossings for segment
	// 0-1 and 0-2.  If y1=y2 (flat-bottomed triangle), the scanline y
	// is included here (and second loop will be skipped, avoiding a /
	// error there), otherwise scanline y1 is skipped here and handle
	// in the second loop...which also avoids a /0 error here if y0=y
	// (flat-topped triangle)
	if y1 == y2 {
		last = y1 // Include y1 scanline
	} else {
		last = y1 - 1
	} // Skip it

	for y = y0; y <= last; y++ {
		a = x0 + sa/dy01
		b = x0 + sb/dy02
		sa += dx01
		sb += dx02
		// longhand a = x0 + (x1 - x0) * (y - y0) / (y1 - y0)
		//          b = x0 + (x2 - x0) * (y - y0) / (y2 - y0)
		hline(a, b, y)
	}

	// For lower part of triangle, find scanline crossings for segment
	// 0-2 and 1-2.  This loop is skipped if y1=y2
	sa = dx12 * (y - y1)
	sb = dx02 * (y - y0)
	for ; y <= y2; y++ {
		a = x1 + sa/dy12
		b = x0 + sb/dy02
		sa += dx12
		sb += dx02
		// longhand a = x1 + (x2 - x1) * (y - y1) / (y2 - y1)
		//          b = x0 + (x2 - x0) * (y - y0) / (y2 - y0)
		hline(a, b, y)
	}
}
