package models

import ti "github.com/gage-technologies/gigo-lib/db"

func init() {
	_, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		panic(err)
	}
}
