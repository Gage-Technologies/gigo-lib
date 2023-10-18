package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileSystemStorage
//
//	Implementation of the Storage interface for filesystems
type FileSystemStorage struct {
	Storage
	root string
}

// CreateFileSystemStorage
//
//	Creates a new FileSystemStorage including creating the configured root path
//	if it doesn't already exist.
func CreateFileSystemStorage(root string) (*FileSystemStorage, error) {
	// create root directory if it doesn't exist
	if _, err := os.Stat(root); os.IsNotExist(err) {
		err := os.MkdirAll(root, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create root directory: %v", err)
		}
	}

	return &FileSystemStorage{
		root: root,
	}, nil
}

// GetFile
//
//			Returns a file from the configured bucket.
//	     Returns nil if the file does not exist.
//
//			Args:
//		       - path (string): The path of the file to retrieve.
//
//		 Returns:
//		       - (io.ReadCloser): The contents of the file.
func (s *FileSystemStorage) GetFile(path string) (io.ReadCloser, error) {
	// open the file
	file, err := os.Open(filepath.Join(s.root, path))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to open file path: %v", err)
	}

	return file, nil
}

// CreateFile
//
//	Creates a new file in the configured bucket.
//
//	Args:
//	   - path (string): The path of the file to create.
//	   - contents ([]byte): The contents of the file.
func (s *FileSystemStorage) CreateFile(path string, contents []byte) error {
	// get parent directory of file
	parent := filepath.Dir(filepath.Join(s.root, path))

	// create parent directory if it doesn't exist
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		err := os.MkdirAll(parent, 0755)
		if err != nil {
			return fmt.Errorf("failed to create file path directory: %v", err)
		}
	}

	// create the file
	file, err := os.Create(filepath.Join(s.root, path))
	if err != nil {
		return fmt.Errorf("failed to create file path: %v", err)
	}
	defer file.Close()

	// write the file
	_, err = file.Write(contents)
	if err != nil {
		return fmt.Errorf("failed to write file contents: %v", err)
	}

	return nil
}

// CreateFileStreamed
//
//	  Creates a new file in the configured bucket reading from an io.ReadCloser.
//
//	Args:
//	      - path (string): The path of the file to create.
//		  - length (int64): The size in bytes of the contents.
//	      - contents (io.ReadCloser): The contents of the file.
func (s *FileSystemStorage) CreateFileStreamed(path string, length int64, contents io.ReadCloser) error {
	// get parent directory of file
	parent := filepath.Dir(filepath.Join(s.root, path))

	// create parent directory if it doesn't exist
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		err := os.MkdirAll(parent, 0755)
		if err != nil {
			return fmt.Errorf("failed to create file path directory: %v", err)
		}
	}

	// create the file
	file, err := os.Create(filepath.Join(s.root, path))
	if err != nil {
		return fmt.Errorf("failed to create file path: %v", err)
	}
	defer file.Close()

	// write the file
	_, err = io.CopyN(file, contents, length)
	if err != nil {
		return fmt.Errorf("failed to write file contents: %v", err)
	}

	return nil
}

// DeleteFile
//
//	    Deletes a file from the configured bucket.
//
//	Args:
//	       - path (string): The path of the file to delete.
func (s *FileSystemStorage) DeleteFile(path string) error {
	// delete the file
	err := os.Remove(filepath.Join(s.root, path))
	if err != nil {
		return fmt.Errorf("failed to delete file path: %v", err)
	}
	return nil
}

// MoveFile
//
//	    Moves a file within the configured bucket.
//
//	Args:
//	       - src (string): The path of the file to move.
//	       - dst (string): The new path of the file.
func (s *FileSystemStorage) MoveFile(src, dst string) error {
	// move the file
	err := os.Rename(filepath.Join(s.root, src), filepath.Join(s.root, dst))
	if err != nil {
		return fmt.Errorf("failed to move file path: %v", err)
	}
	return nil
}

// CopyFile
//
//	    Copies a file within the configured bucket.
//
//	Args:
//	       - src (string): The path of the file to copy.
//	       - dst (string): The new path of the file.
func (s *FileSystemStorage) CopyFile(src, dst string) error {
	// stat file path
	info, err := os.Stat(filepath.Join(s.root, src))
	if err != nil {
		return fmt.Errorf("failed to stat file path: %v", err)
	}

	// ensure file is a copyable file
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	// open the source file
	srcFile, err := os.Open(filepath.Join(s.root, src))
	if err != nil {
		return fmt.Errorf("failed to open file path: %v", err)
	}
	defer srcFile.Close()

	// open the destination file
	dstFile, err := os.Create(filepath.Join(s.root, dst))
	if err != nil {
		return fmt.Errorf("failed to create file path: %v", err)
	}
	defer dstFile.Close()

	// copy the contents
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %v", err)
	}

	return nil
}

// MergeFiles
//
//	    Merges multiple files within the configured bucket.
//
//	Args:
//	       - dst (string): The path of the merged file in the configured bucket.
//	       - paths ([]string): The paths of the files to merge in order of merge.
//	       - smallFiles (bool): This parameter is a no-op in this implementation and only used
//	                            for compatibility with the Storage interface
func (s *FileSystemStorage) MergeFiles(dst string, paths []string, smallFiles bool) error {
	// get parent directory of file
	parent := filepath.Dir(filepath.Join(s.root, dst))

	// create parent directory if it doesn't exist
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		err := os.MkdirAll(parent, 0755)
		if err != nil {
			return fmt.Errorf("failed to create file path directory: %v", err)
		}
	}

	// open the destination file
	dstFile, err := os.Create(filepath.Join(s.root, dst))
	if err != nil {
		return fmt.Errorf("failed to create file path: %v", err)
	}

	defer dstFile.Close()

	// iterate through paths opening the source files and copying them into the destination file
	for _, path := range paths {
		// open source file
		srcFile, err := os.Open(filepath.Join(s.root, path))
		if err != nil {
			return fmt.Errorf("failed to open file path: %v", err)
		}

		// copy source file to destination file
		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			_ = srcFile.Close()
			return fmt.Errorf("failed to copy file contents: %v", err)
		}

		// close source file
		_ = srcFile.Close()
	}

	return nil
}

// Exists
//
//	   Checks whether the path exists in the configured bucket
//	   and returns what type of path it is (file, directory, symlink, etc.).
//
//	Args:
//	    - path (string): The path of the file to check.
//
//	Returns:
//	    - (bool): Whether the path exists or not.
//	    - (string): Path type
func (s *FileSystemStorage) Exists(path string) (bool, string, error) {
	// check if the file exists
	stat, err := os.Lstat(filepath.Join(s.root, path))
	if err != nil {
		// return false if the error is for a non-existent file
		if os.IsNotExist(err) {
			return false, "", nil
		}
		return false, "", fmt.Errorf("failed to stat file path: %v", err)
	}

	// default path type to file
	pathType := "file"

	// conditionally set path type to dir if the path is a directory
	if stat.IsDir() {
		pathType = "dir"
	}

	// conditionally set path type to symlink if the path is a symlink
	if stat.Mode()&os.ModeSymlink == os.ModeSymlink {
		pathType = "symlink"
	}

	return true, pathType, nil
}

// CreateDir
//
//	    Creates a new directory in the configured bucket.
//
//	Args:
//	       - path (string): The path of the directory to create.
func (s *FileSystemStorage) CreateDir(path string) error {
	// create the directory
	err := os.MkdirAll(filepath.Join(s.root, path), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	return nil
}

// ListDir
//
//		       Lists the contents of a directory in the configured bucket.
//
//		   Args:
//		        - path (string): The path of the directory to list.
//				- recursive (bool): Whether to list the directory recursively.
//		   Returns:
//	         - []string: The list of files in the directory.
func (s *FileSystemStorage) ListDir(path string, recursive bool) ([]string, error) {
	// list the directory
	files, err := os.ReadDir(filepath.Join(s.root, path))
	if err != nil {
		// return empty slice for non-existent directory to keep consistent
		// behavior between object storage on filesystem backends
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to list directory: %v", err)
	}

	// create slice to hold filepaths
	filepaths := make([]string, 0)

	// iterate over the files
	for _, file := range files {
		// handle normal files
		if !file.IsDir() {
			// append filepath to slice
			filepaths = append(filepaths, filepath.Join(path, file.Name()))
			continue
		}

		// handle non-recursive directory by adding directory name to slice with a '/' for suffix
		if !recursive {
			// add the directory to the slice
			filepaths = append(filepaths, filepath.Join(path, file.Name())+"/")
			continue
		}

		// recurse directory appending the resulting files to the slice
		dirContents, err := s.ListDir(filepath.Join(path, file.Name()), true)
		if err != nil {
			return nil, fmt.Errorf("failed to list directory contents: %v", err)
		}

		// append the directory contents to the slice
		filepaths = append(filepaths, dirContents...)
	}

	return filepaths, nil
}

// DeleteDir
//
//	    Deletes a directory in the configured bucket.
//
//	Args:
//	       - path (string): The path of the directory to delete.
//		   - recursive (bool): Whether to delete all subdirectories within the passed directory
func (s *FileSystemStorage) DeleteDir(path string, recursive bool) error {
	// conditionally delete the directory recursively
	if recursive {
		err := os.RemoveAll(filepath.Join(s.root, path))
		if err != nil {
			return fmt.Errorf("failed to delete directory: %v", err)
		}
		return nil
	}

	// list the directory
	files, err := os.ReadDir(filepath.Join(s.root, path))
	if err != nil {
		return fmt.Errorf("failed to list directory: %v", err)
	}

	// create boolean to track if we should remove the directory after removing the files
	removeDir := true

	// iterate paths in the directory removing only files
	for _, f := range files {
		// skip directories
		if f.IsDir() {
			// set removeDir false since there are subdirectories that
			// we want to leave in place
			removeDir = false
			continue
		}

		// remove the file
		err = os.Remove(filepath.Join(s.root, path, f.Name()))
		if err != nil {
			return fmt.Errorf("failed to delete file: %v", err)
		}
	}

	// conditionally remove top level directory
	if removeDir {
		err = os.RemoveAll(filepath.Join(s.root, path))
		if err != nil {
			return fmt.Errorf("failed to delete directory: %v", err)
		}
	}

	return nil
}
