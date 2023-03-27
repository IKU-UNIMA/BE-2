package model

type Semester struct {
	ID   int    `gorm:"primaryKey"`
	Nama string `gorm:"varchar(120)"`
}
