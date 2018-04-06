package dao

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/archerwq/photo-viewer/conf"

	"github.com/archerwq/photo-viewer/model"
	"github.com/olivere/elastic"
)

const queryAllTempl = `
{
	"query": {
		"bool": {
			"must": {
				"multi_match": {
					"query": "_KEYWORDS",
					"fields": [
						"tags",
						"story",
						"path"
					]
				}
			},
			"filter": [
				{
					"geo_distance": {
						"distance": "_RADIUS km",
						"location": {
							"lat": _LAT,
							"lon": _LON
						}
					}
				},
				{
					"range": {
						"time": {
							"gte": "_START",
							"lte": "_END"
						}
					}
				}
			]
		}
	},
	"sort": [
		{
			"time": "desc"
		}
	],
	"from": _FROM,
	"size": _SIZE
}`

const queryTemplNoLoc = `
{
	"query": {
		"bool": {
			"must": {
				"multi_match": {
					"query": "_KEYWORDS",
					"fields": [
						"tags",
						"story",
						"path"
					]
				}
			},
			"filter": [
				{
					"range": {
						"time": {
							"gte": "_START",
							"lte": "_END"
						}
					}
				}
			]
		}
	},
	"sort": [
		{
			"time": "desc"
		}
	],
	"from": _FROM,
	"size": _SIZE
}`

const queryTempNoKeywords = `
{
	"query": {
		"bool": {
			"must": {
				"match_all": {}
			},
			"filter": [
				{
					"geo_distance": {
						"distance": "_RADIUS km",
						"location": {
							"lat": _LAT,
							"lon": _LON
						}
					}
				},
				{
					"range": {
						"time": {
							"gte": "_START",
							"lte": "_END"
						}
					}
				}
			]
		}
	},
	"sort": [
		{
			"time": "desc"
		}
	],
	"from": _FROM,
	"size": _SIZE
}`

const queryTemplOnlyTime = `
{
	"query": {
		"bool": {
			"must": {
				"match_all": {}
			},
			"filter": [
				{
					"range": {
						"time": {
							"gte": "_START",
							"lte": "_END"
						}
					}
				}
			]
		}
	},
	"sort": [
		{
			"time": "desc"
		}
	],
	"from": _FROM,
	"size": _SIZE
}`

type PhotoDao struct {
	client *elastic.Client
}

func NewPhotoDao(config conf.ESConfig) (*PhotoDao, error) {
	c, err := elastic.NewClient(elastic.SetURL(config.Endpoint))
	if err != nil {
		return nil, err
	}
	return &PhotoDao{
		client: c,
	}, nil
}

func (p *PhotoDao) GetPhoto(sha1 string) (*model.Photo, error) {
	result, err := p.client.Get().
		Index("files").
		Type("photo").
		Id(sha1).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	if result.Found {
		var photo model.Photo
		err := json.Unmarshal(*result.Source, &photo)
		if err != nil {
			return nil, err
		}
		return &photo, nil
	}
	return nil, nil
}

// QueryPhotos return photos (order by time desc) qualify the given conditions.
func (p *PhotoDao) QueryPhotos(lat, lon string, radius float64, keywords, startTime, endTime string, from, size int) ([]model.Photo, error) {
	queryStr := genQueryStr(lat, lon, radius, keywords, startTime, endTime, from, size)

	searchResult, err := p.client.Search().
		Index("files").
		Type("photo").
		Source(queryStr).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	photos := make([]model.Photo, 0, size)
	var phttyp model.Photo
	for _, item := range searchResult.Each(reflect.TypeOf(phttyp)) {
		if p, ok := item.(model.Photo); ok {
			photos = append(photos, p)
		}
	}
	return photos, nil
}

func genQueryStr(lat, lon string, radius float64, keywords, startTime, endTime string, offset, limit int) string {
	var queryStr string

	withLoc := lat != "" && lon != ""
	withKeywords := keywords != ""
	if startTime == "" {
		startTime = "1983-12-14"
	}
	if endTime == "" {
		endTime = time.Now().Format("2006-01-02")
	}

	if limit == 0 {
		limit = 100
	}

	if withKeywords && withLoc {
		queryStr = strings.Replace(queryAllTempl, "_LAT", lat, 1)
		queryStr = strings.Replace(queryStr, "_LON", lon, 1)
		queryStr = strings.Replace(queryStr, "_RADIUS", strconv.FormatFloat(radius, 'f', -1, 64), 1)
		queryStr = strings.Replace(queryStr, "_KEYWORDS", keywords, 1)
		queryStr = strings.Replace(queryStr, "_START", startTime, 1)
		queryStr = strings.Replace(queryStr, "_END", endTime, 1)
		queryStr = strings.Replace(queryStr, "_FROM", strconv.Itoa(offset), 1)
		queryStr = strings.Replace(queryStr, "_SIZE", strconv.Itoa(limit), 1)
	}

	if withKeywords && !withLoc {
		queryStr = strings.Replace(queryTemplNoLoc, "_KEYWORDS", keywords, 1)
		queryStr = strings.Replace(queryStr, "_START", startTime, 1)
		queryStr = strings.Replace(queryStr, "_END", endTime, 1)
		queryStr = strings.Replace(queryStr, "_FROM", strconv.Itoa(offset), 1)
		queryStr = strings.Replace(queryStr, "_SIZE", strconv.Itoa(limit), 1)
	}

	if !withKeywords && withLoc {
		queryStr = strings.Replace(queryTempNoKeywords, "_LAT", lat, 1)
		queryStr = strings.Replace(queryStr, "_LON", lon, 1)
		queryStr = strings.Replace(queryStr, "_RADIUS", strconv.FormatFloat(radius, 'f', -1, 64), 1)
		queryStr = strings.Replace(queryStr, "_START", startTime, 1)
		queryStr = strings.Replace(queryStr, "_END", endTime, 1)
		queryStr = strings.Replace(queryStr, "_FROM", strconv.Itoa(offset), 1)
		queryStr = strings.Replace(queryStr, "_SIZE", strconv.Itoa(limit), 1)
	}

	if !withKeywords && !withLoc {
		queryStr = strings.Replace(queryTemplOnlyTime, "_START", startTime, 1)
		queryStr = strings.Replace(queryStr, "_END", endTime, 1)
		queryStr = strings.Replace(queryStr, "_FROM", strconv.Itoa(offset), 1)
		queryStr = strings.Replace(queryStr, "_SIZE", strconv.Itoa(limit), 1)
	}

	return queryStr
}

func (p *PhotoDao) ClearUp() {
	p.client.Stop()
}
