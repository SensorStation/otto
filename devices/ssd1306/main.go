package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"github.com/sandbankdisperser/go-i2c-oled/i2c"
	"github.com/sandbankdisperser/go-i2c-oled/ssd1306"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type Stat struct {
	Ip   string
	Cpu  string
	Mem  string
	Disk string
}

func GlitchEffect(img image.Image, intensity float64) image.Image {
	bounds := img.Bounds()
	glitched := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		if rand.Float64() < intensity {
			switch rand.Intn(3) {
			case 0: // Décalage horizontal
				shift := rand.Intn(20) - 10
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					glitched.Set(x, y, img.At((x+shift)%bounds.Max.X, y))
				}
			case 1:
				continue
			case 2:
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, a := img.At(x, y).RGBA()
					glitched.Set(
						x,
						y,
						color.RGBA{
							uint8((r + rand.Uint32()) % 256),
							uint8((g + rand.Uint32()) % 256),
							uint8((b + rand.Uint32()) % 256),
							uint8(a),
						},
					)
				}
			}
		} else {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				glitched.Set(x, y, img.At(x, y))
			}
		}
	}

	return glitched
}

func CreateGlitchAnimation(img image.Image, frames int) []image.Image {
	var sequence []image.Image

	for i := 0; i < frames; i++ {
		intensity := rand.Float64()
		glitchedFrame := GlitchEffect(img, intensity)
		sequence = append(sequence, glitchedFrame)
	}

	return sequence
}

func executeCmd(command string, args ...string) string {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return ""
	}
	return strings.TrimSpace(out.String())
}

func GetStat() Stat {
	return Stat{
		Ip:   executeCmd("bash", "-c", "hostname -I | cut -d' ' -f1"),
		Cpu:  executeCmd("bash", "-c", "top -bn1 | grep load | awk '{printf \"CPU Load: %.2f\", $(NF-2)}'"),
		Mem:  executeCmd("bash", "-c", "free -m | awk 'NR==2{printf \"Mem: %.2f%%\", $3*100/$2 }'"),
		Disk: executeCmd("bash", "-c", "df -h | awk '$NF==\"/\"{printf \"Disk: %d/%dGB %s\", $3,$2,$5}'"),
	}
}

func main() {
	// Initialize the OLED with specific settings
	// oled, err := goi2coled.NewI2c(ssd1306.SSD1306_SWITCHCAPVCC, 64, 128, 0x3C, 1)
	i2c, err := i2c.NewI2c(0x3c, 1)
	if err != nil {
		log.Fatalf("Failed to open i2c device: %v", err)
	}

	oled := ssd1306.NewSSD1306_128_64(i2c, ssd1306.SSD1306_SWITCHCAPVCC)
	if err != nil {
		log.Fatalf("Failed to create a new display: %v", err)
	}

	err = oled.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize the screen: %v", err)
	}

	colWhite := color.RGBA{255, 255, 255, 255}
	// black := color.RGBA{0, 0, 0, 255}
	img := image.NewRGBA(image.Rect(0, 0, 128, 64))
	point := fixed.Point26_6{fixed.Int26_6(0 * 64), fixed.Int26_6(15 * 64)} // x = 0, y = 15

	// Configure the font drawer with the chosen font and color
	drawer := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{colWhite},
		Face: basicfont.Face7x13,
		Dot:  point,
	}

	// Clear the OLED image (making it all black)
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	// Draw the text "Hello" on the OLED image
	drawer.DrawString("Hello, World!")
	oled.DisplayOn()

	time.Sleep(10 * time.Second)
	oled.DisplayOff()

	// // Define a white color for text and drawings
	// colWhite := color.RGBA{255, 255, 255, 255}

	// // Set the starting point for drawing text
	// point := fixed.Point26_6{fixed.Int26_6(0 * 64), fixed.Int26_6(0 * 64)} // x = 0, y = 0 initially

	// img := image.NewRGBA(image.Rect(0, 0, 128, 64))

	// // Configure the font drawer with the chosen font and color
	// drawer := &font.Drawer{
	// 	Dst:  img,
	// 	Src:  &image.Uniform{colWhite},
	// 	Face: basicfont.Face7x13,
	// 	Dot:  point,
	// }

	// // Setting up channel for graceful shutdown
	// done := make(chan os.Signal, 1)
	// stopCh := make(chan bool, 1)

	// signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// // Load Raspberry Pi logo
	// logoFile, err := os.Open("./rpi1.png")
	// if err != nil {
	// 	fmt.Println("Error opening logo file:", err)
	// 	return
	// }
	// defer logoFile.Close()
	// logoImg, err := png.Decode(logoFile)
	// if err != nil {
	// 	fmt.Println("Error decoding logo image:", err)
	// 	return
	// }

	// go func() {
	// 	for {
	// 		select {
	// 		case <-stopCh:
	// 			return
	// 		default:
	// 			// Main loop for updating the display
	// 			// Set the entire OLED image to black
	// 			draw.Draw(img, img.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)

	// 			drawer.Dot.Y = fixed.Int26_6(0 * 64)
	// 			drawer.Dot.X = fixed.Int26_6(0 * 64)
	// 			dstRect := image.Rect(0, 0, logoImg.Bounds().Dx(), logoImg.Bounds().Dy()+15)
	// 			draw.Draw(img, dstRect, logoImg, image.Point{}, draw.Over)

	// 			drawer.Dot.Y = fixed.Int26_6(14 * 64)
	// 			drawer.Dot.X = fixed.Int26_6(
	// 				(logoImg.Bounds().Dx() + 10) * 64,
	// 			)

	// 			drawer.DrawString("Rpi 4")
	// 			drawer.Dot.Y += fixed.Int26_6(30 * 64)
	// 			drawer.Dot.X = fixed.Int26_6(
	// 				(logoImg.Bounds().Dx() + 10) * 64,
	// 			)
	// 			drawer.DrawString(executeCmd("hostname"))
	// 			oled.DrawImage(img)
	// 			oled.DisplayOn()

	// 			time.Sleep(time.Second * 10)

	// 			transistion := CreateGlitchAnimation(img, 10)
	// 			for _, t := range transistion {
	// 				drawer.Dot.Y = fixed.Int26_6(0 * 64)
	// 				drawer.Dot.X = fixed.Int26_6(0 * 64)
	// 				draw.Draw(img, dstRect, t, image.Point{}, draw.Over)
	// 				oled.DrawImage(img)
	// 				oled.DisplayOn()
	// 				time.Sleep(time.Microsecond * 2)
	// 			}

	// 			drawer.Dot.Y = fixed.Int26_6(10 * 64)
	// 			drawer.Dot.X = fixed.Int26_6(0 * 64)
	// 			draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	// 			stat := GetStat()
	// 			drawer.DrawString(fmt.Sprintf("Ip: %s", stat.Ip))
	// 			drawer.Dot.Y += fixed.Int26_6(16 * 64)
	// 			drawer.Dot.X = fixed.Int26_6(0 * 64)
	// 			drawer.DrawString(stat.Cpu)
	// 			drawer.Dot.Y += fixed.Int26_6(14 * 64)
	// 			drawer.Dot.X = fixed.Int26_6(0 * 64)
	// 			drawer.DrawString(stat.Mem)
	// 			drawer.Dot.Y += fixed.Int26_6(14 * 64)
	// 			drawer.Dot.X = fixed.Int26_6(0 * 64)
	// 			drawer.DrawString(stat.Disk)
	// 			oled.DrawImage(img)
	// 			oled.DisplayOn()
	// 			time.Sleep(time.Second * 10)
	// 		}
	// 	}
	// }()
	// <-done
	// stopCh <- true
	fmt.Printf("Stop programme")
}

func Oldmain() {
	// // Initialize the OLED display with the provided parameters
	// oled, err := goi2coled.NewI2c(ssd1306.SSD1306_SWITCHCAPVCC, 32, 128, 0x3C, 1)
	// if err != nil {
	// 	panic(err)
	// }

	// // Ensure the OLED is properly closed at the end of the program
	// defer oled.Close()

	// // Define a black color
	// black := color.RGBA{0, 0, 0, 255}

	// // Set the entire OLED image to black
	// draw.Draw(oled.Img, oled.Img.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)

	// // Define a white color
	// colWhite := color.RGBA{255, 255, 255, 255}

	// // Set the starting point for drawing text
	// point := fixed.Point26_6{fixed.Int26_6(0 * 64), fixed.Int26_6(15 * 64)} // x = 0, y = 15

	// // Configure the font drawer with the chosen font and color
	// drawer := &font.Drawer{
	// 	Dst:  oled.Img,
	// 	Src:  &image.Uniform{colWhite},
	// 	Face: basicfont.Face7x13,
	// 	Dot:  point,
	// }

	// // Clear the OLED image (making it all black)
	// draw.Draw(oled.Img, oled.Img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	// // Draw the text "Hello" on the OLED image
	// drawer.DrawString("Hello")

	// // Move the drawing point down by 10 pixels for the next line of text
	// drawer.Dot.Y += fixed.Int26_6(10 * 64)

	// // Set the drawing point's x coordinate back to 0 for alignment
	// drawer.Dot.X = fixed.Int26_6(0 * 64)

	// // Draw the text "From golang!" on the OLED image
	// drawer.DrawString("From golang!")

	// // Clear the OLED's buffer (if applicable to your library)
	// oled.Clear()

	// // Update the OLED's buffer with the current image data
	// oled.Draw()

	// // Display the buffered content on the OLED screen
	// err = oled.Display()
}
