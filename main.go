package main

import (
	"fmt"
	"net/url"
	"os"
	_ "system_service/routers"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"

	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "github.com/lib/pq"
)

func setup() {
	fmt.Println("Setting up the database connection...")
	// 1. Fetch values from app.conf
	// val, _ := beego.AppConfig.GetSection("default")
	user, err := beego.AppConfig.String("pguser")
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch configuration: %v", err))
	}
	pass, err := beego.AppConfig.String("pgpass")
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch configuration: %v", err))
	}
	host, err := beego.AppConfig.String("pghost")
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch configuration: %v", err))
	}
	port, err := beego.AppConfig.String("pgport")
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch configuration: %v", err))
	}
	dbName, err := beego.AppConfig.String("pgdb")
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch configuration: %v", err))
	}
	pgTest, err := beego.AppConfig.String("pgtest")
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch configuration: %v", err))
	}
	logs.Info("Postgres test value is: ", pgTest)

	// 2. Register the driver name
	orm.RegisterDriver("postgres", orm.DRPostgres)

	if user == "" {
		user = os.Getenv("DB_USER")
	}
	if pass == "" {
		pass = os.Getenv("DB_PASSWORD")
	}
	if host == "" {
		host = os.Getenv("DB_HOST")
	}
	if port == "" {
		port = os.Getenv("DB_PORT")
	}
	if dbName == "" {
		dbName = os.Getenv("DB_NAME")
	}

	fmt.Printf("Database credentials are: user=%s, host=%s, port=%s, dbName=%s\n", user, host, port, dbName)

	// 3. Construct connection string (handle special characters in password)
	dataSource := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=require",
		user, url.QueryEscape(pass), host, port, dbName,
	)

	// 4. Register the default database alias
	// Max connection pool limits can be configured here as the 4th/5th optional arguments
	dberr := orm.RegisterDataBase("default", "postgres", dataSource)
	if dberr != nil {
		panic(fmt.Sprintf("Failed to register database: %v", dberr))
	}
}

func main() {
	logs.SetLogger(logs.AdapterConsole)

	beego.LoadAppConfig("ini", "/app/conf/app.conf")
	// sqlConn, err := beego.AppConfig.String("sqlconn")
	// if err != nil {
	// 	logs.Error("%s", err)
	// }
	// orm.RegisterDataBase("default", "mysql", sqlConn)
	// logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/system_application.log"}`)

	setup()

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:8000"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	orm.Debug = true

	beego.Run()
}
