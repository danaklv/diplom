package utils

import (
	"errors"
	"mime/multipart"
	"path/filepath"
)

// Разрешённые расширения
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

// Максимальный размер файла — 5 МБ
const maxFileSize = 5 * 1024 * 1024 // 5MB

// ValidateImage проверяет тип и размер загружаемого файла
func ValidateImage(file multipart.File, handler *multipart.FileHeader) error {
	ext := filepath.Ext(handler.Filename)
	if !allowedExtensions[ext] {
		return errors.New("invalid file type: only .jpg, .jpeg, .png, .gif allowed")
	}

	if handler.Size > maxFileSize {
		return errors.New("file too large: max 5MB allowed")
	}

	// Вернуть курсор в начало (иначе io.Copy не сработает)
	if _, err := file.Seek(0, 0); err != nil {
		return errors.New("failed to reset file pointer")
	}

	return nil
}
