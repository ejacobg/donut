package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"os"
	"time"

	"github.com/quincinia/donut/frame"
	"golang.org/x/term"
)

var (
	GIF    = flag.String("gif", "", "Generate a GIF file with the given name")
	steps  = flag.Int("steps", 200, "Number of frames to generate")
	delay  = flag.Int("delay", 50, "Delay between frames, in milliseconds")
	width  = flag.Int("width", 100, "Width of the GIF file, if needed")
	height = flag.Int("height", 100, "Height of the GIF file, if needed")
)

func main() {
	flag.Parse()
	genGIF := GIF != nil && *GIF != ""

	var f frame.Frame
	var min int
	rect := image.Rect(0, 0, *width, *height)
	palette := frame.Gray16Palette()
	g := &gif.GIF{
		Config: image.Config{ColorModel: palette, Width: *width, Height: *height},
		BackgroundIndex: 0,
	}
	if genGIF {
		f = &frame.GIF{
			Image: image.NewPaletted(rect, palette),
			SF:    frame.SetFunc[color.Color](frame.Gray16SF),
		}
	} else {
		width, height, err := term.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatalln(err)
		}

		if width < height {
			min = width
		} else {
			// If we're constrained by height, we subtract 1 because fmt.Println adds an extra line.
			// If we don't do this, we get artifacts when multiple renders occur.
			min = height - 1
		}

		f = frame.NewASCII(min, min, frame.ASCIISF)
	}

	torus := frame.DefaultTorus()
	scene := frame.Scene{
		Width:  min,
		Height: min,
		K2:     5.0,
		LX:     0.0,
		LY:     1.0,
		LZ:     -1.0,
	}
	if genGIF {
		scene.Width, scene.Height = *width, *height
	}
	scene.CalculateK1(torus.R1, torus.R2)

	// How much the angles should change between each step.
	// Move this inside the Torus struct?
	startA, startB := 15.0, 25.0
	stepA, stepB := 0.07, 0.03

	for {
		var s int
		for s, torus.A, torus.B = 0, startA, startB; s < *steps; s, torus.A, torus.B = s+1, torus.A+stepA, torus.B+stepB {
			f.Reset() // Should the Render function clear the frame instead of doing it here?
			frame.Render(f, torus, scene)
			switch v := f.(type) {
			case *frame.GIF:
				g.Image = append(g.Image, v.Image)
				g.Delay = append(g.Delay, *delay/10)
				g.Disposal = append(g.Disposal, gif.DisposalBackground)
				v.Image = image.NewPaletted(rect, palette)
			case *frame.ASCII:
				v.Print(os.Stdout)
				fmt.Printf("\033[%dF", min)
			}
			if !genGIF {
				time.Sleep(time.Duration(*delay) * time.Millisecond)
			}
		}
		if genGIF {
			break
		}
	}
	file, _ := os.OpenFile(*GIF, os.O_WRONLY|os.O_CREATE, 0600)
	defer file.Close()
	gif.EncodeAll(file, g)
}
