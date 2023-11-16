package model

import "rule_engine/global"

type Config struct {
	Base
	Operator string               `db:"operator" json:"operator" binding:"required"`
	Mode     global.WafConfigType `db:"mode" json:"mode" binding:"required,oneof=1 2 3"`
}

//建表语句
//`
//		CREATE TABLE IF NOT EXISTS config (
//			id INT AUTO_INCREMENT PRIMARY KEY COMMENT '配置模式主键ID',
//			operator VARCHAR(10) NOT NULL COMMENT '操作人' ,
//			mode INT NOT NULL COMMENT '配置模式',
//			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
//			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
//		)
//	`
