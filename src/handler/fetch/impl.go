package fetch

import (
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	"Yearning-go/src/parser"
	"encoding/hex"
	"fmt"
	"github.com/jinzhu/gorm"
	"net"
	"strconv"
	"strings"
)

const (
	UNDO_EXPR            = "username =? AND work_id =? AND `status` =? "
	UNDO_MESSAGE_ERROR   = "工单状态已更改！无法撤销"
	UNDO_MESSAGE_SUCCESS = "工单已撤销！"
	AUDITOR_IS_NOT_EXIST = "流程信息缺失,请检查该数据源流程配置!"
)

type referOrder struct {
	Data model.CoreSqlOrder `json:"data"`
	SQLs string             `json:"sqls"`
	Tp   int                `json:"tp"`
}

type _FetchBind struct {
	IDC      string             `json:"idc"`
	Tp       string             `json:"tp"`
	Source   string             `json:"source"`
	DataBase string             `json:"data_base"`
	Table    string             `json:"table"`
	Rows     []parser.FieldInfo `json:"rows"`
	Idx      []parser.IndexInfo `json:"idx"`
}

func (u *_FetchBind) FetchTableFieldsOrIndexes() error {
	var s model.CoreDataSource

	model.DB().Where("source =?", u.Source).First(&s)

	ps := lib.Decrypt(s.Password)
	host := s.IP
	if "" != s.Proxy {
		host = strings.Join([]string{lib.PROXY_PREFIX, s.Proxy, "_", hex.EncodeToString([]byte(net.JoinHostPort(s.IP, strconv.Itoa(s.Port))))}, "")
	}

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", s.Username, ps, host, s.Port, u.DataBase))
	if err != nil {
		return err
	}

	defer db.Close()

	if err := db.Raw(fmt.Sprintf("SHOW FULL FIELDS FROM `%s`.`%s`", u.DataBase, u.Table)).Scan(&u.Rows).Error; err != nil {
		return err
	}

	if err := db.Raw(fmt.Sprintf("SHOW INDEX FROM `%s`.`%s`", u.DataBase, u.Table)).Scan(&u.Idx).Error; err != nil {
		return err
	}
	return nil
}
