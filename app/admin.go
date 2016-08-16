package app

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/qor/admin"
	"github.com/qor/i18n"
	"github.com/qor/l10n"
	"github.com/qor/media_library"
	"github.com/qor/publish"
	"github.com/qor/qor"
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
		l10n.RegisterCallbacks(DB)
		sorting.RegisterCallbacks(DB)
		validations.RegisterCallbacks(DB)
		media_library.RegisterCallbacks(DB)
	} else {
		panic(err)
	}

	Admin = admin.New(&qor.Config{DB: DB})
	Admin.SetSiteName("Auto Testing")
	Admin.AddResource(&Project{})
	Admin.AddResource(&TestCase{})
	AddWorker()
	Admin.AddResource(&ScheduledTest{})

	// Admin.SetAuth(Auth{})

}

func HostUrl() string {
	if Config.Host != "80" {
		return Config.Scheme + "://" + Config.Host + ":" + Config.Port
	}
	return Config.Scheme + "://" + Config.Host
}
