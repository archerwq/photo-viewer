package dao

import (
	"testing"

	"github.com/archerwq/photo-viewer/conf"
)

func TestAddTag(t *testing.T) {
	ugiDao, err := NewUGIDao(conf.DBConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "qwang",
		Password: "qwer1234",
		DBName:   "photo",
	})
	if err != nil {
		t.Errorf("failed to create db client: %v", err)
		return
	}

	err = ugiDao.AddTag("09e7f623b8b7f441d128c920c140a3931e0fb82c", "life")
	if err != nil {
		t.Errorf("failed to add a tag: %v", err)
		return
	}
}
