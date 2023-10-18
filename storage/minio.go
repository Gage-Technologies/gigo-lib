package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gage-technologies/gigo-lib/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"os"
	"strings"
)

const MinioNotExistsError = "The specified key does not exist."

// MinioObjectStorage
//
//	Implementation of the Storage interface for Minio and all S3 compliant
//	object storage systems.
type MinioObjectStorage struct {
	Storage
	client *minio.Client
	config config.StorageS3Config
}

// CreateMinioObjectStorage
//
//	Creates a new MinioObjectStorage including initialization of the Minio client
//	and creating the configured bucket if it doesn't already exist.
func CreateMinioObjectStorage(config config.StorageS3Config) (*MinioObjectStorage, error) {
	// create minio client options
	opts := &minio.Options{
		Secure: config.UseSSL,
	}

	// conditionally add access credentials
	if config.AccessKey != "" && config.SecretKey != "" {
		opts.Creds = credentials.NewStaticV4(config.AccessKey, config.SecretKey, "")
	}

	// conditionally add region to minio client options
	if config.Region != "" {
		opts.Region = config.Region
	}

	// create minio client
	client, err := minio.New(config.Endpoint, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	// check if the bucket exists
	exists, err := client.BucketExists(context.TODO(), config.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %v", err)
	}

	// create bucket if it doesn't exist
	if !exists {
		err = client.MakeBucket(context.TODO(), config.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %v", err)
		}
	}

	return &MinioObjectStorage{
		client: client,
		config: config,
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
func (s *MinioObjectStorage) GetFile(path string) (io.ReadCloser, error) {
	exists, _, err := s.Exists(path)
	if err != nil {
		return nil, fmt.Errorf("check if object exists: %v", err)
	}
	if !exists {
		return nil, nil
	}
	file, err := s.client.GetObject(context.TODO(), s.config.Bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve object: %v", err)
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
func (s *MinioObjectStorage) CreateFile(path string, contents []byte) error {
	_, err := s.client.PutObject(context.TODO(), s.config.Bucket, path, bytes.NewReader(contents), int64(len(contents)), minio.PutObjectOptions{})
	return err
}

// CreateFileStreamed
//
//	  Creates a new file in the configured bucket reading from an io.ReadCloser.
//
//	Args:
//	      - path (string): The path of the file to create.
//		  - length (int64): The size in bytes of the contents.
//	      - contents (io.ReadCloser): The contents of the file.
func (s *MinioObjectStorage) CreateFileStreamed(path string, length int64, contents io.ReadCloser) error {
	_, err := s.client.PutObject(context.TODO(), s.config.Bucket, path, contents, length, minio.PutObjectOptions{})
	return err
}

// DeleteFile
//
//	    Deletes a file from the configured bucket.
//
//	Args:
//	       - path (string): The path of the file to delete.
func (s *MinioObjectStorage) DeleteFile(path string) error {
	return s.client.RemoveObject(context.TODO(), s.config.Bucket, path, minio.RemoveObjectOptions{})
}

// MoveFile
//
//	    Moves a file within the configured bucket.
//
//	Args:
//	       - src (string): The path of the file to move.
//	       - dst (string): The new path of the file.
func (s *MinioObjectStorage) MoveFile(src, dst string) error {
	// copy file to destination
	_, err := s.client.CopyObject(
		context.TODO(),
		minio.CopyDestOptions{
			Bucket: s.config.Bucket,
			Object: dst,
		},
		minio.CopySrcOptions{
			Bucket: s.config.Bucket,
			Object: src,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	// delete source file
	err = s.client.RemoveObject(context.TODO(), s.config.Bucket, src, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
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
func (s *MinioObjectStorage) CopyFile(src, dst string) error {
	// copy file to destination
	_, err := s.client.CopyObject(
		context.TODO(),
		minio.CopyDestOptions{
			Bucket: s.config.Bucket,
			Object: dst,
		},
		minio.CopySrcOptions{
			Bucket: s.config.Bucket,
			Object: src,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	return nil
}

// MergeFiles
//
//		   Merges multiple files within the configured bucket.
//
//	    NOTE:
//	    The S3 api does not permit the server-side composition (merging) of
//	    files smaller than 5MB excluding the final file. This function has
//	    the necessary logic to manually merge the files in the case that
//	    any of the files (excluding the final) are less than 5MB. However,
//	    manually merging files requires that the files be downloaded to the
//	    local file system, merged into single file locally, and re-uploaded
//	    as the final merged file.
//
//	    IT IS THE RESPONSIBILITY OF THE CALLER TO ENSURE SUFFICIENT SPACE AND
//		   BANDWIDTH ARE AVAILABLE TO PERFORM A LOCAL MERGE OF THE FILES
//
//		   Args:
//	           - dst (string): The path of the merged file in the configured bucket.
//	           - paths ([]string): The paths of the files to merge in order of merge.
//	           - smallFiles (bool): Whether to manually merge small files if they cannot be merged on the server (see note above)
func (s *MinioObjectStorage) MergeFiles(dst string, paths []string, smallFiles bool) error {
	// format paths into a slice minio.CopySrcOptions using the config bucket
	var src []minio.CopySrcOptions
	for _, path := range paths {
		src = append(src, minio.CopySrcOptions{
			Bucket: s.config.Bucket,
			Object: path,
		})
	}

	// use compose api to merge the files on the server side
	_, err := s.client.ComposeObject(
		context.TODO(),
		minio.CopyDestOptions{
			Bucket: s.config.Bucket,
			Object: dst,
		},
		src...,
	)
	// return for successful merge
	if err == nil {
		return nil
	}

	// return for unexpected errors
	if err != nil && !(smallFiles && strings.Contains(err.Error(), "is too small")) {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	// merge small files locally

	// create local temporary file to merge into
	tmpFile, err := os.CreateTemp(os.TempDir(), "gigo-local-file-merge")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	// create int64 to track the total size of the merged file
	fileLength := int64(0)

	// copy small files to local temporary file
	for _, path := range paths {
		// get file from storage
		file, err := s.GetFile(path)
		if err != nil {
			return fmt.Errorf("failed to get part file: %v", err)
		}

		// copy file to local temporary file
		n, err := io.Copy(tmpFile, file)
		if err != nil {
			_ = file.Close()
			return fmt.Errorf("failed to copy part into temp file: %v", err)
		}

		// increment file length by written bytes
		fileLength += n

		// close file
		_ = file.Close()
	}

	// reset temporary file cursor to the beginning of the file for upload
	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to reset temporary file cursor: %v", err)
	}

	// upload local temporary file
	err = s.CreateFileStreamed(dst, fileLength, tmpFile)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	return nil
}

// Exists
//
//		   Checks whether the path exists in the configured bucket
//		   and returns what type of path it is (file, directory, symlink, etc.).
//
//		   NOTE:
//		   Since object storage doesn't have a concept of directories this
//	    function will only check if a key (file) exists. If a valid
//		   prefix (directory) is passed the function will return false.
//		Args:
//		    - path (string): The path of the file to check.
//
//		Returns:
//		    - (bool): Whether the path exists or not.
//		    - (string): Path type
func (s *MinioObjectStorage) Exists(path string) (bool, string, error) {
	// stat the requested path to determine if it exists
	_, err := s.client.StatObject(context.TODO(), s.config.Bucket, path, minio.StatObjectOptions{})

	// handle failed call caused by an error that is not for a non-existent file
	if err != nil && err.Error() != MinioNotExistsError {
		return false, "", fmt.Errorf("failed to stat file: %v", err)
	}

	// handle successful call since we know that this is a valid file
	if err == nil {
		return true, "file", nil
	}

	return false, "", nil
}

// CreateDir
//
//	    Creates a new directory in the configured bucket.
//
//		NOTE:
//		Since object storage doesn't have the concept of directories
//		only prefixes, this function is a no-op. Creating a "directory"
//		in object storage is as simple as creating a key (file) at the
//		desired directory and the entire prefix (subdirectories) will
//		be automatically created.
//	Args:
//	       - path (string): The path of the directory to create.
func (s *MinioObjectStorage) CreateDir(path string) error {
	// object storage doesn't have directories only prefixes
	// so this functions is simply a pass-through for compliance
	// with the Storage interface
	return nil
}

// ListDir
//
//		       Lists the contents of a directory in the configured bucket.
//
//		   Args:
//		        - path (string): The path of the directory to list.
//		   		- recursive (bool): Whether to list the directory recursively.
//		   Returns:
//	         - []string: The list of files in the directory.
func (s *MinioObjectStorage) ListDir(path string, recursive bool) ([]string, error) {
	// create cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// conditionally append final slash to path if it was not passed
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	// call list api on client to get channel for iteration of the directory
	objects := s.client.ListObjects(ctx, s.config.Bucket, minio.ListObjectsOptions{Prefix: path, Recursive: recursive})

	// iterate through the objects
	var contents []string
	for object := range objects {
		// handle error for object
		if object.Err != nil {
			return nil, fmt.Errorf("failed to list directory: %v", object.Err)
		}
		contents = append(contents, object.Key)
	}

	return contents, nil
}

// DeleteDir
//
//	    Deletes a directory in the configured bucket.
//
//	Args:
//	       - path (string): The path of the directory to delete.
//		   - recursive (bool): Whether to delete all subdirectories within the passed directory
func (s *MinioObjectStorage) DeleteDir(path string, recursive bool) error {
	// conditionally append final slash to path if it was not passed
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	// call list api on client to get channel for iteration of the directory
	objects := s.client.ListObjects(context.TODO(), s.config.Bucket, minio.ListObjectsOptions{Prefix: path, Recursive: recursive})

	// create channel to pass removals to the removal function
	removeChannel := make(chan minio.ObjectInfo)

	// create context with cancel for this operation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// begin removal function via the minio agentsdk
	errChan := s.client.RemoveObjects(ctx, s.config.Bucket, removeChannel, minio.RemoveObjectsOptions{})

	// iterate through the objects sending them to the removal channel in a go routine
	go func() {
		// close the channel
		defer close(removeChannel)

		for object := range objects {
			// exit if context is done
			select {
			case <-ctx.Done():
				return
			default:
			}

			// skip directories
			if strings.HasSuffix(object.Key, "/") {
				continue
			}

			// send object to removal channel
			removeChannel <- object
		}
	}()

	// iterate error channel handling error
	for deletionErr := range errChan {
		cancel()
		return fmt.Errorf("failed to delete directory: %v\n    object: %v\n    version: %v", deletionErr.Err, deletionErr.ObjectName, deletionErr.VersionID)
	}

	return nil
}
