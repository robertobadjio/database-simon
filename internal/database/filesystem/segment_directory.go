package filesystem

import (
	"fmt"
	"os"
)

// SegmentsDirectory ...
type SegmentsDirectory struct {
	directory string
}

// NewSegmentsDirectory ...
func NewSegmentsDirectory(directory string) *SegmentsDirectory {
	return &SegmentsDirectory{
		directory: directory,
	}
}

// ForEach ...
func (d *SegmentsDirectory) ForEach(action func([]byte) error) error {
	files, err := os.ReadDir(d.directory)
	if err != nil {
		// TODO: need to create a directory if it is missing
		return fmt.Errorf("failed to scan directory with segments: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := fmt.Sprintf("%s/%s", d.directory, file.Name())
		var errReadFile error
		data, errReadFile := os.ReadFile(filename) // nolint : TODO: G304: Potential file inclusion via variable
		if errReadFile != nil {
			return errReadFile
		}

		if err = action(data); err != nil {
			return err
		}
	}

	return nil
}
