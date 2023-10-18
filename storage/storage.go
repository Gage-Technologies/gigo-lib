package storage

import "io"

// Storage
// Interface for accessing remote object storage systems and local file systems
type Storage interface {
	// GetFile
	//
	//		Returns a file from the configured bucket.
	//      Returns nil if the file does not exist.
	//
	//		Args:
	//	       - path (string): The path of the file to retrieve.
	//
	//	 Returns:
	//	       - (io.ReadCloser): The contents of the file.
	GetFile(path string) (io.ReadCloser, error)

	// CreateFile
	//
	//	Creates a new file in the configured bucket.
	//
	//	Args:
	//	   - path (string): The path of the file to create.
	//	   - contents ([]byte): The contents of the file.
	CreateFile(path string, contents []byte) error

	// CreateFileStreamed
	//
	//	  Creates a new file in the configured bucket reading from an io.ReadCloser.
	//
	//	Args:
	//	      - path (string): The path of the file to create.
	//		  - length (int64): The size in bytes of the contents.
	//	      - contents (io.ReadCloser): The contents of the file.
	CreateFileStreamed(path string, length int64, contents io.ReadCloser) error

	// DeleteFile
	//
	//	    Deletes a file from the configured bucket.
	//
	//	Args:
	//	       - path (string): The path of the file to delete.
	DeleteFile(path string) error

	// MoveFile
	//
	//	    Moves a file within the configured bucket.
	//
	//	Args:
	//	       - src (string): The path of the file to move.
	//	       - dst (string): The new path of the file.
	MoveFile(src, dst string) error

	// CopyFile
	//
	//        Copies a file within the configured bucket.
	//
	//    Args:
	//           - src (string): The path of the file to copy.
	//           - dst (string): The new path of the file.
	CopyFile(src, dst string) error

	// MergeFiles
	//
	//        Merges multiple files within the configured bucket.
	//
	//    Args:
	//           - dst (string): The path of the merged file in the configured bucket.
	//           - paths ([]string): The paths of the files to merge in order of merge.
	//           - smallFiles (bool): Used for handling local merges of small files in some object stores
	//                                check the header doc for this function on each implementation to
	//                                determine if this parameter is used and what the implications of its
	//                                use are
	MergeFiles(dst string, paths []string, smallFiles bool) error

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
	Exists(path string) (bool, string, error)

	// CreateDir
	//
	//	    Creates a new directory in the configured bucket.
	//
	//	Args:
	//	       - path (string): The path of the directory to create.
	CreateDir(path string) error

	// ListDir
	//
	//	       Lists the contents of a directory in the configured bucket.
	//
	//	   Args:
	//	        - path (string): The path of the directory to list.
	//			- recursive (bool): Whether to list the directory recursively.
	//	   Returns:
	//          - []string: The list of files in the directory.
	ListDir(path string, recursive bool) ([]string, error)

	// DeleteDir
	//
	//	    Deletes a directory in the configured bucket.
	//
	//	Args:
	//	       - path (string): The path of the directory to delete.
	//		   - recursive (bool): Whether to delete all subdirectories within the passed directory
	DeleteDir(path string, recursive bool) error
}
