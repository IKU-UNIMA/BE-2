package response

type (
	Dashboard struct {
		Target         string                       `json:"target"`
		Total          int                          `json:"total"`
		TotalMahasiswa int                          `json:"total_mahasiswa"`
		Pencapaian     string                       `json:"pencapaian"`
		Detail         []DashboardDetailPerFakultas `json:"detail"`
	}

	DashboardDetailPerFakultas struct {
		ID              int    `json:"id"`
		Fakultas        string `json:"fakultas"`
		JumlahMahasiswa int    `json:"jumlah_mahasiswa"`
		Jumlah          int    `json:"jumlah"`
		Persentase      string `json:"persentase"`
	}

	DashboardPerProdi struct {
		Fakultas       string                    `json:"fakultas"`
		Total          int                       `json:"total"`
		TotalMahasiswa int                       `json:"total_mahasiswa"`
		Pencapaian     string                    `json:"pencapaian"`
		Detail         []DashboardDetailPerProdi `json:"detail"`
	}

	DashboardDetailPerProdi struct {
		Prodi           string `json:"prodi"`
		JumlahMahasiswa int    `json:"jumlah_mahasiswa"`
		Jumlah          int    `json:"jumlah"`
		Persentase      string `json:"persentase"`
	}

	KategoriDashboard struct {
		Nama   string `json:"nama"`
		Jumlah int    `json:"jumlah"`
	}

	TotalDashboard struct {
		Nama  string `json:"nama"`
		Total int    `json:"total"`
	}

	DashboardUmum struct {
		Fakultas  int `json:"fakultas"`
		Prodi     int `json:"prodi"`
		Dosen     int `json:"dosen"`
		Mahasiswa int `json:"mahasiswa"`
	}
)
