package filesystem

import (
	"fmt"
	"os"
)

// SegmentNext ...
func SegmentNext(directory string, segmentName string) (string, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return "", fmt.Errorf("failed to scan WAL directory: %w", err)
	}

	filenames := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filenames = append(filenames, file.Name())
	}

	idx := upperBound(filenames, segmentName)
	if idx <= len(filenames)-1 {
		return filenames[idx], nil
	}

	return "", nil
}

// SegmentLast ...
func SegmentLast(directory string) (string, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return "", fmt.Errorf("failed to scan WAL directory: %w", err)
	}

	filename := ""
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		if file.IsDir() {
			continue
		}

		filename = file.Name()
		break
	}

	return filename, nil
}

// CreateFile ...
func CreateFile(filename string) (*os.File, error) {
	flags := os.O_CREATE | os.O_WRONLY
	file, err := os.OpenFile(filename, flags, 0600) // nolint : TODO: 0644? G304: Potential file inclusion via variable
	if err != nil {
		return nil, err
	}

	return file, err
}

// WriteFile ...
func WriteFile(file *os.File, data []byte) (int, error) {
	writtenBytes, err := file.Write(data)
	if err != nil {
		return 0, err
	}

	if err = file.Sync(); err != nil {
		return 0, err
	}

	return writtenBytes, nil
}

func upperBound(array []string, target string) int {
	low, high := 0, len(array)-1

	for low <= high {
		mid := (low + high) / 2
		if array[mid] > target {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	return low
}
