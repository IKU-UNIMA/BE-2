package request

import "be-2/src/model"

type Prestasi struct {
	IdSemester        int    `form:"id_semester"`
	IdDosenPembimbing int    `form:"id_dosen_pembimbing"`
	Nama              string `form:"nama"`
	TingkatPrestasi   string `form:"tingkat_prestasi"`
	Penyelenggara     string `form:"penyelenggara"`
	Peringkat         string `form:"peringkat"`
}

func (r *Prestasi) MapRequest(idMahasiswa, idFakultas, idProdi int, sertifikat string) *model.Prestasi {
	return &model.Prestasi{
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
