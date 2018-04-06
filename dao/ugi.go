package dao

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/archerwq/photo-viewer/conf"
	"github.com/archerwq/photo-viewer/model"

	_ "github.com/go-sql-driver/mysql"

	. "github.com/archerwq/photo-viewer/pvlog"
)

type UGIDao struct {
	db *sql.DB
}

func NewUGIDao(config conf.DBConfig) (*UGIDao, error) {
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	dataSourceName := config.User + ":" +
		config.Password + "@tcp(" +
		config.Host + ":" +
		strconv.Itoa(config.Port) + ")/" +
		config.DBName
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &UGIDao{db}, nil
}

func (u *UGIDao) AddTag(sha1, tag string) error {
	ugi, err := u.Get(sha1)
	if err != nil {
		return err
	}

	updatedOn := time.Now().UnixNano() / 1000 / 1000
	if ugi == nil {
		stmt, err := u.db.Prepare("INSERT INTO ugi(sha1, tags, updated_on) VALUES(?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(sha1, tag, updatedOn)
		if err != nil {
			return err
		}
	} else {
		var newTags string
		if ugi.Tags == "" {
			newTags = tag
		} else {
			newTags = ugi.Tags + "," + tag
		}
		stmt, err := u.db.Prepare("UPDATE ugi SET tags=?,updated_on=? WHERE sha1=?")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(newTags, updatedOn, sha1)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *UGIDao) AddStory(sha1, story string) error {
	ugi, err := u.Get(sha1)
	if err != nil {
		return err
	}

	updatedOn := time.Now().UnixNano() / 1000 / 1000
	if ugi == nil {
		stmt, err := u.db.Prepare("INSERT INTO ugi(sha1, story, updated_on) VALUES(?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(sha1, story, updatedOn)
		if err != nil {
			return err
		}
	} else {
		stmt, err := u.db.Prepare("UPDATE ugi SET story=?,updated_on=? WHERE sha1=?")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(story, updatedOn, sha1)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *UGIDao) Get(sha1 string) (*model.UGI, error) {
	var tags, story sql.NullString
	var updatedOn sql.NullInt64

	stmt, err := u.db.Prepare("SELECT tags,story,updated_on FROM ugi WHERE sha1=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(sha1)
	err = row.Scan(&tags, &story, &updatedOn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	ugi := model.UGI{Sha1: sha1}
	if tags.Valid {
		ugi.Tags = tags.String
	}
	if story.Valid {
		ugi.Story = story.String
	}
	if updatedOn.Valid {
		ugi.UpdatedOn = updatedOn.Int64
	}
	return &ugi, nil
}

func (u *UGIDao) ClearUp() {
	err := u.db.Close()
	if err != nil {
		PVLog.Println("failed to close DB")
	}
}
