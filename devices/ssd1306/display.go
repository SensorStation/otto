package ssd1306

import (
	"image"
	"image/draw"
	"image/gif"
	"os"
	"time"

	"github.com/nfnt/resize"
	"github.com/sensorstation/otto/devices"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ssd1306"
	"periph.io/x/devices/v3/ssd1306/image1bit"
	"periph.io/x/host/v3"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type Bit bool

const (
	On  Bit = true
	Off Bit = false
)

type Display struct {
	bus        string
	addr       int
	Dev        *ssd1306.Dev
	Height     int
	Width      int
	Font       *basicfont.Face
	Background *image1bit.VerticalLSB
}

func NewDisplay(name string, width, height int) (*Display, error) {
	d := &Display{
		Height: height,
		Width:  width,
		bus:    "/dev/i2c-1",
		addr:   0x3c,
	}

	devices.NewI2CDevice(name, d.bus, d.addr)

	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	// Open a handle to the first available I²C bus:
	bus, err := i2creg.Open("/dev/i2c-1")
	if err != nil {
		return nil, err
	}

	// Open a handle to a ssd1306 connected on the I²C bus:
	opts := &ssd1306.DefaultOpts
	opts.H = height
	opts.W = width

	d.Dev, err = ssd1306.NewI2C(bus, opts)
	if err != nil {
		return nil, err
	}

	d.Background = image1bit.NewVerticalLSB(image.Rect(0, 0, opts.W, opts.H))
	return d, nil
}

func (d *Display) Clear() {
	// got to be a better way
	d.Rectangle(0, 0, d.Width, d.Height, Off)
	d.Draw()
}

func (d *Display) Draw() error {
	err := d.Dev.Draw(d.Background.Bounds(), d.Background, image.Point{})
	if err != nil {
		return err
	}
	return nil
}
func (d *Display) Rectangle(x0, y0, x1, y1 int, value Bit) {
	d.Clip(&x0, &y0, &x1, &y1)
	for x := x0; x < x1; x++ {
		for y := y0; y < y1; y++ {
			d.SetBit(x, y, value)
		}
	}
}

func (d *Display) Line(x0, y0, len, width int, value Bit) {
	x1 := x0 + len
	y1 := y0 + width
	d.Rectangle(x0, y0, x1, y1, value)
}

func (d *Display) SetBit(x, y int, value Bit) {
	d.Background.SetBit(x, y, image1bit.Bit(value))
}

func (d *Display) Clip(x0, y0, x1, y1 *int) {
	if x0 != nil {
		if *x0 < 0 {
			*x0 = 0
		}
		if *x0 > d.Width {
			*x0 = d.Width
		}
	}
	if x1 != nil {
		if *x1 < 0 {
			*x1 = 0
		}
		if *x1 > d.Width {
			*x1 = d.Width
		}
	}

	if y0 != nil {
		if *y0 < 0 {
			*y0 = 0
		}
		if *y0 > d.Height {
			*y0 = d.Height
		}
	}

	if y1 != nil {
		if *y1 < 0 {
			*y1 = 0
		}
		if *y1 > d.Height {
			*y1 = d.Height
		}
	}
}

func (d *Display) DrawString(x, y int, str string) {
	d.Font = basicfont.Face7x13
	drawer := &font.Drawer{
		Dst:  d.Background,
		Src:  image.White,
		Face: d.Font,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)},
	}
	drawer.DrawString(str)
}

func (d *Display) AnimatedGIF(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	g, err := gif.DecodeAll(f)
	f.Close()
	if err != nil {
		return err
	}

	// Converts every frame to image.Gray and resize them:
	imgs := make([]*image.Gray, len(g.Image))
	for i := range g.Image {
		imgs[i] = d.convertAndResizeAndCenter(g.Image[i])
	}

	// Display the frames in a loop:
	for i := 0; ; i++ {
		index := i % len(imgs)
		c := time.After(time.Duration(10*g.Delay[index]) * time.Millisecond)
		img := imgs[index]
		d.Dev.Draw(img.Bounds(), img, image.Point{})
		<-c
	}

	return nil
}

// convertAndResizeAndCenter takes an image, resizes and centers it on a
// image.Gray of size w*h.
func (d *Display) convertAndResizeAndCenter(src image.Image) *image.Gray {
	w := d.Width
	h := d.Height
	src = resize.Thumbnail(uint(w), uint(h), src, resize.Bicubic)
	img := image.NewGray(image.Rect(0, 0, w, h))
	r := src.Bounds()
	r = r.Add(image.Point{(w - r.Max.X) / 2, (h - r.Max.Y) / 2})
	draw.Draw(img, r, src, image.Point{}, draw.Src)
	return img
}
