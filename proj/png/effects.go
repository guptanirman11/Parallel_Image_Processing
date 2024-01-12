// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import (
	"image/color"
)

// Grayscale applies a grayscale filtering effect to the image
func (img *Image) Grayscale(yMin int, yMax int) {

	// Bounds returns defines the dimensions of the image. Always
	// use the bounds Min and Max fields to get out the width
	// and height for the image
	bounds := img.Out.Bounds()
	for y := yMin; y < yMax; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			//Returns the pixel (i.e., RGBA) value at a (x,y) position
			// Note: These get returned as int32 so based on the math you'll
			// be performing you'll need to do a conversion to float64(..)
			r, g, b, a := img.In.At(x, y).RGBA()

			//Note: The values for r,g,b,a for this assignment will range between [0, 65535].
			//For certain computations (i.e., convolution) the values might fall outside this
			// range so you need to clamp them between those values.
			greyC := clamp(float64(r+g+b) / 3)

			//Note: The values need to be stored back as uint16 (I know weird..but there's valid reasons
			// for this that I won't get into right now).
			img.Out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
		}
	}
}

// Sharpen applies a sharpening filtering effect to the image
func (img *Image) Sharpen(yMin int, yMax int) {
	sharpenKernel := []float64{0, -1, 0, -1, 5, -1, 0, -1, 0}
	img.ApplyFilter(sharpenKernel, yMin, yMax)
}

// EdgeDetect applies an edge detection filtering effect to the image
func (img *Image) EdgeDetect(yMin int, yMax int) {
	edgeKernel := []float64{-1, -1, -1, -1, 8, -1, -1, -1, -1}
	img.ApplyFilter(edgeKernel, yMin, yMax)
}

// Blur applies a blur filtering effect to the image
func (img *Image) Blur(yMin int, yMax int) {
	blurKernel := []float64{1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0}
	img.ApplyFilter(blurKernel, yMin, yMax)
}

func (img *Image) ApplyFilter(kernel []float64, yMin int, yMax int) {
	bounds := img.Out.Bounds()
	xMin := bounds.Min.X
	xMax := bounds.Max.X

	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			var currentA, currentB, currentG, currentR uint32
			var rC, gC, bC float64

			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					row := x + i
					col := y + j

					if row >= xMin && row < xMax && col >= bounds.Min.Y && col < bounds.Max.Y {
						currentR, currentG, currentB, currentA = img.In.At(row, col).RGBA()
						index := (i+1)*3 + (j + 1)
						rC += float64(currentR) * kernel[index]
						gC += float64(currentG) * kernel[index]
						bC += float64(currentB) * kernel[index]
					}
				}
			}

			img.Out.Set(x, y, color.RGBA64{clamp(rC), clamp(gC), clamp(bC), uint16(currentA)})
		}
	}
}
