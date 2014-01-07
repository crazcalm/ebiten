package texture

import (
	"github.com/hajimehoshi/go-ebiten/graphics"
	"image"
	"image/draw"
)

func nextPowerOf2(x uint64) uint64 {
	x -= 1
	x |= (x >> 1)
	x |= (x >> 2)
	x |= (x >> 4)
	x |= (x >> 8)
	x |= (x >> 16)
	x |= (x >> 32)
	return x + 1
}

type Texture struct {
	native interface{}
	width  int
	height int
}

func AdjustSize(size int) int {
	return int(nextPowerOf2(uint64(size)))
}

func AdjustImage(img image.Image) *image.NRGBA {
	width, height := img.Bounds().Size().X, img.Bounds().Size().Y
	adjustedImageBounds := image.Rectangle{
		image.ZP,
		image.Point{
			AdjustSize(width),
			AdjustSize(height),
		},
	}
	if nrgba, ok := img.(*image.NRGBA); ok &&
		img.Bounds() == adjustedImageBounds {
		return nrgba
	}

	adjustedImage := image.NewNRGBA(adjustedImageBounds)
	dstBounds := image.Rectangle{
		image.ZP,
		img.Bounds().Size(),
	}
	draw.Draw(adjustedImage, dstBounds, img, image.ZP, draw.Src)
	return adjustedImage
}

func New(native interface{}, width, height int) *Texture {
	return &Texture{native, width, height}
}

func (texture *Texture) u(x int) float32 {
	return float32(x) / float32(AdjustSize(texture.width))
}

func (texture *Texture) v(y int) float32 {
	return float32(y) / float32(AdjustSize(texture.height))
}

type Drawable interface {
	Draw(native interface{}, quads []graphics.TextureQuad)
}

func (texture *Texture) Draw(drawable Drawable) {
	x1 := float32(0)
	x2 := float32(texture.width)
	y1 := float32(0)
	y2 := float32(texture.height)
	u1 := texture.u(0)
	u2 := texture.u(texture.width)
	v1 := texture.v(0)
	v2 := texture.v(texture.height)
	quad := graphics.TextureQuad{x1, x2, y1, y2, u1, u2, v1, v2}
	drawable.Draw(texture.native, []graphics.TextureQuad{quad})
}

func (texture *Texture) DrawParts(parts []graphics.TexturePart, drawable Drawable) {
	quads := []graphics.TextureQuad{}
	for _, part := range parts {
		x1 := float32(part.LocationX)
		x2 := float32(part.LocationX + part.Source.Width)
		y1 := float32(part.LocationY)
		y2 := float32(part.LocationY + part.Source.Height)
		u1 := texture.u(part.Source.X)
		u2 := texture.u(part.Source.X + part.Source.Width)
		v1 := texture.v(part.Source.Y)
		v2 := texture.v(part.Source.Y + part.Source.Height)
		quad := graphics.TextureQuad{x1, x2, y1, y2, u1, u2, v1, v2}
		quads = append(quads, quad)
	}
	drawable.Draw(texture.native, quads)
}

type FramebufferCreator interface {
	Create(native interface{}) interface{}
}

func (texture *Texture) CreateRenderTarget(creator FramebufferCreator) *RenderTarget {
	return NewRenderTarget(creator.Create(texture.native), texture.width, texture.height)
}
