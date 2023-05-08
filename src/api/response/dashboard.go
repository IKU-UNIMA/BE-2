package response

type (
	Dashboard struct {
		ID     int    `json:"-"`
		Nama   string `json:"nama"`
		Jumlah int    `json:"jumlah"`
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
		}
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
