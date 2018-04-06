package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/archerwq/photo-viewer/dao"
)

type quryParams struct {
	lat       string
	lon       string
	radius    float64
	keywords  string
	startTime string
	endTime   string
	offset    int
	limit     int
}

func parseParams(r *http.Request) *quryParams {
	qp := &quryParams{}
	params := r.URL.Query()
	if params["lat"] != nil {
		qp.lat = params["lat"][0]
	}
	if params["lon"] != nil {
		qp.lon = params["lon"][0]
	}
	if params["r"] != nil {
		r, err := strconv.ParseFloat(params["r"][0], 64)
		if err != nil {
			r = 1.0
		}
		qp.radius = r
	} else {
		qp.radius = 1.0
	}
	if params["kw"] != nil {
		qp.keywords = params["kw"][0]
	}
	if params["start"] != nil {
		qp.startTime = params["start"][0]
	}
	if params["end"] != nil {
		qp.endTime = params["end"][0]
	}
	if params["offset"] != nil {
		o, err := strconv.Atoi(params["offset"][0])
		if err != nil {
			o = 0
		}
		qp.offset = o
	} else {
		qp.offset = 0
	}
	if params["limit"] != nil {
		l, err := strconv.Atoi(params["limit"][0])
		if err != nil {
			l = 500
		}
		qp.limit = l
	} else {
		qp.limit = 500
	}
	return qp
}

// GET /api/photos?lat=&lon=&r=&kw=&start=&end=&offset=&limit=
func QueryPhotos(w http.ResponseWriter, r *http.Request) {
	quryParams := parseParams(r)

	photos, err := dao.Manager.PhotoDao.QueryPhotos(quryParams.lat, quryParams.lon, quryParams.radius,
		quryParams.keywords, quryParams.startTime, quryParams.endTime, quryParams.offset, quryParams.limit)

	js, err := json.Marshal(map[string]interface{}{
		"success": true,
		"photos":  photos,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
