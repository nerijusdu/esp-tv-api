package providers

import (
	"encoding/json"
	"fmt"
	"image/color"
	"strings"

	"net/http"
	"net/url"

	"io"

	"github.com/fogleman/gg"
	"github.com/nerijusdu/esp-tv-api/src/constants"
	"github.com/nerijusdu/esp-tv-api/src/util"
)

type BskyProvider struct {
	config      BskyConfig
	currentPost *BskyPost
}

type BskyConfig struct {
	Feed         string `json:"feed"`
	RenderImages bool   `json:"renderImages"`
}

func (p *BskyProvider) Init(config any) error {
	c, err := util.CastConfig[BskyConfig](config)
	if err != nil {
		return err
	}

	p.config = c

	return nil
}

const MAX_LINES = 4

func (p *BskyProvider) GetView(cursor string) (ViewResponse, error) {
	if p.currentPost == nil {
		res, err := getBskyPost(p.config.Feed)
		if err != nil {
			return ViewResponse{}, err
		}
		p.currentPost = res
	}

	if cursor == "image" {
		data, err := drawImage(p.currentPost.Embed.Images[0].Thumb)
		if err != nil {
			return ViewResponse{}, err
		}

		p.currentPost = nil
		return ViewResponse{
			Cursor:     cursor,
			NextCursor: "",
			View: View{
				RefreshAfter: 5000,
				Data:         *data,
			},
		}, nil
	}

	author, stats, postText := p.getPostText()

	dc := gg.NewContext(constants.DISPLAY_WIDTH, constants.DISPLAY_HEIGHT)
	pages := dc.WordWrap(postText, constants.DISPLAY_WIDTH)
	paging, err := util.ParsePaging(cursor, len(pages))
	if err != nil {
		return ViewResponse{}, err
	}

	refreshAfter := 1000
	if len(pages) <= MAX_LINES {
		paging.NextCursor = ""
	}
	if paging.NextCursor == "" {
		if p.config.RenderImages && len(p.currentPost.Embed.Images) > 0 {
			paging.NextCursor = "image"
		} else {
			p.currentPost = nil
		}
		refreshAfter = 3000
	} else if paging.IntCursor == 0 {
		refreshAfter = 3000
	}

	renderPostText(dc, author, stats, pages[paging.IntCursor:])
	bytes := util.GraphicToBytes(dc)

	return ViewResponse{
		Cursor:     paging.Cursor,
		NextCursor: paging.NextCursor,
		View: View{
			RefreshAfter: refreshAfter,
			Data:         *bytes,
		},
	}, nil
}

func (p *BskyProvider) getPostText() (author string, stats string, postText string) {
	author = p.currentPost.Author.DisplayName
	if author == "" {
		author = p.currentPost.Author.Handle
	}
	stats = fmt.Sprintf("%dL %dC", p.currentPost.LikeCount, p.currentPost.ReplyCount)
	postText = strings.ReplaceAll(p.currentPost.Record.Text, "\n", "")
	return
}

func renderPostText(dc *gg.Context, author, stats string, pages []string) {
	dc.SetColor(color.White)
	dc.DrawString(author, 0, 9)
	dc.DrawStringAnchored(stats, constants.DISPLAY_WIDTH, 9, 1, 0)
	dc.DrawLine(0, 10.5, constants.DISPLAY_WIDTH, 10.5)
	for i, line := range pages {
		if i >= MAX_LINES {
			break
		}
		dc.DrawString(line, 0, 21+float64(i)*12)
	}
	dc.Stroke()

}

func getBskyPost(feed string) (*BskyPost, error) {
	feedUrl := fmt.Sprintf("https://public.api.bsky.app/xrpc/app.bsky.feed.getFeed?limit=1&feed=%s", url.QueryEscape(feed))
	req, err := http.NewRequest("GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseText := ""
		if resp.Body != nil {
			bodyBytes, _ := io.ReadAll(resp.Body)
			responseText = string(bodyBytes)
		}
		return nil, fmt.Errorf("unexpected status code: %d. %s", resp.StatusCode, responseText)
	}

	var result BskyFeedResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result.Feed[0].Post, nil
}
