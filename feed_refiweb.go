package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/securityfirst/umbrella-api/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/gosexy/to"
)

type RefiWebFetcher struct {
	Country *models.Country
}

func (r *RefiWebFetcher) Fetch() ([]models.FeedItem, error) {
	body, err := makeRequest(fmt.Sprintf("https://api.reliefweb.int/v1/countries/%v", r.Country.ReliefWeb), "get", nil)
	if err != nil {
		return nil, err
	}
	var rwResp RWResponse
	if err = json.Unmarshal(body, &rwResp); err != nil {
		return nil, err
	}
	if len(rwResp.Data) < 1 {
		return nil, errors.New("No data received")
	}
	if rwResp.Data[0].Fields.DescriptionHTML == "" {
		return nil, nil
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rwResp.Data[0].Fields.DescriptionHTML))
	if err != nil {
		return nil, err
	}
	var items []models.FeedItem
	doc.Find("h3").First().Next().Children().Each(func(i int, t *goquery.Selection) {
		item, err := r.parseItem(t)
		if err != nil {
			log.Println("", err)
			return
		}
		items = append(items, *item)
	})
	return items, nil
}

func (r *RefiWebFetcher) parseItem(t *goquery.Selection) (*models.FeedItem, error) {
	href, ok := t.Contents().Attr("href")
	if !ok {
		return nil, errors.New("no href")
	}
	item := models.FeedItem{
		Title:     t.Contents().Text(),
		URL:       href,
		Country:   r.Country.Iso2,
		Source:    ReliefWeb,
		UpdatedAt: time.Now().Unix(),
	}
	segments := strings.Split(href, "/")
	if len(segments) == 0 || to.Int64(segments[len(segments)-1]) == 0 {
		return &item, nil
	}
	nodeUrl := fmt.Sprintf("http://api.rwlabs.org/v0/report/%v", segments[len(segments)-1])
	body, err := makeRequest(nodeUrl, "get", nil)
	if err != nil {
		return nil, err
	}
	var rwReport RWReport
	if err = json.Unmarshal(body, &rwReport); err != nil {
		return nil, err
	}
	if rwReport.Data[0].Fields.Headline.Summary != "" {
		item.Description = rwReport.Data[0].Fields.Headline.Summary
	} else {
		item.Description = rwReport.Data[0].Fields.BodyHTML
	}
	item.UpdatedAt = rwReport.Data[0].Fields.Date.Changed.Unix()
	return &item, nil
}

// type RWResponse struct {
// 	Version string `json:"version"`
// 	Status  int    `json:"status"`
// 	Time    int    `json:"time"`
// 	Data    struct {
// 		Type string `json:"type"`
// 		ID   int    `json:"id"`
// 		Item struct {
// 			ID              int    `json:"id"`
// 			Name            string `json:"name"`
// 			Description     string `json:"description"`
// 			Status          string `json:"status"`
// 			Iso3            string `json:"iso3"`
// 			Featured        bool   `json:"featured"`
// 			URL             string `json:"url"`
// 			DescriptionHTML string `json:"description-html"`
// 			Current         bool   `json:"current"`
// 			Location        struct {
// 				Lat  float64 `json:"lat"`
// 				Long float64 `json:"long"`
// 			} `json:"location"`
// 		} `json:"item"`
// 	} `json:"data"`
// }

type RWResponse struct {
	Href  string `json:"href"`
	Time  int    `json:"time"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Collection struct {
			Href string `json:"href"`
		} `json:"collection"`
	} `json:"links"`
	TotalCount int `json:"totalCount"`
	Count      int `json:"count"`
	Data       []struct {
		Fields struct {
			ID              int    `json:"id"`
			Name            string `json:"name"`
			Description     string `json:"description"`
			Status          string `json:"status"`
			Iso3            string `json:"iso3"`
			Featured        bool   `json:"featured"`
			VideoPlaylist   string `json:"video_playlist"`
			URL             string `json:"url"`
			URLAlias        string `json:"url_alias"`
			DescriptionHTML string `json:"description-html"`
			Current         bool   `json:"current"`
			Location        struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			} `json:"location"`
		} `json:"fields"`
		ID string `json:"id"`
	} `json:"data"`
}

type RWReport struct {
	Href  string `json:"href"`
	Time  int    `json:"time"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Collection struct {
			Href string `json:"href"`
		} `json:"collection"`
	} `json:"links"`
	TotalCount int `json:"totalCount"`
	Count      int `json:"count"`
	Data       []struct {
		Fields struct {
			ID       int    `json:"id"`
			Title    string `json:"title"`
			Status   string `json:"status"`
			Body     string `json:"body"`
			Headline struct {
				Title   string `json:"title"`
				Summary string `json:"summary"`
				Image   struct {
					ID        string `json:"id"`
					Width     string `json:"width"`
					Height    string `json:"height"`
					URL       string `json:"url"`
					Filename  string `json:"filename"`
					Mimetype  string `json:"mimetype"`
					Filesize  string `json:"filesize"`
					Copyright string `json:"copyright"`
					Caption   string `json:"caption"`
					URLLarge  string `json:"url-large"`
					URLSmall  string `json:"url-small"`
					URLThumb  string `json:"url-thumb"`
				} `json:"image"`
			} `json:"headline"`
			File []struct {
				ID          string `json:"id"`
				Description string `json:"description"`
				URL         string `json:"url"`
				Filename    string `json:"filename"`
				Mimetype    string `json:"mimetype"`
				Filesize    string `json:"filesize"`
				Preview     struct {
					URL      string `json:"url"`
					URLLarge string `json:"url-large"`
					URLSmall string `json:"url-small"`
					URLThumb string `json:"url-thumb"`
				} `json:"preview"`
			} `json:"file"`
			PrimaryCountry struct {
				Href     string `json:"href"`
				ID       int    `json:"id"`
				Name     string `json:"name"`
				Iso3     string `json:"iso3"`
				Location struct {
					Lat float64 `json:"lat"`
					Lon float64 `json:"lon"`
				} `json:"location"`
			} `json:"primary_country"`
			Country []struct {
				Href     string `json:"href"`
				ID       int    `json:"id"`
				Name     string `json:"name"`
				Iso3     string `json:"iso3"`
				Location struct {
					Lat float64 `json:"lat"`
					Lon float64 `json:"lon"`
				} `json:"location"`
				Primary bool `json:"primary"`
			} `json:"country"`
			Source []struct {
				Href        string `json:"href"`
				ID          int    `json:"id"`
				Name        string `json:"name"`
				Shortname   string `json:"shortname"`
				Longname    string `json:"longname,omitempty"`
				SpanishName string `json:"spanish_name,omitempty"`
				Homepage    string `json:"homepage"`
				Type        struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"type"`
			} `json:"source"`
			Language []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Code string `json:"code"`
			} `json:"language"`
			Theme []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"theme"`
			Format []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"format"`
			OchaProduct []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"ocha_product"`
			Disaster []struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				Glide string `json:"glide"`
				Type  []struct {
					ID      int    `json:"id"`
					Name    string `json:"name"`
					Code    string `json:"code"`
					Primary bool   `json:"primary"`
				} `json:"type"`
			} `json:"disaster"`
			DisasterType []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Code string `json:"code"`
			} `json:"disaster_type"`
			URL      string `json:"url"`
			URLAlias string `json:"url_alias"`
			BodyHTML string `json:"body-html"`
			Date     struct {
				Original time.Time `json:"original"`
				Changed  time.Time `json:"changed"`
				Created  time.Time `json:"created"`
			} `json:"date"`
		} `json:"fields"`
		ID string `json:"id"`
	} `json:"data"`
}
