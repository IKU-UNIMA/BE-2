package request

import "be-2/src/model"

type KategoriProgramKm struct {
	Nama string `json:"nama" validate:"required"`
}

func (r *KategoriProgramKm) MapRequest() *model.KategoriProgramKm {
	return &model.KategoriProgramKm{Nama: r.Nama}
}
