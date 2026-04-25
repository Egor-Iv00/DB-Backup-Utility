package DBinterface

import "fmt"

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	FilePath string
	DBtype   string
}

func CreateConnectionString(conf Config) string {
	switch conf.DBtype {
	case "postgres":
		{
			return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
				conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)
		}
	case "mysql":
		{
			return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=5s",
				conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)
		}
	}
	return "nil"
}
