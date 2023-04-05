package request

import "be-2/src/model"

type Prestasi struct {
	IdMahasiswa       int    `form:"id_mahasiswa"`
	IdSemester        int    `form:"id_semester"`
	IdDosenPembimbing int    `form:"id_dosen_pembimbing"`
	Nama              string `form:"nama"`
	TingkatPrestasi   string `form:"tingkat_prestasi"`
	Penyelenggara     string `form:"penyelenggara"`
	Peringkat         string `form:"peringkat"`
}

func (r *Prestasi) MapRequest(idFakultas, idProdi int, sertifikat string) *model.Prestasi {
	return &model.Prestasi{
		IdMahasiswa:       r.IdMahasiswa,
		IdFakultas:        idFakultas,
		IdProdi:           idProdi,
		IdSemester:        r.IdSemester,
		IdDosenPembimbing: r.IdDosenPembimbing,
		Nama:              r.Nama,
		TingkatPrestasi:   r.TingkatPrestasi,
		Penyelenggara:     r.Penyelenggara,
		Peringkat:         r.Peringkat,
		Sertifikat:        sertifikat,
	}
}
