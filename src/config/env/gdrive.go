package env

import "os"

func GetPrestasiFolderId() string {
	return os.Getenv("GOOGLE_DRIVE_PRESTASI_FOLDER_ID")
}
