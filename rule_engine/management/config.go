package management

import (
	"fmt"
	"rule_engine/global"
	"rule_engine/model"
	"rule_engine/mysql"
)

var ManagerConfig *_ManagerConfig

type _ManagerConfig struct {
}

func (*_ManagerConfig) Post(param *model.Config) (err error) {
	var count int
	if err = mysql.Database.QueryRow("SELECT COUNT(*) FROM config").Scan(&count); err != nil {
		return
	}
	if count == 1 {
		return fmt.Errorf("waf configuration already exists. ")
	}
	if _, err = mysql.Database.Exec(`INSERT INTO config (mode,operator) VALUES (?,?)`, param.Mode, param.Operator); err != nil {
		return
	}
	return
}

func (*_ManagerConfig) Put(query *model.QueryID, param *model.Config) (err error) {
	defer func() {
		if err != nil {
			return
		}
		go func() {
			//TODO  模式切换发送消息通知相关人
		}()
		go func() {
			model.Message{Event: global.Config}.Update()
		}()
	}()
	stmt, err := mysql.Database.Prepare("UPDATE config SET mode = ? ,operator = ? WHERE id = ?")
	if err != nil {
		return
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err = stmt.Exec(param.Mode, param.Operator, query.ID)
	if err != nil {
		return
	}
	return

}

func (*_ManagerConfig) Get() (interface{}, error) {
	var c model.Config
	type newConfig struct {
		model.Config
		ModeDesc string `json:"mode_desc"`
		Comment  string `json:"comment"`
	}
	query := `SELECT id,mode,operator,updated_at,created_at FROM config `
	if err := mysql.Database.Get(&c, query); err != nil {
		return nil, err
	}
	result := newConfig{
		Config:   c,
		ModeDesc: c.Mode.String(),  //前端模式开关按钮名称
		Comment:  c.Mode.Comment(), //高亮详细描述此模式功能
	}
	return result, nil
}
