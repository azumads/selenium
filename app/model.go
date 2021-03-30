package app

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/media_library"
)

type ScheduledTest struct {
	gorm.Model
	Project   Project
	ProjectID uint
	JobId     string
	LoopHour  string
	NextRun   time.Time
}

type Project struct {
	gorm.Model
	Name        string
	NotifyEmail string
}

type TestCase struct {
	gorm.Model
	Project   Project
	ProjectID uint
	Name      string
	TestFile  media_library.FileSystem
	CsvFile   media_library.FileSystem
}
