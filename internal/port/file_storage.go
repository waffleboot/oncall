package port

type FileStorage interface {
	UploadFile(src string) (string, error)
	DownloadFile(src, dst string) error
}
