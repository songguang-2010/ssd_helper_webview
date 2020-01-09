package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"lib/serror"
)

type Model struct {
	db        *gorm.DB
	tableName string
	Error     error
}

func (m *Model) SetError(msg string) error {
	m.Error = serror.New(msg)
	return nil
}

func (m *Model) SetTableName(tableName string) error {
	m.tableName = tableName
	return nil
}

func (m *Model) GetTableName() (string, error) {
	if m.tableName == "" {
		return "", serror.New("table name is empty")
	}
	return m.tableName, nil
}

func (m *Model) OpenDB(dbname string) error {
	//获取数据库配置
	db_name := viper.GetString(dbname + ".name")
	db_username := viper.GetString(dbname + ".username")
	db_password := viper.GetString(dbname + ".password")
	db_address := viper.GetString(dbname + ".address")
	db_port := viper.GetString(dbname + ".port")
	db_charset := viper.GetString(dbname + ".charset")

	msn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", db_username, db_password, db_address, db_port, db_name, db_charset)
	db, err := gorm.Open("mysql", msn)
	if err != nil {
		return err
	}
	m.db = db
	return nil
}

func (m *Model) CloseDB() error {
	err := m.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) GetDB() (*gorm.DB, error) {
	return m.db, nil
}

func (m *Model) SetDB(db *gorm.DB) error {
	m.db = db
	return nil
}

func (oi *Model) TrancCommit() error {
	db, err := oi.GetDB()
	if err != nil {
		return err
	}
	resCommit := db.Commit()
	if resCommit.Error != nil {
		return resCommit.Error
	}
	return nil
}

func (oi *Model) TrancRollback() error {
	db, err := oi.GetDB()
	if err != nil {
		return err
	}
	db.Rollback()
	return nil
}

func (oi *Model) TrancBegin() error {
	db, err := oi.GetDB()
	if err != nil {
		return err
	}
	tx := db.Begin()
	err = oi.SetDB(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
