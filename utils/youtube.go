package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_proxy_worker/db"
	"go_proxy_worker/models"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type YouTubeVideo_VideoDescriptionHeaderRenderer struct {
	Title struct {
		Runs []struct {
			Text string `json:"text"`
		} `json:"runs"`
	} `json:"title"`
	Channel struct {
		SimpleText string `json:"simpleText"`
	} `json:"channel"`
	Views struct {
		SimpleText string `json:"simpleText"`
	} `json:"views"`
	PublishDate struct {
		SimpleText string `json:"simpleText"`
	} `json:"publishDate"`
	Factoid []struct {
		FactoidRenderer struct {
			Value struct {
				SimpleText string `json:"simpleText"`
			} `json:"value"`
			Label struct {
				SimpleText string `json:"simpleText"`
			} `json:"label"`
		} `json:"factoidRenderer"`
	} `json:"factoid"`
	ChannelNavigationEndpoint struct {
		BrowseEndpoint struct {
			BrowseId string `json:"browseId"`
		} `json:"browseEndpoint"`
	} `json:"channelNavigationEndpoint"`
}

type YouTubeVideo_ExpandableVideoDescriptionBodyRenderer struct {
	AttributedDescriptionBodyText struct {
		Content     string `json:"content"`
		CommandRuns []struct {
			OnTap struct {
				InnertubeCommand struct {
					UrlEndpoint struct {
						Url string `json:"url"`
					} `json:"urlEndpoint"`
				} `json:"innertubeCommand"`
			} `json:"onTap"`
		} `json:"commandRuns"`
	} `json:"attributedDescriptionBodyText"`
}

type YouTubeVideo_StructuredDescriptionContentRenderer struct {
	Items []map[string]interface{} `json:"items"`
}

type YouTubeVideo_SegmentedLikeDislikeButtonViewModel struct {
	Like struct {
		Like struct {
			Toggle struct {
				Toggle struct {
					Default struct {
						Button struct {
							Text string `json:"accessibilityText"`
						} `json:"buttonViewModel"`
					} `json:"defaultButtonViewModel"`
				} `json:"toggleButtonViewModel"`
			} `json:"toggleButtonViewModel"`
		} `json:"likeButtonViewModel"`
	} `json:"likeButtonViewModel"`
}

type YouTubeVideo_VideoPrimaryInfoRenderer struct {
	VideoActions struct {
		MenuRenderer struct {
			TopLevelButtons []map[string]interface{} `json:"topLevelButtons"`
		} `json:"menuRenderer"`
	} `json:"videoActions"`
}

type YouTubeVideo_InitialData struct {
	EngagementPanels []struct {
		EngagementPanelSectionListRenderer struct {
			TargetId string                 `json:"targetId"`
			Content  map[string]interface{} `json:"content"`
		} `json:"engagementPanelSectionListRenderer"`
	} `json:"engagementPanels"`
	Contents struct {
		TwoColumnWatchNextResults struct {
			Results struct {
				Results struct {
					Contents []map[string]interface{} `json:"contents"`
				} `json:"results"`
			} `json:"results"`
		} `json:"twoColumnWatchNextResults"`
	} `json:"contents"`
}

type YouTubeVideo_InitialResponse struct {
	VideoDetails struct {
		LengthSeconds string `json:"lengthSeconds"`
	} `json:"videoDetails"`
	Captions struct {
		PlayerCaptionsTracklistRenderer struct {
			CaptionTracks []struct {
				BaseUrl string `json:"baseUrl"`
			} `json:"captionTracks"`
		} `json:"playerCaptionsTracklistRenderer"`
	} `json:"captions"`
}

type YouTubeVideo_Info struct {
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	Channel    string    `json:"channel"`
	ChannelURL string    `json:"channel_url"`
	Date       time.Time `json:"date"`
	Desc       string    `json:"desc"`
	Views      int       `json:"views"`
	Likes      int       `json:"likes"`
	Duration   int       `json:"duration"`
	Repo       string    `json:"repo"`
	Transcript string    `json:"caption"`
	Thumbnail  string    `json:"thumbnail"`
}

func YouTube_ParseDescriptionBody(content interface{}, data *YouTubeVideo_Info) error {
	var renderer YouTubeVideo_ExpandableVideoDescriptionBodyRenderer
	text, err := json.Marshal(content)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(text, &renderer); err != nil {
		return err
	}
	for _, el := range renderer.AttributedDescriptionBodyText.CommandRuns {
		if strings.Contains(el.OnTap.InnertubeCommand.UrlEndpoint.Url, "github.com") {
			if decodedUrl, err := url.QueryUnescape(el.OnTap.InnertubeCommand.UrlEndpoint.Url); err == nil {
				if github, err := ParseTextFromPattern(decodedUrl, `(https\:\/\/github.com\/[^\/]+\/[^\/\&\s\.]+)`); err == nil {
					data.Repo = github
				}
			}
		}
	}
	data.Desc = renderer.AttributedDescriptionBodyText.Content
	return nil
}

func YouTube_ParseDescriptionInfo(content interface{}, data *YouTubeVideo_Info) error {
	var renderer YouTubeVideo_VideoDescriptionHeaderRenderer
	text, err := json.Marshal(content)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(text, &renderer); err != nil {
		return err
	}
	return nil
}

func YouTube_ParseDescriptionHeader(content interface{}, data *YouTubeVideo_Info) error {
	var renderer YouTubeVideo_VideoDescriptionHeaderRenderer
	text, err := json.Marshal(content)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(text, &renderer); err != nil {
		return err
	}
	data.Channel = renderer.Channel.SimpleText
	data.ChannelURL = fmt.Sprintf("https://youtube.com/channel/%s", renderer.ChannelNavigationEndpoint.BrowseEndpoint.BrowseId)
	if views, err := ParseIntFromPattern(renderer.Views.SimpleText, `([\d\,]+)`); err == nil {
		data.Views = views
	}
	if time, err := time.Parse("Jan 2, 2006", renderer.PublishDate.SimpleText); err == nil {
		data.Date = time
	}
	if len(renderer.Title.Runs) > 0 {
		data.Name = renderer.Title.Runs[0].Text
	}
	return nil
}

func YouTube_ParseLikeButtonView(content interface{}, data *YouTubeVideo_Info) error {
	var view YouTubeVideo_SegmentedLikeDislikeButtonViewModel
	text, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	if err := json.Unmarshal(text, &view); err != nil {
		return err
	}
	like, err := ParseIntFromPattern(view.Like.Like.Toggle.Toggle.Default.Button.Text, `([\d\,]+)`)
	if err != nil {
		return err
	}
	data.Likes = like
	return nil
}

func YouTube_ParsePrimaryInfo(content interface{}, data *YouTubeVideo_Info) error {
	var renderer YouTubeVideo_VideoPrimaryInfoRenderer
	text, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	if err := json.Unmarshal(text, &renderer); err != nil {
		return err
	}
	for _, el := range renderer.VideoActions.MenuRenderer.TopLevelButtons {
		for key, value := range el {
			if key == "segmentedLikeDislikeButtonViewModel" {
				if err := YouTube_ParseLikeButtonView(value, data); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func YouTube_ParseStructuredDescription(content interface{}, data *YouTubeVideo_Info) error {
	var render YouTubeVideo_StructuredDescriptionContentRenderer
	text, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	if err := json.Unmarshal(text, &render); err != nil {
		return err
	}
	for _, item := range render.Items {
		for key, value := range item {
			if key == "videoDescriptionHeaderRenderer" {
				if err := YouTube_ParseDescriptionHeader(value, data); err != nil {
					return err
				}
			} else if key == "expandableVideoDescriptionBodyRenderer" {
				if err := YouTube_ParseDescriptionBody(value, data); err != nil {
					return err
				}
			} else if key == "videoDescriptionInfocardsSectionRenderer" {
				if err := YouTube_ParseDescriptionInfo(value, data); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func YouTube_ParseVideo(link string, data *YouTubeVideo_Info, simple bool) error {
	body, err := GetHtmlByProxy(link, false)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	dataText := ""
	respText := ""
	imgTag := doc.Find("link[as='image']").First()
	if imgTag.Length() == 0 {
		return errors.New("cannot find image")
	}
	data.Thumbnail = imgTag.AttrOr("href", "")
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "var ytInitialPlayerResponse =") {
			scriptText := strings.ReplaceAll(s.Text(), "var ", "^")
			respText, err = ParseTextFromPattern(scriptText, `ytInitialPlayerResponse\s+\=\s+([^\^]+)`)
			respText = strings.TrimSpace(respText)
		}
		if strings.Contains(s.Text(), "var ytInitialData =") {
			scriptText := strings.ReplaceAll(s.Text(), "var ", "^")
			dataText, err = ParseTextFromPattern(scriptText, `ytInitialData\s+\=\s+([^\^]+)`)
			dataText = strings.TrimSpace(dataText)
		}
	})
	if dataText == "" {
		return errors.New("cannot find initial data")
	}
	var info YouTubeVideo_InitialData
	if err := json.Unmarshal([]byte(dataText[0:len(dataText)-1]), &info); err != nil {
		return err
	}
	for _, el := range info.EngagementPanels {
		for key, content := range el.EngagementPanelSectionListRenderer.Content {
			if key == "structuredDescriptionContentRenderer" {
				if err := YouTube_ParseStructuredDescription(content, data); err != nil {
					return err
				}
			}
		}
	}
	for _, el := range info.Contents.TwoColumnWatchNextResults.Results.Results.Contents {
		for key, content := range el {
			if key == "videoPrimaryInfoRenderer" {
				if err := YouTube_ParsePrimaryInfo(content, data); err != nil {
					return err
				}
			}
		}
	}
	if simple {
		return nil
	}
	if respText == "" {
		return errors.New("cannot find initial responose")
	}
	var resp YouTubeVideo_InitialResponse
	if err := json.Unmarshal([]byte(respText[0:len(respText)-1]), &resp); err != nil {
		return err
	}
	if resp.VideoDetails.LengthSeconds == "" {
		return errors.New("cannot find duration")
	}
	if length, err := ParseIntFromPattern(resp.VideoDetails.LengthSeconds, `([\d\,]+)`); err == nil {
		data.Duration = length
	}
	if len(resp.Captions.PlayerCaptionsTracklistRenderer.CaptionTracks) > 0 {
		url := resp.Captions.PlayerCaptionsTracklistRenderer.CaptionTracks[0].BaseUrl
		if transcript, err := GetHtmlByProxy(url, false); err == nil {
			data.Transcript = transcript
		} else {
			return err
		}
	} else {
		return errors.New("cannot find transcript")
	}
	return nil
}

func YouTube_ParseYoutubeLinks(url string, page int) (youtubeLinks []string, nextUrl string, err error) {
	body, err := GetHtmlByProxy(url, false)
	if err != nil {
		return youtubeLinks, nextUrl, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return youtubeLinks, nextUrl, err
	}
	doc.Find("#search #rso").Find("a[jsname='UWckNb']").Each(func(i int, s *goquery.Selection) {
		link := s.AttrOr("href", "")
		if strings.Contains(link, "youtube.com/watch?") {
			youtubeLinks = append(youtubeLinks, link)
		}
	})
	nextLinkTag := doc.Find("#botstuff").Find("[role='navigation']").Find("a#pnnext").First()
	if nextLinkTag.Length() > 0 {
		nextUrl = fmt.Sprintf("https://www.google.com%s", nextLinkTag.AttrOr("href", ""))
	}
	return youtubeLinks, nextUrl, nil
}

func YouTube_FindRecord(url string, record *models.MarketingWebsiteVideo) bool {
	dbInst := db.GetDB()
	result := dbInst.Where("video_url = ?", url).First(record)
	return result.Error == nil
}

func YouTube_UpdateRecord(record *models.MarketingWebsiteVideo) error {
	dbInst := db.GetDB()
	result := dbInst.Save(&record)
	return result.Error
}

func YouTube_CreateRecord(record *models.MarketingWebsiteVideo) error {
	dbInst := db.GetDB()
	result := dbInst.Create(&record)
	return result.Error
}
