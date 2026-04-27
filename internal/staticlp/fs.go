package staticlp

import (
	"os"
)

// WriteFileUTF8 writes text as UTF-8 without a byte-order mark (Hugo TOML loader rejects BOM).
func WriteFileUTF8(path string, content string, perm os.FileMode) error {
	return os.WriteFile(path, []byte(content), perm)
}
