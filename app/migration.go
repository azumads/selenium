package app

// Tables returns a slice of all tables.
var Tables = []interface{}{
	&AutoTest{},
}

func init() {
	AutoMigrate(Tables...)
}

func AutoMigrate(tables ...interface{}) {
	for _, table := range tables {
		DB.AutoMigrate(table)
	}
}
