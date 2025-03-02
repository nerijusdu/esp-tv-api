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
	Feed string `json:"feed"`
}

func (p *BskyProvider) Init(config any) error {
	c, err := util.CastConfig[BskyConfig](config)
	if err != nil {
		return err
	}

	p.config = c

	return nil
}

func (p *BskyProvider) GetView(cursor string) (ViewResponse, error) {
	if p.currentPost == nil {
		res, err := getBskyPost(p.config.Feed)
		if err != nil {
			return ViewResponse{}, err
		}
		p.currentPost = res
	}

	dc := gg.NewContext(constants.DISPLAY_WIDTH, constants.DISPLAY_HEIGHT)
	dc.SetColor(color.White)

	author := p.currentPost.Author.DisplayName
	if author == "" {
		author = p.currentPost.Author.Handle
	}
	stats := fmt.Sprintf("%dL %dC", p.currentPost.LikeCount, p.currentPost.ReplyCount)
	postText := strings.ReplaceAll(p.currentPost.Record.Text, "\n", "")

	pages := dc.WordWrap(postText, constants.DISPLAY_WIDTH)
	paging, err := util.ParsePaging(cursor, len(pages))
	if err != nil {
		return ViewResponse{}, err
	}

	dc.DrawString(author, 0, 9)
	dc.DrawStringAnchored(stats, constants.DISPLAY_WIDTH, 9, 1, 0)
	dc.DrawLine(0, 10.5, constants.DISPLAY_WIDTH, 10.5)

	maxLines := 4
	for i, line := range pages[paging.IntCursor:] {
		if i >= maxLines {
			break
		}

		dc.DrawString(line, 0, 21+float64(i)*12)
	}

	dc.Stroke()

	if len(pages) <= maxLines {
		paging.NextCursor = ""
	}

	refreshAfter := 1000
	if paging.NextCursor == "" {
		p.currentPost = nil
		refreshAfter = 5000
	}
	if paging.IntCursor == 0 {
		refreshAfter = 3000
	}

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
