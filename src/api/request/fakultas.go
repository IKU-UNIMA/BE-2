package request

import "be-2/src/model"

type Fakultas struct {
	Nama string `json:"nama" validate:"required"`
}

func (r *Fakultas) MapRequest() *model.Fakultas {
	return &model.Fakultas{
		Nama: r.Nama,
	}
}
