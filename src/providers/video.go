package providers

import (
	"image"
	"image/draw"
	"image/gif"
	"os"
	"strconv"

	"github.com/nerijusdu/esp-tv-api/src/constants"
	"github.com/nerijusdu/esp-tv-api/src/util"
)

type VideoProvider struct {
	currentVideo *gif.GIF
}

func (p *VideoProvider) Init(config any) error {
	return nil
}

func (p *VideoProvider) GetName() string {
	return "video"
}

func (p *VideoProvider) GetView(cursor string) (ViewResponse, error) {
	// use this ffmpeg command to generate gif from video
	// ffmpeg -i input.mp4 -f lavfi -i color=gray:s=1280x720 -f lavfi -i color=black:s=1280x720 -f lavfi -i color=white:s=1280x720 -filter_complex threshold,scale=128:64,fps=5 -y output.gif

	if p.currentVideo == nil {
		file, err := os.Open("output.gif")
		if err != nil {
			return ViewResponse{}, err
		}
		defer file.Close()

		gifImage, err := gif.DecodeAll(file)
		if err != nil {
			return ViewResponse{}, err
		}
		p.currentVideo = gifImage
	}

	result := ViewResponse{
		Cursor: cursor,
		View: View{
			Data:         make([]byte, 0),
			RefreshAfter: 200,
		},
	}

	paging, err := util.ParsePaging(cursor, len(p.currentVideo.Image))
	if err != nil {
		return result, err
	}
	nextCursor := paging.IntCursor + 5

	overpaintImage := image.NewRGBA(image.Rect(0, 0, constants.DISPLAY_WIDTH, constants.DISPLAY_HEIGHT))
	draw.Draw(overpaintImage, overpaintImage.Bounds(), p.currentVideo.Image[0], image.Point{}, draw.Src)
	for i, frame := range p.currentVideo.Image {
		draw.Draw(overpaintImage, overpaintImage.Bounds(), frame, image.Point{}, draw.Over)
		if i <= nextCursor && i >= paging.IntCursor {
			result.View.Data = *util.AppendBWImageToBytes(overpaintImage, &result.View.Data)
		}
	}

	if nextCursor >= len(p.currentVideo.Image) {
		nextCursor = 0
	}
	result.NextCursor = strconv.Itoa(nextCursor)

	return result, nil
}
