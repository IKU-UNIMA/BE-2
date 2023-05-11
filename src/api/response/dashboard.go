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

	DetailDashboard struct {
		Prodi struct {
			ID        int    `json:"id"`
			KodeProdi string `json:"kode_prodi"`
			Nama      string `json:"nama"`
			Jenjang   string `json:"jenjang"`
		} `json:"prodi"`
		Fakultas struct {
			ID   int    `json:"id"`
			Nama string `json:"nama"`
		} `json:"fakultas"`
		Semester struct {
			Id   int    `json:"id"`
			Nama string `json:"nama"`
		} `json:"semester"`
		Jumlah int `json:"jumlah"`
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
