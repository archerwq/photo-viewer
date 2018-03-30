package dao

import (
	"testing"
)

func TestQueryPhotos(t *testing.T) {
	photoDao, err := New()
	if err != nil {
		t.Errorf("failed to create es client")
		return
	}

	photos, err := photoDao.QueryPhotos("22.275754", "113.578935", 1, "Users", "2018-02-20", "2018-03-29", 0, 10)
	if err != nil {
		t.Errorf("failed to query photos")
		return
	}
	if len(photos) < 1 {
		t.Errorf("expected: >0 actually: %v", len(photos))
	}
}
