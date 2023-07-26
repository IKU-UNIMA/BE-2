package request

import "be-2/src/model"

type Operator struct {
	Nama    string `json:"nama" validate:"required"`
	Nip     string `json:"nip" validate:"required"`
	Email   string `json:"email" validate:"required"`
	IdProdi int    `json:"id_prodi" validate:"required"`
}

func (r *Operator) MapRequest() *model.Operator {
	return &model.Operator{
		Nama:    r.Nama,
		Nip:     r.Nip,
		IdProdi: r.IdProdi,
	}
}
