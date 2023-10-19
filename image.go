package sourceafis

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"sourceafis/primitives"

	"github.com/jtejido/go-wsq"
)

type ImageOptions func(*Image) (*Image, error)

func WithResolution(dpi float64) ImageOptions {
	return func(f *Image) (*Image, error) {
		if dpi < 20 || dpi > 20_000 {
			return nil, fmt.Errorf("dpi cannot be less than 20 and greater than 20000")
		}
		f.dpi = dpi
		return f, nil
	}
}

type Image struct {
	dpi    float64
	matrix *primitives.Matrix
}

func (i *Image) Matrix() *primitives.Matrix {
	return i.matrix
}

func New(bounds image.Rectangle, opts ...ImageOptions) (img *Image, err error) {
	dx, dy := bounds.Dx(), bounds.Dy()
	img = new(Image)
	for _, opt := range opts {
		img, err = opt(img)
		if err != nil {
			return
		}
	}
	if img.dpi == 0 {
		img.dpi = 500
	}
	img.matrix = primitives.NewMatrix(dx, dy)
	return
}

type Pixel struct {
	R int
	G int
	B int
	A int
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func NewFromImage(img image.Image, opts ...ImageOptions) (*Image, error) {
	bounds := img.Bounds()
	m, err := New(bounds, opts...)
	if err != nil {
		return nil, err
	}
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			pix := rgbaToPixel(img.At(x, y).RGBA())
			colorSum := pix.R + pix.G + pix.B
			m.matrix.Set(x, y, 1-float64(colorSum)*(1.0/(3.0*255.0)))
		}
	}
	return m, nil
}

func NewFromGray(img *image.Gray, opts ...ImageOptions) (*Image, error) {
	bounds := img.Bounds()
	m, err := New(bounds, opts...)
	if err != nil {
		return nil, err
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			m.matrix.Set(x, y, float64(img.GrayAt(x, y).Y))
		}
	}
	return m, nil
}

func LoadImage(fname string, opts ...ImageOptions) (*Image, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var img image.Image

	ext := path.Ext(fname)
	if ext == ".wsq" {
		img, err = wsq.Decode(f)
	} else if ext == ".jpg" {
		img, err = jpeg.Decode(f)
	} else if ext == ".png" {
		img, err = png.Decode(f)
	} else {
		return nil, fmt.Errorf("%q extension not supported", ext)
	}
	if err != nil {
		return nil, fmt.Errorf("cannot decode image %s, err: %s", fname, err)
	}

	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			gray.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}

	return NewFromImage(img, opts...)
}
