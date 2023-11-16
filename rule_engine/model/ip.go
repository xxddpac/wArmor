package model

import "rule_engine/global"

// Ip 黑白名单
type Ip struct {
	Base
	Operator      string              `db:"operator" json:"operator" binding:"required"`
	Comment       string              `db:"comment" json:"comment" binding:"required"`
	BlockType     global.WafBlockType `db:"block_type" json:"block_type"`
	ExpireTime    string              `db:"expire_time" json:"expire_time"`
	IpAddress     string              `db:"ip_address" json:"ip_address" binding:"required"`
	IpType        global.WafIpType    `db:"ip_type" json:"ip_type" binding:"required,oneof=1 2"`
	ExpireTimeTag int                 `json:"expire_time_tag"`
}

//建表语句
//`
//		CREATE TABLE IF NOT EXISTS ip (
//			id INT AUTO_INCREMENT PRIMARY KEY COMMENT '黑白名单主键ID',
//			operator VARCHAR(10) NOT NULL COMMENT '操作人' ,
//		    comment VARCHAR(255) NOT NULL COMMENT '备注信息' ,
//			ip_type INT NOT NULL COMMENT '黑白名单类型',
//			block_type INT DEFAULT 1 COMMENT '封禁类型',
//			ip_address VARCHAR(15) NOT NULL COMMENT '黑白名单IP' ,
//		    expire_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '黑名单封禁时间',
//			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
//			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
//		)
//	`

type IpQuery struct {
	QueryPage
	IpType    global.WafIpType    `form:"ip_type"`
	BlockType global.WafBlockType `form:"block_type"`
	Keyword   string              `form:"keyword"`
}

func IpQueryFunc() *IpQuery {
	return &IpQuery{
		QueryPage: QueryPage{Page: 1, Size: 10},
	}
}

type IpQueryDto struct {
	Ip
	Base
	IpTypeDesc    string `json:"ip_type_desc"`
	BlockTypeDesc string `json:"block_type_desc"`
}

type IpQueryResult struct {
	QueryResult
	Items []IpQueryResultDto `json:"items"`
}

type IpQueryResultDto struct {
	IpQueryDto
}

func IpQueryResultFunc(i []Ip) []IpQueryResultDto {
	var (
		resp   IpQueryResultDto
		result []IpQueryResultDto
	)
	for _, item := range i {
		resp.Ip = item
		resp.Base = item.Base
		resp.BlockTypeDesc = item.BlockType.String()
		resp.IpTypeDesc = item.IpType.String()
		result = append(result, resp)
	}
	return result
}
