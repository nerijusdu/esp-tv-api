package providers

import (
	"fmt"
	"image"
	"net/http"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/nerijusdu/esp-tv-api/src/constants"
	"github.com/nerijusdu/esp-tv-api/src/util"
	"golang.org/x/image/draw"
)

type ImageProvider struct {
	config ImageConfig
}

type ImageConfig struct {
	Urls []string `json:"urls"`
}

func (p *ImageProvider) Init(config any) error {
	c, err := util.CastConfig[ImageConfig](config)
	p.config = c

	return err
}

func (p *ImageProvider) GetView(cursor string) (ViewResponse, error) {
	response := ViewResponse{
		Cursor:     cursor,
		NextCursor: "",
		View: View{
			Data:         make([]byte, constants.DISPLAY_SIZE),
			RefreshAfter: 5000,
		},
	}

	if cursor == "" {
		cursor = "0"
	}

	intCursor, err := strconv.Atoi(cursor)
	if err != nil {
		return response, err
	}

	err = drawImage(p.config.Urls[intCursor], &response.View.Data)
	if err != nil {
		return response, err
	}

	intCursor++
	response.NextCursor = fmt.Sprint(intCursor)
	if intCursor >= len(p.config.Urls) {
		response.NextCursor = ""
	}

	return response, nil
}

func drawImage(url string, result *[]byte) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	if err != nil {
		return err
	}

	dc := gg.NewContext(constants.DISPLAY_WIDTH, constants.DISPLAY_HEIGHT)

	//scale to fit
	dst := image.NewRGBA(image.Rect(0, 0, constants.DISPLAY_WIDTH, constants.DISPLAY_HEIGHT))
	draw.NearestNeighbor.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	dc.DrawImage(dst, 0, 0)
	dc.Stroke()

	util.ImageToBytes(dc, result)

	return nil
}
