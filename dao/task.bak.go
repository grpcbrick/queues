package database

import (
	"bytes"

	_ "github.com/go-sql-driver/mysql" // mysql 驱动
	"github.com/jmoiron/sqlx"
	"github.com/yinxulai/goutils/config"
)

var (
	createQueuesTableStmt        *sqlx.Stmt
	InsertTaskByChannelNamedStmt *sqlx.NamedStmt
	CountTaskByChannelNamedStmt  *sqlx.NamedStmt
	UpdateTaskByIDNamedStmt      *sqlx.NamedStmt
	DeleteTaskByIDNamedStmt      *sqlx.NamedStmt
	QueryTaskByIDNamedStmt       *sqlx.NamedStmt
	QueryTaskByOwnerNamedStmt    *sqlx.NamedStmt
	QueryTaskByHashCodeNamedStmt *sqlx.NamedStmt
	CountTaskByHashCodeNamedStmt *sqlx.NamedStmt
	CountTaskByIDNamedStmt       *sqlx.NamedStmt
	CountTaskByOwnerNamedStmt    *sqlx.NamedStmt
)

func Init() {
	var err error
	database, err := sqlx.Connect("mysql", config.MustGet("mysql"))
	if err != nil {
		panic(err)
	}

	// 设置 Name 映射方法
	database.MapperFunc(func(field string) string { return field })

	// 创建文章表
	createQueuesTableStmt = MustPreparex(database,
		" CREATE TABLE IF NOT EXISTS `queues` (",
		" `ID` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID' COMMENT 'ID',",
		" `Next` int(11) DEFAULT 0 COMMENT '下一个任务',",
		" `Prior` int(11) DEFAULT 0 COMMENT '上一个任务',",
		" `Owner` int(11) NOT NULL COMMENT '所属',",
		" `State` varchar(128) DEFAULT '' COMMENT '状态',",
		" `Input` text DEFAULT '' COMMENT '输入',",
		" `Output` text DEFAULT '' COMMENT '输出',",
		" `HashCode` varchar(128) NOT NULL COMMENT 'HashCode',",
		" `Channel` varchar(128) NOT NULL COMMENT '频道',",
		" `RetryCount` tinyint(8) DEFAULT 0 COMMENT '重试次数统计',",
		" `CreateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',",
		" `UpdateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',",
		" `RetryMaxLimit` tinyint(8) DEFAULT 3 COMMENT '最大重试次数统计',",
		" PRIMARY KEY (`ID`)",
		" ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4",
		" ;",
	)

	_, err = createQueuesTableStmt.Exec()
	if err != nil {
		panic(err)
	}

	InsertTaskByChannelNamedStmt = MustPreparexNamed(
		database,
		" INSERT INTO `queues`",
		" (`ID`, `Next`, `Prior`, `Owner`, `Input`, `HashCode`, `Channel`, `RetryMaxLimit`)",
		" VALUES",
		" (:ID, :Next, :Prior, :Owner, :Input, :HashCode, :Channel, :RetryMaxLimit)",
		" ;",
	)

	CountTaskByChannelNamedStmt = MustPreparexNamed(
		database,
		" SELECT COUNT(*) FROM `queues`",
		" WHERE `Channel` = :Channel ;",
	)

	CountTaskByIDNamedStmt = MustPreparexNamed(
		database,
		" SELECT COUNT(*) FROM `queues`",
		" WHERE `ID` = :ID",
		" ;",
	)

	CountTaskByOwnerNamedStmt = MustPreparexNamed(
		database,
		" SELECT COUNT(*) FROM `queues`",
		" WHERE `Owner` = :Owner",
		" ;",
	)

	CountTaskByHashCodeNamedStmt = MustPreparexNamed(
		database,
		" SELECT COUNT(*) FROM `queues`",
		" WHERE `HashCode` = :HashCode",
		" ;",
	)

	UpdateTaskByIDNamedStmt = MustPreparexNamed(
		database,
		" UPDATE `queues` SET",
		" `State` = :State,",
		" `Output` = :Output,",
		" `RetryCount` = :RetryCount",
		" WHERE `ID` = :ID",
		" ;",
	)

	DeleteTaskByIDNamedStmt = MustPreparexNamed(
		database,
		" DELETE FROM `queues`",
		" WHERE`ID` = :ID",
		" ;",
	)

	QueryTaskByIDNamedStmt = MustPreparexNamed(
		database,
		" SELECT * FROM `queues`",
		" WHERE `ID` = :ID",
		" ;",
	)

	QueryTaskByHashCodeNamedStmt = MustPreparexNamed(
		database,
		" SELECT * FROM `queues`",
		" WHERE `HashCode` = :HashCode",
		" ;",
	)

	QueryTaskByOwnerNamedStmt = MustPreparexNamed(
		database,
		" SELECT * FROM `queues`",
		" WHERE `Owner` = :Owner",
		" LIMIT :Limit",
		" OFFSET :Offset",
		" ;",
	)
}

// MustPreparex 解析 query
func MustPreparex(database *sqlx.DB, querys ...string) *sqlx.Stmt {
	var queryBuf bytes.Buffer

	for _, s := range querys {
		queryBuf.WriteString(s)
	}

	stmp, err := database.Preparex(queryBuf.String())
	if err != nil {
		panic(err)
	}
	return stmp
}

// MustPreparexNamed 解析 query
func MustPreparexNamed(database *sqlx.DB, querys ...string) *sqlx.NamedStmt {
	var queryBuf bytes.Buffer

	for _, s := range querys {
		queryBuf.WriteString(s)
	}

	stmp, err := database.PrepareNamed(queryBuf.String())
	if err != nil {
		panic(err)
	}
	return stmp
}
