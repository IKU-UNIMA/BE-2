package request

import "be-2/src/model"

type Prestasi struct {
	IdSemester        int    `form:"id_semester" validate:"required"`
	IdDosenPembimbing int    `form:"id_dosen_pembimbing"`
	Nama              string `form:"nama" validate:"required"`
	TingkatPrestasi   string `form:"tingkat_prestasi" validate:"required"`
	Penyelenggara     string `form:"penyelenggara" validate:"required"`
	Peringkat         string `form:"peringkat" validate:"required"`
}

func (r *Prestasi) MapRequest(idMahasiswa int, sertifikat string) *model.Prestasi {
	return &model.Prestasi{
		IdMahasiswa:       idMahasiswa,
		IdSemester:        r.IdSemester,
		IdDosenPembimbing: r.IdDosenPembimbing,
		Nama:              r.Nama,
		TingkatPrestasi:   r.TingkatPrestasi,
		Penyelenggara:     r.Penyelenggara,
		Peringkat:         r.Peringkat,
		Sertifikat:        sertifikat,
	}
}
