package graphicadapter

type GraphicBackend interface {
	Init(int32, int32)
	GetRealResolution() (int32, int32)
	GetRenderResolution() (int32, int32)
	BeginFrame()
	EndFrame()
	SetInternalResolution(w, h int32)
	Flush()
	Clear()
	SetColor(r, g, b uint8)
	DrawPoint(x, y int32)
	FillRect(x, y, w, h int)
	VerticalLine(x, y0, y1 int)
}
