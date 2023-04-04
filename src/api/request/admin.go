package request

import "be-2/src/model"

type Admin struct {
	Nama   string `json:"nama"`
	Nip    string `json:"nip"`
	Email  string `json:"email"`
	Bagian string `json:"bagian"`
}

func (r *Admin) MapRequest() *model.Admin {
	return &model.Admin{
		Nama:   r.Nama,
		Nip:    r.Nip,
		Bagian: r.Bagian,
	}
}
