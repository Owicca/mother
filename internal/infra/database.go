package infra

import (
	"os"

	"gorm.io/gorm"
)

const (
	DefaultDbName = "imageboard"
)

func CreateDbSchema(db *gorm.DB) {
	data, _ := os.ReadFile("./db_schema.my.sql")

	_ = db.Exec(string(data))
}

func DeleteDb(db *gorm.DB) {
	db.Exec("DROP DATABASE " + DefaultDbName)
}

func ClearDb(db *gorm.DB) {
	tables := []string{
		//"pair_to_role",
		//"action_to_object",
		"links",

		"posts",
		"threads",
		"boards",
		"topics",
		//"users",

		//"roles",
		"log_actions",
		//"objects",
		//"actions",
		"media",
	}

	for _, name := range tables {
		db.Exec("DELETE FROM " + name)
	}
}
