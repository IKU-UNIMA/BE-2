package response

type (
	Dashboard struct {
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

	TotalDashboard struct {
		TotalKM       int `json:"total_kampus_merdeka"`
		TotalPrestasi int `json:"total_prestasi"`
	}
)
