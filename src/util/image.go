package util

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/nerijusdu/esp-tv-api/src/constants"
)

func ImageToBytes(dc *gg.Context, result *[]byte) {
	img := dc.Image()
	for i := 0; i < constants.DISPLAY_HEIGHT; i++ {
		for j := 0; j < constants.DISPLAY_WIDTH; j++ {
			r, _, _, _ := img.At(j, i).RGBA()
			if r > 0 {
				(*result)[i*constants.DISPLAY_WIDTH+j] = []byte(fmt.Sprint(1))[0]
			} else {
				(*result)[i*constants.DISPLAY_WIDTH+j] = []byte(fmt.Sprint(0))[0]
			}
		}
	}
}
