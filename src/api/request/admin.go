package request

import "be-2/src/model"

type Admin struct {
	Nama   string `json:"nama" validate:"required"`
	Nip    string `json:"nip" validate:"required"`
	Email  string `json:"email" validate:"required"`
	Bagian string `json:"bagian" validate:"required"`
}

func (r *Admin) MapRequest() *model.Admin {
	return &model.Admin{
		Nama:   r.Nama,
		Nip:    r.Nip,
		Bagian: r.Bagian,
	}
}
