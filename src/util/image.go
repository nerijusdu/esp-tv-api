package util

import (
	"fmt"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/nerijusdu/esp-tv-api/src/constants"
)

func GraphicToBytes(dc *gg.Context) *[]byte {
	img := dc.Image()
	result := make([]byte, constants.DISPLAY_SIZE)
	for i := range constants.DISPLAY_HEIGHT {
		for j := range constants.DISPLAY_WIDTH {
			r, g, b, _ := img.At(j, i).RGBA()
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			pixel := color.Gray{uint8(lum / 256)}
			isOn := pixel.Y >= 128
			if isOn {
				result[i*constants.DISPLAY_WIDTH+j] = fmt.Append(nil, "1")[0]
			} else {
				result[i*constants.DISPLAY_WIDTH+j] = fmt.Append(nil, "0")[0]
			}
		}
	}

	return &result
}

type Image interface {
	At(x int, y int) color.Color
}

func AppendBWImageToBytes(img Image, starting *[]byte) *[]byte {
	temp := make([]byte, constants.DISPLAY_SIZE)
	for i := range constants.DISPLAY_HEIGHT {
		for j := range constants.DISPLAY_WIDTH {
			r, g, b, _ := img.At(j, i).RGBA()
			isOn := r+g+b >= 128
			if isOn {
				temp[i*constants.DISPLAY_WIDTH+j] = fmt.Append(nil, "1")[0]
			} else {
				temp[i*constants.DISPLAY_WIDTH+j] = fmt.Append(nil, "0")[0]
			}
		}
	}

	full := append(*starting, temp...)
	return &full
}
