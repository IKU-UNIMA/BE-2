package database

import (
	"be-2/src/config/env"
	"be-2/src/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitMySQL() {
	var err error
	DB, err = gorm.Open(mysql.Open(env.GetMySQLEnv()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
}

func MigrateMySQL() {
	DB.AutoMigrate(
		&model.Fakultas{},
		&model.Prodi{},
		&model.Semester{},
		&model.Akun{},
		&model.Admin{},
		&model.Rektor{},
		&model.Operator{},
		&model.Mahasiswa{},
		&model.Dosen{},
		&model.Prestasi{},
		&model.KategoriProgramKm{},
		&model.KampusMerdeka{},
		&model.Target{},
	)
}
