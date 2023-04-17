package env

import "os"

func GetPrestasiFolderId() string {
	return os.Getenv("GOOGLE_DRIVE_PRESTASI_FOLDER_ID")
}

func GetSuratTugasFolderId() string {
	return os.Getenv("GOOGLE_DRIVE_SURAT_TUGAS_FOLDER_ID")
}

func GetBeritaAcaraFolderId() string {
	return os.Getenv("GOOGLE_DRIVE_BERITA_ACARA_FOLDER_ID")
}
