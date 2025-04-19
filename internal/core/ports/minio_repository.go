package ports

type MinioRepository interface {
	UploadBase64Image(base64Data string, fileName string) (string, error)
	DeleteImage(fileName string) error
}
