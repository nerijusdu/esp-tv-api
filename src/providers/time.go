package providers

import (
	"image/color"
	"time"

	"github.com/fogleman/gg"
	"github.com/nerijusdu/esp-tv-api/src/constants"
	"github.com/nerijusdu/esp-tv-api/src/util"
)

type TimeProvider struct{}

func (p *TimeProvider) Init(config any) error {
	return nil
}

func (p *TimeProvider) GetView(cursor string) (ViewResponse, error) {
	response := ViewResponse{
		Cursor:     cursor,
		NextCursor: "",
		View: View{
			Data:         make([]byte, constants.DISPLAY_SIZE),
			RefreshAfter: 5000,
		},
	}

	err := drawTime(&response.View.Data)
	if err != nil {
		return response, err
	}

	return response, nil
}

func drawTime(result *[]byte) error {
	timeText := time.Now().Format("15:04")
	dateText := time.Now().Format("2006-01-02")
	dayOfWeek := time.Now().Weekday().String()
	dc := gg.NewContext(constants.DISPLAY_WIDTH, constants.DISPLAY_HEIGHT)
	dc.SetColor(color.White)
	dc.DrawStringAnchored(dateText, constants.DISPLAY_WIDTH/2, constants.DISPLAY_HEIGHT/2-20, 0.5, 0.5)
	dc.DrawStringAnchored(timeText, constants.DISPLAY_WIDTH/2, constants.DISPLAY_HEIGHT/2-3, 0.5, 0.5)
	dc.DrawStringAnchored(dayOfWeek, constants.DISPLAY_WIDTH/2, constants.DISPLAY_HEIGHT/2+15, 0.5, 0.5)
	dc.Stroke()

	util.ImageToBytes(dc, result)

	return nil
}
