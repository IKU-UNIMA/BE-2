package request

import "be-2/src/model"

type Semester struct {
	Nama string `json:"nama"`
}

func (r *Semester) MapRequest() *model.Semester {
	return &model.Semester{Nama: r.Nama}
}
