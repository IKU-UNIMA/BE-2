package model

type Akun struct {
	ID        int       `gorm:"primaryKey"`
	Email     string    `gorm:"type:varchar(255);unique"`
	Password  string    `gorm:"type:varchar(255)"`
	Role      string    `gorm:"type:enum('rektor', 'admin', 'operator', 'mahasiswa')"`
	Admin     Admin     `gorm:"foreignKey:ID;constraint:OnDelete:CASCADE"`
	Operator  Operator  `gorm:"foreignKey:ID;constraint:OnDelete:CASCADE"`
	Rektor    Rektor    `gorm:"foreignKey:ID;constraint:OnDelete:CASCADE"`
	Mahasiswa Mahasiswa `gorm:"foreignKey:ID;constraint:OnDelete:CASCADE"`
}
