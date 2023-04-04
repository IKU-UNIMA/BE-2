package request

import "be-2/src/model"

type Mahasiswa struct {
	IdProdi      int    `json:"id_prodi"`
	Nim          string `json:"nim"`
	Nama         string `json:"nama"`
	Email        string `json:"email"`
	JenisKelamin string `json:"jenis_kelamin"`
}

func (r *Mahasiswa) MapRequest() *model.Mahasiswa {
	return &model.Mahasiswa{
		IdProdi:      r.IdProdi,
		Nim:          r.Nim,
		Nama:         r.Nama,
		JenisKelamin: r.JenisKelamin,
	}
}
