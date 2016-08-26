package app

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/qor/admin"
	"github.com/qor/i18n"
	"github.com/qor/media_library"
	"github.com/qor/publish"
	"github.com/qor/qor"
	"github.com/qor/roles"
	"github.com/qor/sorting"
	"github.com/qor/validations"
)

var Config = struct {
	Scheme string `default:"http"`
	Host   string `default:"localhost"`
	Port   string `default:"8000"`
	DB     struct {
		Name     string `default:"testing"`
		Host     string `default:"localhost"`
		Port     string `default:"5432"`
		Adapter  string `default:"postgres"`
		User     string `default:"app"`
		Password string `default:"1234"`
	}
	I18n *i18n.I18n
}{}
var Root string

var Admin *admin.Admin

var (
	DB      *gorm.DB
	Publish *publish.Publish
)

func init() {
	Root = path.Join(os.Getenv("GOPATH"), "/src/github.com/azumads/selenium")
	if err := configor.Load(&Config, filepath.Join(Root, "app/config.yml")); err != nil {
		panic(err)
	}

	var err error
	var db *gorm.DB

	dbConfig := Config.DB
	if Config.DB.Adapter == "mysql" {
		db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name))
	} else if Config.DB.Adapter == "postgres" {
		fmt.Printf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
		db, err = gorm.Open("postgres", fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name))
	} else {
		panic(errors.New("not supported database adapter"))
	}

	if err == nil {
		DB = db
		DB.LogMode(true)
		Publish = publish.New(DB)
		sorting.RegisterCallbacks(DB)
		validations.RegisterCallbacks(DB)
		media_library.RegisterCallbacks(DB)
	} else {
		panic(err)
	}
	roles.Register("user", func(req *http.Request, currentUser interface{}) bool {
		return true
	})

	Admin = admin.New(&qor.Config{DB: DB})
	Admin.SetSiteName("Auto Testing")
	Admin.AddResource(&Project{})
	Admin.AddResource(&TestCase{})
	AddWorker()
	scheduledTest := Admin.AddResource(&ScheduledTest{}, &admin.Config{Permission: roles.Deny(roles.Create, "user")})
	scheduledTest.Meta(&admin.Meta{Name: "JobId", Permission: roles.Allow(roles.Read, "user")})
	scheduledTest.Meta(&admin.Meta{Name: "Project", Permission: roles.Allow(roles.Read, "user")})

	// Admin.SetAuth(Auth{})

}

func IsProd() bool {
	return configor.ENV() == "production"
}

func HostUrl() string {
	if Config.Host != "80" {
		return Config.Scheme + "://" + Config.Host + ":" + Config.Port
	}
	return Config.Scheme + "://" + Config.Host
}
