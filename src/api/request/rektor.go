package request

import "be-2/src/model"

type Rektor struct {
	Nama  string `json:"nama" validate:"required"`
	Nip   string `json:"nip" validate:"required"`
	Email string `json:"email" validate:"required"`
}

func (r *Rektor) MapRequest() *model.Rektor {
	return &model.Rektor{
		Nama: r.Nama,
		Nip:  r.Nip,
	}
}
