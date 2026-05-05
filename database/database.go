package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func InitDB() {
	dsn := "host=localhost port=5432 user=hero password=123456 dbname=pg_practice sslmode=disable"
	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("数据库连接失败: ", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("数据库连接失败: ", err)
	}
	fmt.Println("数据库连接成功")

	createTables()
	migrateTables()
	seedData()
}

func createTables() {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id   SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			age  INT NOT NULL CHECK (age >= 0 AND age <= 150)
		)`,
		`CREATE TABLE IF NOT EXISTS modules (
			id         VARCHAR(32) PRIMARY KEY,
			name       VARCHAR(200) NOT NULL,
			type       VARCHAR(50) NOT NULL,
			version    VARCHAR(50) NOT NULL,
			file_name  VARCHAR(500) NOT NULL,
			file_path  VARCHAR(1000) NOT NULL,
			file_size  BIGINT NOT NULL,
			changelog  TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS version_histories (
			id         SERIAL PRIMARY KEY,
			module_id  VARCHAR(32) NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
			version    VARCHAR(50) NOT NULL,
			file_name  VARCHAR(500) NOT NULL,
			file_size  BIGINT NOT NULL,
			changelog  TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS app_updates (
			id           SERIAL PRIMARY KEY,
			has_update   BOOLEAN NOT NULL DEFAULT TRUE,
			version_code INT NOT NULL,
			version_name VARCHAR(50) NOT NULL,
			download_url VARCHAR(1000) NOT NULL,
			changelog    TEXT,
			force_update BOOLEAN NOT NULL DEFAULT FALSE,
			file_size    BIGINT NOT NULL
		)`,
	}
	for _, s := range stmts {
		if _, err := DB.Exec(s); err != nil {
			log.Fatal("建表失败: ", err)
		}
	}
	fmt.Println("数据表创建完成")
}

func migrateTables() {
	migrations := []string{
		`ALTER TABLE modules ADD COLUMN IF NOT EXISTS code VARCHAR(100)`,
		`ALTER TABLE modules ADD COLUMN IF NOT EXISTS download_url VARCHAR(1000)`,
	}
	for _, m := range migrations {
		DB.Exec(m)
	}
}

func seedData() {
	// 插入默认用户
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		DB.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", "张三", 25)
		DB.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", "李四", 30)
		fmt.Println("已插入默认用户数据")
	}

	// 插入默认App更新信息
	DB.QueryRow("SELECT COUNT(*) FROM app_updates").Scan(&count)
	if count == 0 {
		DB.Exec(`INSERT INTO app_updates (has_update, version_code, version_name, download_url, changelog, force_update, file_size)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			true, 10006, "1.0.6",
			"http://192.168.5.67:8080/download/HotelRepack-release-arm64-v8a.apk",
			"修复若干问题", false, 74442780)
		fmt.Println("已插入默认App更新信息")
	}
}
