package main

import (
	"image"
	"image/color"
	"minecraftRemade/renderer" // Ensure this matches your go.mod name
)

var offset int

func main() {
	// The third argument is the callback function that provides each frame
	renderer.Start(500, 500, updateFrame)
}

func updateFrame() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 500, 500))
	offset++

	// Create a simple moving gradient for testing 60fps
	for x := 0; x < 500; x++ {
		for y := 0; y < 500; y++ {
			c := color.RGBA{
				R: uint8((x + offset) % 255),
				G: uint8((y + offset) % 255),
				B: 150,
				A: 255,
			}
			img.Set(x, y, c)
		}
	}
	return img
}
