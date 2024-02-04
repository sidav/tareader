package renderer

func (r *ModelRenderer) DrawLine(fromx, fromy, tox, toy int32) {
	deltax := fastInt32Abs(tox - fromx)
	deltay := fastInt32Abs(toy - fromy)
	var xmod int32 = 1
	var ymod int32 = 1
	if tox < fromx {
		xmod = -1
	}
	if toy < fromy {
		ymod = -1
	}
	var error int32 = 0
	if deltax >= deltay {
		y := fromy
		deltaerr := deltay
		for x := fromx; x != tox+xmod; x += xmod {
			r.gAdapter.DrawPoint(x, y)
			error += deltaerr
			if 2*error >= deltax {
				y += ymod
				error -= deltax
			}
		}
	} else {
		x := fromx
		deltaerr := deltax
		for y := fromy; y != toy+ymod; y += ymod {
			r.gAdapter.DrawPoint(x, y)
			error += deltaerr
			if 2*error >= deltay {
				x += xmod
				error -= deltay
			}
		}
	}
}

func fastInt32Abs(n int32) int32 {
	y := n >> 31
	return (n ^ y) - y
}
