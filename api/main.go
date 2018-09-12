package main

import (
	"errors"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"./core"
)

// DB DB
var (
	Host string
	Name string
	User string
	Pass string
	Port string
)

// Server Server Settings
type Server struct {
	UseSoftDelete bool
}

func loadEnvFile(envPath string) (err error) {
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return errors.New("Not Exist")
	}
	e := godotenv.Load(envPath)
	mustExist := []string{"API_PORT", "DB_HOST", "DB_NAME", "DB_USER", "DB_PASS", "DB_PORT"}
	for _, key := range mustExist {
		if os.Getenv(key) == "" {
			panic("This key must exist: " + key)
		}
	}
	if e != nil {
		return e
	}
	return nil
}

func setDB() {
	Host = os.Getenv("DB_HOST")
	Name = os.Getenv("DB_NAME")
	User = os.Getenv("DB_USER")
	Pass = os.Getenv("DB_PASS")
	Port = os.Getenv("DB_PORT")
}

func setRequestAndResponse() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "apikey", "Content-Type"},
		ExposeHeaders:   []string{"Content-Length", "Access-Control-Allow-Origin"},
	}))
	beego.BConfig.CopyRequestBody = true
}

func init() {
	loadEnvFile("../.env")
	setDB()
	setRequestAndResponse()
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", User+":"+Pass+"@tcp("+Host+":"+Port+")/"+Name+"?charset=utf8mb4&sql_mode=TRADITIONAL")
	core.InitLogger()
}

func main() {
	beego.Run(":" + os.Getenv("API_PORT"))
}
