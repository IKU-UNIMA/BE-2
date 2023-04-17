package database

import (
	"be-2/src/config/env"
	"be-2/src/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitMySQL() *gorm.DB {
	db, err := gorm.Open(mysql.Open(env.GetMySQLEnv()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	return db
}

func MigrateMySQL() {
	InitMySQL().AutoMigrate(
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
	)
}
