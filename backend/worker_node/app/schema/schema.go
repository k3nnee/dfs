package schema

type Response struct {
	Message string
	Success bool
}

type FileUpload struct {
	FileName string
	File     []byte
	FileType string
}
