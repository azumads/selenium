package app

// Tables returns a slice of all tables.
var Tables = []interface{}{
	&ScheduledTest{},
	&Project{},
	&TestCase{},
}

func init() {
	AutoMigrate(Tables...)
}

func AutoMigrate(tables ...interface{}) {
	for _, table := range tables {
		DB.AutoMigrate(table)
	}
}
