package request

import (
	"be-2/src/model"
	"be-2/src/util"
	"errors"
)

type KampusMerdeka struct {
	IdSemester              int     `form:"id_semester" validate:"required"`
	IdDosenPembimbing       int     `form:"id_dosen_pembimbing" validate:"required"`
	IdKategoriProgram       int     `form:"id_kategori_program" validate:"required"`
	StatusKeikutsertaan     string  `form:"status_keikutsertaan" validate:"required"`
	KontrakKrs              bool    `form:"kontrak_krs" validate:"required"`
	JudulAktivitasMahasiswa string  `form:"judul_aktivitas_mahasiswa" validate:"required"`
	NoSkTugas               string  `form:"no_sk_tugas" validate:"required"`
	TanggalSkTugas          string  `form:"tanggal_sk_tugas" validate:"required"`
	JenisAnggota            string  `form:"jenis_anggota" validate:"required"`
	Ips                     float32 `form:"ips" validate:"required"`
	Ipk                     float32 `form:"ipk" validate:"required"`
	JumlahSks               int     `form:"jumlah_sks" validate:"required"`
	TotalSKS                int     `form:"total_sks" validate:"required"`
	BiayaKuliah             float32 `form:"biaya_kuliah" validate:"required"`
}

func (km *KampusMerdeka) MapRequest(idMahasiswa int, suratTugas string) (*model.KampusMerdeka, error) {
	tanggalSKTugas, err := util.ConvertStringToDate(km.TanggalSkTugas)
	if err != nil {
		return nil, errors.New("format tanggal salah")
	}
	return &model.KampusMerdeka{
		IdMahasiswa:             idMahasiswa,
		SuratTugas:              suratTugas,
		IdSemester:              km.IdSemester,
		IdDosenPembimbing:       km.IdDosenPembimbing,
		IdKategoriProgram:       km.IdKategoriProgram,
		StatusKeikutsertaan:     km.StatusKeikutsertaan,
		KontrakKrs:              km.KontrakKrs,
		JudulAktivitasMahasiswa: km.JudulAktivitasMahasiswa,
		NoSkTugas:               km.NoSkTugas,
		TanggalSkTugas:          tanggalSKTugas,
		JenisAnggota:            km.JenisAnggota,
		Ips:                     km.Ips,
		Ipk:                     km.Ipk,
		JumlahSks:               km.JumlahSks,
		TotalSks:                km.TotalSKS,
		BiayaKuliah:             km.BiayaKuliah,
	}, nil
}
