package model

import (
	"rule_engine/global"
)

type Rule struct {
	Base
	RuleVariable   global.WafRuleVariable `db:"rule_variable" json:"rule_variable" binding:"required" enums:"1,2,3,4,5,6"`
	Operator       string                 `db:"operator" json:"operator" binding:"required"`
	RuleType       global.WafRuleType     `db:"rule_type" json:"rule_type" binding:"required" enums:"1,2,3,4,5,6,7,8,9,10,11,12,13,14,15"`
	Status         *bool                  `db:"status" json:"status" binding:"required"`
	Rules          []string               `db:"rules" json:"rules" binding:"required" example:"规则正则表达式列表"`
	RuleAction     global.WafRuleAction   `db:"rule_action" json:"rule_action" binding:"required" enums:"1,2,3"`
	Description    string                 `db:"description" json:"description" binding:"required"`
	RulesOperation string                 `db:"rules_operation" json:"rules_operation" binding:"required" example:"输入 and 或 or,and需全部匹配Rules列表,or匹配Rules列表任意一个"`
	Severity       global.WafRuleSeverity `db:"severity" json:"severity" binding:"required" enums:"1,2,3,4,5"`
}

//建表语句
//主表
//`
//		CREATE TABLE IF NOT EXISTS rule (
//			id INT AUTO_INCREMENT PRIMARY KEY COMMENT '规则主键ID',
//			operator VARCHAR(10) NOT NULL COMMENT '操作人' ,
//			rule_variable INT NOT NULL COMMENT '规则变量',
//		    rule_type INT NOT NULL COMMENT '规则类型',
//		    status BOOLEAN NOT NULL COMMENT '规则开启状态',
//		    rule_action INT NOT NULL COMMENT '规则执行动作',
//		    description VARCHAR(255) NOT NULL COMMENT '规则描述',
//		    severity INT NOT NULL COMMENT '风险级别',
//		    rules_operation VARCHAR(3) NOT NULL COMMENT '规则匹配条件(and/or)',
//			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
//			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
//		)
//	`

//关联表
//`CREATE TABLE rules (
//    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '规则匹配详情主键ID',
//    rule_id INT NOT NULL COMMENT '规则表外键',
//    rules VARCHAR(255) NOT NULL COMMENT '规则详情',
//    FOREIGN KEY (rule_id) REFERENCES rule(id) ON DELETE CASCADE,
//    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
//    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
//)`

type RuleQuery struct {
	QueryPage
	Keyword        string                 `form:"keyword"`
	Status         string                 `form:"status"`
	RulesOperation string                 `form:"rules_operation"`
	RuleVariable   global.WafRuleVariable `form:"rule_variable"`
	RuleType       global.WafRuleType     `form:"rule_type"`
	RuleAction     global.WafRuleAction   `form:"rule_action"`
	Severity       global.WafRuleSeverity `form:"severity"`
}

func RuleQueryFunc() *RuleQuery {
	return &RuleQuery{
		QueryPage: QueryPage{Page: 1, Size: 10},
	}
}

type RuleQueryDto struct {
	Rule
	Base
	RuleVariableDesc string `json:"rule_variable_desc"`
	RuleTypeDesc     string `json:"rule_type_desc"`
	RuleActionDesc   string `json:"rule_action_desc"`
	SeverityDesc     string `json:"severity_desc"`
}

type RuleQueryResult struct {
	QueryResult
	Items []RuleQueryResultDto `json:"items"`
}

type RuleQueryResultDto struct {
	RuleQueryDto
}

func RuleQueryResultFunc(r []*Rule) []RuleQueryResultDto {
	var (
		resp   RuleQueryResultDto
		result []RuleQueryResultDto
	)
	for _, item := range r {
		resp.Rule = *item
		resp.Base = item.Base
		resp.RuleVariableDesc = item.RuleVariable.String()
		resp.RuleTypeDesc = item.RuleType.String()
		resp.RuleActionDesc = item.RuleAction.String()
		resp.SeverityDesc = item.Severity.String()
		result = append(result, resp)
	}
	return result
}
