package dao

import (
	"github.com/archerwq/photo-viewer/conf"
)

type DaoManager struct {
	PhotoDao *PhotoDao
	UGIDao   *UGIDao
}

var Manager DaoManager

func InitManager(config *conf.Config) error {
	p, err := NewPhotoDao(config.ES)
	if err != nil {
		return err
	}
	u, err := NewUGIDao(config.DB)
	if err != nil {
		return err
	}

	Manager = DaoManager{
		PhotoDao: p,
		UGIDao:   u,
	}
	return nil
}

func CleanManager() {
	Manager.PhotoDao.ClearUp()
	Manager.UGIDao.ClearUp()
}
