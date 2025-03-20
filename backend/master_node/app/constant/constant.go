package constant

const MaxFileSize = 10 << 20
const BlockSize = 5 << 20

var AllowedFileTypes = map[string]bool{
	"application/pdf": true,
	"image/jpeg":      true,
	"image/png":       true,
}
