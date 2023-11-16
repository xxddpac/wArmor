package management

import (
	"fmt"
	"math"
	"rule_engine/global"
	"rule_engine/log"
	"rule_engine/model"
	"rule_engine/mysql"
	"strings"
	"time"
)

var (
	ManagerIp *_ManagerIp
	expireMap = map[int]func() time.Time{
		1: func() time.Time {
			return time.Now().Add(time.Hour * 1)
		},
		2: func() time.Time {
			return time.Now().Add(time.Hour * 8)
		},
		3: func() time.Time {
			return time.Now().AddDate(0, 0, 1)
		},
		4: func() time.Time {
			return time.Now().AddDate(0, 0, 7)
		},
	}
)

type _ManagerIp struct {
}

func (*_ManagerIp) Post(param *model.Ip) (err error) {
	defer func() {
		if err != nil {
			return
		}
		go model.Message{Event: global.Ip}.Update()
	}()
	if param.IpType == global.WhiteList || (param.IpType == global.BlackList && param.BlockType == global.Permanent) {
		if _, err = mysql.Database.Exec(
			`
		INSERT INTO ip (operator, comment, ip_type, ip_address)
		VALUES (?, ?, ?, ?)`,
			param.Operator,
			param.Comment,
			param.IpType,
			param.IpAddress,
		); err != nil {
			return
		}
	} else {
		if param.BlockType != global.Temporary {
			return fmt.Errorf("error type for block_type")
		}
		if _, ok := expireMap[param.ExpireTimeTag]; !ok {
			return fmt.Errorf("error type for expire_time")
		}
		param.ExpireTime = expireMap[param.ExpireTimeTag]().Format(global.DefaultTimeLayout)
		if _, err = mysql.Database.Exec(`
		INSERT INTO ip (operator, comment, ip_type, ip_address,block_type,expire_time)
		VALUES (?, ?, ?, ?,?,?)`, param.Operator, param.Comment, param.IpType, param.IpAddress, param.BlockType, param.ExpireTime); err != nil {
			return
		}
	}
	return
}

func (*_ManagerIp) Get(param *model.IpQuery) (interface{}, error) {
	var (
		result         model.IpQueryResult
		resp           []model.Ip
		conditionCount int
		whereClause    string
	)
	conditions := make([]string, 0)
	args := make([]interface{}, 0)

	if param.IpType != 0 {
		conditions = append(conditions, "ip_type = ?")
		args = append(args, param.IpType)
	}

	if param.BlockType != 0 {
		conditions = append(conditions, "block_type = ?")
		args = append(args, param.BlockType)
	}
	if param.Keyword != "" {
		keyword := "%" + param.Keyword + "%"
		keywordCondition := `
        (ip_address LIKE ? OR comment LIKE ? OR operator LIKE ?)
    	`
		conditions = append(conditions, keywordCondition)
		args = append(args, keyword, keyword, keyword)
	}
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	conditionCountQuery := fmt.Sprintf(`SELECT count(*)  FROM ip %s`, whereClause)
	if err := mysql.Database.Get(&conditionCount, conditionCountQuery, args...); err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`SELECT * FROM ip %s ORDER BY updated_at DESC LIMIT ? OFFSET ?`, whereClause)
	args = append(args, param.Size, (param.Page-1)*param.Size)
	if err := mysql.Database.Select(&resp, query, args...); err != nil {
		return nil, err
	}
	result.Total = conditionCount
	result.Page = param.Page
	result.Size = param.Size
	result.Items = model.IpQueryResultFunc(resp)
	result.Pages = int(math.Ceil(float64(result.Total) / float64(param.Size)))
	return result, nil
}

func (*_ManagerIp) Remove() {
	var (
		removeERR error
		ids       []int
	)
	defer func() {
		if removeERR == nil {
			go model.Message{Event: global.Ip}.Update()
		}
	}()
	query := `
        SELECT id
        FROM ip
        WHERE block_type = 2 AND expire_time < ?
    `
	if err := mysql.Database.Select(&ids, query, time.Now().Format(global.DefaultTimeLayout)); err != nil {
		log.Errorf("select for remove err:%s", err)
		removeERR = err
		return
	}
	if len(ids) == 0 {
		removeERR = fmt.Errorf("not found")
		return
	}
	for _, id := range ids {
		deleteQuery := `
        DELETE FROM ip
        WHERE id = ?
    	`
		if _, err := mysql.Database.Exec(deleteQuery, id); err != nil {
			log.Errorf("delete by id :%d err:%s", id, err)
			removeERR = err
		}
	}
}

func (*_ManagerIp) Delete(id int) (err error) {
	defer func() {
		if err != nil {
			return
		}
		go func() {
			model.Message{Event: global.Ip}.Update()
		}()
	}()
	deleteQuery := `
        DELETE FROM ip
        WHERE id = ?
    	`
	if _, err = mysql.Database.Exec(deleteQuery, id); err != nil {
		return
	}
	return
}
