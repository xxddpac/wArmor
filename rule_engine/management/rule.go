package management

import (
	"database/sql"
	"fmt"
	"math"
	"rule_engine/global"
	"rule_engine/model"
	"rule_engine/mysql"
	"rule_engine/utils"
	"strings"
)

var ManagerRule *_ManagerRule

type _ManagerRule struct {
}

func (*_ManagerRule) Get(param *model.RuleQuery) (interface{}, error) {
	var (
		result         model.RuleQueryResult
		resp           []*model.Rule
		conditionCount int
		whereClause    string
	)
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	if len(param.Status) != 0 {
		conditions = append(conditions, "status = ?")
		args = append(args, utils.StrToBool(param.Status))
	}
	if len(param.RulesOperation) != 0 {
		conditions = append(conditions, "rules_operation = ?")
		args = append(args, param.RulesOperation)
	}
	if param.RuleVariable != 0 {
		conditions = append(conditions, "rule_variable = ?")
		args = append(args, param.RuleVariable)
	}
	if param.RuleType != 0 {
		conditions = append(conditions, "rule_type = ?")
		args = append(args, param.RuleType)
	}
	if param.RuleAction != 0 {
		conditions = append(conditions, "rule_action = ?")
		args = append(args, param.RuleAction)
	}
	if param.Severity != 0 {
		conditions = append(conditions, "severity = ?")
		args = append(args, param.Severity)
	}
	if param.Keyword != "" {
		keyword := "%" + param.Keyword + "%"
		keywordCondition := `
        (description LIKE ? OR operator LIKE ?)
    	`
		conditions = append(conditions, keywordCondition)
		args = append(args, keyword, keyword)
	}
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	conditionCountQuery := fmt.Sprintf(`SELECT count(*)  FROM rule %s`, whereClause)
	if err := mysql.Database.Get(&conditionCount, conditionCountQuery, args...); err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`SELECT * FROM rule %s ORDER BY updated_at DESC LIMIT ? OFFSET ?`, whereClause)
	args = append(args, param.Size, (param.Page-1)*param.Size)
	if err := mysql.Database.Select(&resp, query, args...); err != nil {
		return nil, err
	}
	for _, item := range resp {
		newQuery := fmt.Sprintf(`SELECT rules FROM rules WHERE rule_id = %d`, item.ID)
		if err := mysql.Database.Select(&item.Rules, newQuery); err != nil {
			return nil, err
		}
	}
	result.Total = conditionCount
	result.Page = param.Page
	result.Size = param.Size
	result.Items = model.RuleQueryResultFunc(resp)
	result.Pages = int(math.Ceil(float64(result.Total) / float64(param.Size)))
	return result, nil
}

func (*_ManagerRule) Post(param *model.Rule) (err error) {
	var (
		lastInsertID int64
		tx           *sql.Tx
	)
	//开启事务
	tx, err = mysql.Database.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		if err = tx.Commit(); err != nil {
			return
		}
		go model.Message{Event: global.Rule}.Update()
	}()
	result, err := tx.Exec(`
        INSERT INTO rule (
            operator, rule_variable, rule_type, status, rule_action, description, severity,rules_operation
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `, param.Operator, param.RuleVariable, param.RuleType, param.Status, param.RuleAction, param.Description, param.Severity, param.RulesOperation)
	if err != nil {
		return
	}
	lastInsertID, err = result.LastInsertId()
	if err != nil {
		return
	}
	for _, item := range param.Rules {
		_, err = tx.Exec(`
            INSERT INTO rules (rule_id, rules) VALUES (?, ?)
        `, lastInsertID, item)
		if err != nil {
			return
		}
	}
	return
}

func (*_ManagerRule) Put(query *model.QueryID, param *model.Rule) (err error) {
	var tx *sql.Tx
	tx, err = mysql.Database.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		if err = tx.Commit(); err != nil {
			return
		}
		go model.Message{Event: global.Rule}.Update()
	}()
	_, err = tx.Exec(`
        UPDATE rule
        SET operator=?, rule_variable=?, rule_type=?, status=?, rule_action=?, description=?, severity=?, rules_operation=?
        WHERE id=?
    `, param.Operator, param.RuleVariable, param.RuleType, param.Status, param.RuleAction, param.Description, param.Severity, param.RulesOperation, query.ID)
	if err != nil {
		return
	}
	_, err = tx.Exec(`
        DELETE FROM rules WHERE rule_id=?
    `, query.ID)
	if err != nil {
		return
	}
	for _, item := range param.Rules {
		_, err = tx.Exec(`
            INSERT INTO rules (rule_id, rules) VALUES (?, ?)
        `, query.ID, item)
		if err != nil {
			return
		}
	}
	return
}

func (*_ManagerRule) Delete(query *model.QueryID) (err error) {
	defer func() {
		if err != nil {
			return
		}
		go model.Message{Event: global.Rule}.Update()
	}()
	deleteQuery := `
        DELETE FROM rule
        WHERE id = ?
    	`
	if _, err = mysql.Database.Exec(deleteQuery, query.ID); err != nil {
		return
	}
	return
}
