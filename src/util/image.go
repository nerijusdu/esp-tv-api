package util

import (
	"fmt"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/nerijusdu/esp-tv-api/src/constants"
)

func ImageToBytes(dc *gg.Context, result *[]byte) {
	img := dc.Image()
	for i := 0; i < constants.DISPLAY_HEIGHT; i++ {
		for j := 0; j < constants.DISPLAY_WIDTH; j++ {
			r, g, b, _ := img.At(j, i).RGBA()
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			pixel := color.Gray{uint8(lum / 256)}
			isOn := pixel.Y >= 128
			if isOn {
				(*result)[i*constants.DISPLAY_WIDTH+j] = []byte(fmt.Sprint(1))[0]
			} else {
				(*result)[i*constants.DISPLAY_WIDTH+j] = []byte(fmt.Sprint(0))[0]
			}
		}
	}
}
