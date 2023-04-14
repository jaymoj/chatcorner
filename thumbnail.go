package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run thumbnail.go <image file>")
		return
	}

	filename := os.Args[1]
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	// Read the header to determine the image format without moving the file pointer
	header := make([]byte, 8)
	if _, err = io.ReadFull(f, header); err != nil {
		fmt.Println("Error reading file header:", err)
		return
	}
	f.Seek(0, 0)

	// Decode the image using the appropriate format based on the header
	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	thumb := image.NewRGBA(image.Rect(0, 0, 100, 100))
	if err := resize(thumb, img); err != nil {
		fmt.Println("Error resizing image:", err)
		return
	}

	if err := saveThumb(thumb, "thumb.png"); err != nil {
		fmt.Println("Error saving thumbnail:", err)
		return
	}

	fmt.Println("Thumbnail generated successfully!")
}

func resize(dst *image.RGBA, src image.Image) error {
	srcBounds := src.Bounds()
	dstBounds := dst.Bounds()
	sw, sh := float64(srcBounds.Dx()), float64(srcBounds.Dy())
	dw, dh := float64(dstBounds.Dx()), float64(dstBounds.Dy())
	scaleX, scaleY := sw/dw, sh/dh

	for x := 0; x < dstBounds.Dx(); x++ {
		for y := 0; y < dstBounds.Dy(); y++ {
			px := int(float64(x)*scaleX + 0.5)
			py := int(float64(y)*scaleY + 0.5)
			dst.Set(x, y, src.At(px, py))
		}
	}

	return nil
}

func saveThumb(img *image.RGBA, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
