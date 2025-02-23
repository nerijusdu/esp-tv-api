package providers

import (
	"fmt"

	"github.com/nerijusdu/esp-tv-api/src/constants"
)

type PosthogProvider struct{}

var pages = []string{
	"0",
	"1",
}

func (p *PosthogProvider) GetView(cursor string) ViewResponse {
	view := View{
		Data:         make([]byte, constants.DISPLAY_SIZE),
		RefreshAfter: 5000,
	}
	nextCursor := ""

	for i := 0; i < constants.DISPLAY_SIZE; i++ {
		if cursor == "" || cursor == "0" {
			view.Data[i] = []byte(fmt.Sprint(0))[0]
			nextCursor = "1"
		} else {
			view.Data[i] = []byte(fmt.Sprint(1))[0]
		}
	}

	return ViewResponse{
		View:       view,
		Cursor:     cursor,
		NextCursor: nextCursor,
	}
}
