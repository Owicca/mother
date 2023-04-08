package infra

import (
	"fmt"

	"gocloud.dev/mysql"
	"gocloud.dev/postgres"
	"gorm.io/gorm"
)

func GetDbConn(DbHost string, DbPort string, DbName string, DbUser string, DbPassword string) (*gorm.DB, error) {
	connectionString := GetConnString("mysql", DbHost, DbPort, DbName, DbUser, DbPassword)
	dialector := GetDialector("mysql", connectionString)
	conn, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Database connection failed %s", err)
	}

	// coupling
	var tableList []string
	conn.Raw("SHOW TABLES LIKE 'posts'").Scan(&tableList)
	if len(tableList) == 0 {
		CreateDbSchema(conn)
	}

	return conn, nil
}

func GetDialector(db string, connString string) gorm.Dialector {
	if db == "postgresql" {
		return postgres.Open(connString)
	}

	return mysql.Open(connString)
}

func GetConnString(db string, DbHost string, DbPort string, DbName string, DbUser string, DbPassword string) string {
	if db == "postgresql" {
		return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
			DbHost, DbPort, DbName, DbUser, DbPassword)
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		DbUser, DbPassword, DbHost, DbPort, DbName)
}
