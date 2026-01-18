package main

import (
	"image"
	"image/color"
	"minecraftRemade/renderer" // Replace 'go-gl-image' with your module name in go.mod
)

func main() {
	// 1. Generate an image variable (example: a simple 500x500 blue square)
	myImage := generateImage(500, 500)

	// 2. Pass it to the renderer
	renderer.Show(myImage)
}

func generateImage(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	blue := color.RGBA{0, 100, 255, 255}

	// Fill the image with color
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, blue)
		}
	}
	return img
}
