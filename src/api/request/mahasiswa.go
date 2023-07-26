package request

import "be-2/src/model"

type Mahasiswa struct {
	IdProdi      int    `json:"id_prodi" validate:"required"`
	Nim          string `json:"nim" validate:"required"`
	Nama         string `json:"nama" validate:"required"`
	JenisKelamin string `json:"jenis_kelamin" validate:"required"`
}

func (r *Mahasiswa) MapRequest() *model.Mahasiswa {
	return &model.Mahasiswa{
		IdProdi:      r.IdProdi,
		Nim:          r.Nim,
		Nama:         r.Nama,
		JenisKelamin: r.JenisKelamin,
	}
}
