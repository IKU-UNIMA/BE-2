package request

import "be-2/src/model"

type Prodi struct {
	IdFakultas int    `json:"id_fakultas" validate:"required"`
	Nama       string `json:"nama" validate:"required"`
	Jenjang    string `json:"jenjang" validate:"required"`
}

func (r *Prodi) MapRequest() *model.Prodi {
	return &model.Prodi{
		IdFakultas: r.IdFakultas,
		Nama:       r.Nama,
		Jenjang:    r.Jenjang,
	}
}
