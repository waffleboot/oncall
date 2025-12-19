package port

type FileStorage interface {
	UploadFile(src string) (string, error)
	DownloadFile(fileID, dst string) error
}
