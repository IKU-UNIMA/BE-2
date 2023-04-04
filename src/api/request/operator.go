package request

import "be-2/src/model"

type Operator struct {
	Nama    string `json:"nama"`
	Nip     string `json:"nip"`
	Email   string `json:"email"`
	IdProdi int    `json:"id_prodi"`
}

func (r *Operator) MapRequest() *model.Operator {
	return &model.Operator{
		Nama:    r.Nama,
		Nip:     r.Nip,
		IdProdi: r.IdProdi,
	}
}
