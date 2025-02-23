package providers

import (
	"encoding/json"
	"fmt"
	"image/color"
	"net/http"
	"os"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/nerijusdu/esp-tv-api/src/constants"
	"github.com/nerijusdu/esp-tv-api/src/util"
)

type PosthogProvider struct{}

type PosthogSite struct {
	title     string
	projectId string
	insightId string
}

var pages = map[string]PosthogSite{
	"0": {
		title:     "lingvistas.lt",
		projectId: "43890",
		insightId: "687096",
	},
	"1": {
		title:     "manoakimirka.lt",
		projectId: "41142",
		insightId: "645533",
	},
}

type PosthogInsightResponse struct {
	Result []struct {
		Data   []int    `json:"data"`
		Labels []string `json:"labels"`
		Days   []string `json:"days"`
	} `json:"result"`
}

func (p *PosthogProvider) getSiteStats(projectId string, insightId string) (PosthogInsightResponse, error) {
	token := os.Getenv("POSTHOG_API_KEY")
	url := fmt.Sprintf("https://eu.i.posthog.com/api/projects/%s/insights/%s", projectId, insightId)
	result := PosthogInsightResponse{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return result, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (p *PosthogProvider) renderData(data PosthogInsightResponse, site PosthogSite, result *[]byte) error {
	if len(data.Result) == 0 {
		return nil
	}

	const w = constants.DISPLAY_WIDTH
	const h = constants.DISPLAY_HEIGHT
	values := data.Result[0].Data
	maxValue := 0
	for _, d := range values {
		if d > maxValue {
			maxValue = d
		}
	}

	maxValueTextSize := 6
	if maxValue > 99 {
		maxValueTextSize = 20
	} else if maxValue > 9 {
		maxValueTextSize = 13
	}

	dc := gg.NewContext(w, h)
	dc.SetColor(color.White)
	dc.DrawString(site.title, 0, h-1)
	dc.DrawString(fmt.Sprint(maxValue), w-float64(maxValueTextSize), 9)

	dc.DrawLine(0, h-11.5, w, h-11.5)

	lastIndex := 0
	hStep := w / len(values)
	vStep := h / maxValue
	for i := 1; i < len(values); i++ {
		dc.DrawLine(
			float64(lastIndex*hStep)+0.5,
			float64(h-13-(values[lastIndex]*vStep))+0.5,
			float64(i*hStep)+0.5,
			float64(h-13-(values[i]*vStep))+0.5,
		)
		lastIndex = i
	}

	dc.Stroke()

	util.ImageToBytes(dc, result)

	return nil
}

func (p *PosthogProvider) GetView(cursor string) (ViewResponse, error) {
	view := View{
		Data:         make([]byte, constants.DISPLAY_SIZE),
		RefreshAfter: 5000,
	}
	if cursor == "" {
		cursor = "0"
	}

	intCursor, err := strconv.Atoi(cursor)
	if err != nil {
		return ViewResponse{}, fmt.Errorf("invalid cursor: %s", cursor)
	}
	nextCursor := intCursor + 1

	result := ViewResponse{
		Cursor:     cursor,
		NextCursor: fmt.Sprint(nextCursor),
		View:       view,
	}
	if nextCursor >= len(pages) {
		result.NextCursor = ""
	}

	site := pages[cursor]
	data, err := p.getSiteStats(site.projectId, site.insightId)
	if err != nil {
		return result, err
	}

	err = p.renderData(data, site, &view.Data)
	if err != nil {
		return result, err
	}

	return result, nil
}
