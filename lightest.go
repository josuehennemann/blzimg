package main

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
)

type LightestOperation struct{}

func (c LightestOperation) lightest(color1, color2 color.RGBA) color.Color {
	if c.luminance(color1) > c.luminance(color2) {
		return color1
	} else {
		return color2
	}
}

func (c LightestOperation) Result(images []ImageContainer) (image.Image, error) {
	if len(images) == 0 {
		return nil, nil
	}

	firstImage := images[0].getImage()
	bounds := firstImage.Bounds()
	lightest := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(lightest, bounds, firstImage, bounds.Min, draw.Src)

	for _, currentImageContainer := range images[1:] {
		currentImage := currentImageContainer.getImage()
		if currentImage.Bounds() != bounds {
			return nil, errors.New("The images have different size!")
		}

		imageToCompare := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
		draw.Draw(imageToCompare, bounds, currentImage, bounds.Min, draw.Src)
		c.getLightestImageBetweenTwo(lightest, imageToCompare)
	}

	return lightest, nil
}

func (c LightestOperation) getLightestImageBetweenTwo(current, other *image.RGBA) {
	for i := current.Bounds().Min.X; i < current.Bounds().Max.X; i++ {
		for j := current.Bounds().Min.Y; j < current.Bounds().Max.Y; j++ {
			currentLightestImagePixel := current.At(i, j).(color.RGBA)
			otherImagePixel := other.At(i, j).(color.RGBA)

			lightestColor := c.lightest(currentLightestImagePixel, otherImagePixel)
			if currentLightestImagePixel != lightestColor {
				current.Set(i, j, lightestColor)
			}
		}
	}
}

// http://stackoverflow.com/questions/596216/formula-to-determine-brightness-of-rgb-color
func (c LightestOperation) luminance(someColor color.Color) uint32 {
	r, g, b, _ := someColor.RGBA()
	return uint32(0.2126*float32(r) + 0.7152*float32(g) + 0.0722*float32(b))
}

func (c LightestOperation) average(someColor color.Color) uint32 {
	r, g, b, _ := someColor.RGBA()
	return (r + g + b) / 3
}
