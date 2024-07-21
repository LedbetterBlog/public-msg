package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // MySQL 驱动
	"log"
)

// MySQLPoolManager 结构体
type MySQLPoolManager struct {
	db *sql.DB
}

// NewMySQLPoolManager 创建 MySQLPoolManager 实例
func NewMySQLPoolManager(dsn string) (*MySQLPoolManager, error) {
	// 尝试连接到 MySQL 数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接 MySQL 失败: %v", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("测试 MySQL 连接失败: %v", err)
	}

	// 成功连接并测试通过
	log.Println("成功连接 MySQL")

	return &MySQLPoolManager{db: db}, nil
}

// InsertData 向 MySQL 插入数据
func (m *MySQLPoolManager) InsertData(query string, args ...interface{}) (int64, error) {
	result, err := m.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Close 关闭 MySQL 连接
func (m *MySQLPoolManager) Close() error {
	return m.db.Close()
}
