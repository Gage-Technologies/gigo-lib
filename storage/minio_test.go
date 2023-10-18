package storage

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gage-technologies/gigo-lib/config"
	"github.com/minio/minio-go/v7"
	"golang.org/x/crypto/sha3"
	"io"
	"reflect"
	"sort"
	"testing"
)

// we duplicate the hash data func to prevent an import cycle
func hashData(data []byte) (string, error) {
	// create SHA256 hasher
	hasher := sha3.New256()

	// add data to hasher
	hasher.Write(data)

	// create output buffer
	buff := make([]byte, 32)
	// sum hash slices into buffer
	hasher.Sum(buff[:0])

	// hex encode hash and return
	return hex.EncodeToString(buff), nil
}

func TestCreateMinioObjectStorage(t *testing.T) {
	s, err := CreateMinioObjectStorage(config.StorageS3Config{
		Endpoint:  "localhost:9000",
		Bucket:    "gigo-test",
		AccessKey: "gigo-tests",
		SecretKey: "jDX2FrdsJpsfm64zJpy8uL7ADD7YO4bx",
	})
	if err != nil {
		t.Fatalf("\nCreateMinioObjectStorage failed\n    Error: %v", err)
	}

	if s == nil {
		t.Fatalf("\nCreateMinioObjectStorage failed\n    Error: storage was nil")
	}

	t.Log("\nCreateMinioObjectStorage succeeded")
}

func TestMinioObjectStorage_CreateFile(t *testing.T) {
	s, err := CreateMinioObjectStorage(config.StorageS3Config{
		Endpoint:  "localhost:9000",
		Bucket:    "gigo-test",
		AccessKey: "gigo-tests",
		SecretKey: "jDX2FrdsJpsfm64zJpy8uL7ADD7YO4bx",
	})
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_CreateFile failed\n    Error: %v", err)
	}

	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "create-test", minio.RemoveObjectOptions{})

	err = s.CreateFile("create-test", []byte("create-test-file"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_CreateFile failed\n    Error: %v", err)
	}

	obj, err := s.client.GetObject(context.TODO(), s.config.Bucket, "create-test", minio.GetObjectOptions{})
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_CreateFile failed\n    Error: %v", err)
	}

	data, err := io.ReadAll(obj)
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_CreateFile failed\n    Error: %v", err)
	}

	if string(data) != "create-test-file" {
		t.Fatalf("\nMinioObjectStorage_CreateFile failed\n    Error: corrupt file")
	}

	t.Log("\nMinioObjectStorage_CreateFile succeeded")
}

func TestMinioObjectStorage_GetFile(t *testing.T) {
	s, err := CreateMinioObjectStorage(config.StorageS3Config{
		Endpoint:  "localhost:9000",
		Bucket:    "gigo-test",
		AccessKey: "gigo-tests",
		SecretKey: "jDX2FrdsJpsfm64zJpy8uL7ADD7YO4bx",
	})
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_GetFile failed\n    Error: %v", err)
	}

	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "get-test", minio.RemoveObjectOptions{})

	err = s.CreateFile("get-test", []byte("get-test-file"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_GetFile failed\n    Error: %v", err)
	}

	obj, err := s.GetFile("get-test")
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_GetFile failed\n    Error: %v", err)
	}

	data, err := io.ReadAll(obj)
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_GetFile failed\n    Error: %v", err)
	}

	if string(data) != "get-test-file" {
		t.Fatalf("\nMinioObjectStorage_GetFile failed\n    Error: corrupt file")
	}

	obj, err = s.GetFile("get-test-no-exist")
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_GetFile failed\n    Error: %v", err)
	}

	if obj != nil {
		t.Fatalf("\nMinioObjectStorage_GetFile failed\n    Error: object was not nil")
	}

	t.Log("\nMinioObjectStorage_GetFile succeeded")
}

func TestMinioObjectStorage_CreateFileStreamed(t *testing.T) {
	s, err := CreateMinioObjectStorage(config.StorageS3Config{
		Endpoint:  "localhost:9000",
		Bucket:    "gigo-test",
		AccessKey: "gigo-tests",
		SecretKey: "jDX2FrdsJpsfm64zJpy8uL7ADD7YO4bx",
	})
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_CreateFileStreamed failed\n    Error: %v", err)
	}

	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "create-file-streamed", minio.RemoveObjectOptions{})

	buf := []byte("create-file-streamed-contents")
	err = s.CreateFileStreamed("create-file-streamed", int64(len(buf)), io.NopCloser(bytes.NewBuffer(buf)))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_CreateFileStreamed failed\n    Error: %v", err)
	}

	obj, err := s.client.GetObject(context.TODO(), s.config.Bucket, "create-file-streamed", minio.GetObjectOptions{})
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_CreateFileStreamed failed\n    Error: %v", err)
	}

	data, err := io.ReadAll(obj)
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_CreateFileStreamed failed\n    Error: %v", err)
	}

	if string(data) != string(buf) {
		t.Fatalf("\nMinioObjectStorage_CreateFileStreamed failed\n    Error: corrupt file")
	}

	t.Log("\nMinioObjectStorage_CreateFileStreamed succeeded")
}

func TestMinioObjectStorage_DeleteFile(t *testing.T) {
	s, err := CreateMinioObjectStorage(config.StorageS3Config{
		Endpoint:  "localhost:9000",
		Bucket:    "gigo-test",
		AccessKey: "gigo-tests",
		SecretKey: "jDX2FrdsJpsfm64zJpy8uL7ADD7YO4bx",
	})
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_DeleteFile failed\n    Error: %v", err)
	}

	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "delete-test", minio.RemoveObjectOptions{})

	err = s.CreateFile("delete-test", []byte("delete-test-contents"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_DeleteFile failed\n    Error: %v", err)
	}

	err = s.DeleteFile("delete-test")
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_DeleteFile failed\n    Error: %v", err)
	}

	_, err = s.client.StatObject(context.TODO(), s.config.Bucket, "delete-test", minio.StatObjectOptions{})
	if err != nil && err.Error() != MinioNotExistsError {
		t.Fatalf("\nMinioObjectStorage_DeleteFile failed\n    Error: %v", err)
	}

	if err == nil {
		t.Fatalf("\nMinioObjectStorage_DeleteFile failed\n    Error: file was not deleted")
	}

	t.Log("\nMinioObjectStorage_DeleteFile succeeded")
}

func TestMinioObjectStorage_MoveFile(t *testing.T) {
	s, err := CreateMinioObjectStorage(config.StorageS3Config{
		Endpoint:  "localhost:9000",
		Bucket:    "gigo-test",
		AccessKey: "gigo-tests",
		SecretKey: "jDX2FrdsJpsfm64zJpy8uL7ADD7YO4bx",
	})
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_MoveFile failed\n    Error: %v", err)
	}

	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "move-test-src", minio.RemoveObjectOptions{})
	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "move-test-dst", minio.RemoveObjectOptions{})

	err = s.CreateFile("move-test-src", []byte("move-test-contents"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_MoveFile failed\n    Error: %v", err)
	}

	err = s.MoveFile("move-test-src", "move-test-dst")
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_MoveFile failed\n    Error: %v", err)
	}

	_, err = s.client.StatObject(context.TODO(), s.config.Bucket, "move-test-src", minio.StatObjectOptions{})
	if err != nil && err.Error() != MinioNotExistsError {
		t.Fatalf("\nMinioObjectStorage_MoveFile failed\n    Error: %v", err)
	}

	if err == nil {
		t.Fatalf("\nMinioObjectStorage_MoveFile failed\n    Error: file was not moved")
	}

	obj, err := s.client.GetObject(context.TODO(), s.config.Bucket, "move-test-dst", minio.GetObjectOptions{})
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_MoveFile failed\n    Error: %v", err)
	}

	data, err := io.ReadAll(obj)
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_MoveFile failed\n    Error: %v", err)
	}

	if string(data) != "move-test-contents" {
		t.Fatalf("\nMinioObjectStorage_MoveFile failed\n    Error: file was corrupted")
	}

	t.Log("\nMinioObjectStorage_MoveFile succeeded")
}

func TestMinioObjectStorage_CopyFile(t *testing.T) {
	s, err := CreateMinioObjectStorage(config.StorageS3Config{
		Endpoint:  "localhost:9000",
		Bucket:    "gigo-test",
		AccessKey: "gigo-tests",
		SecretKey: "jDX2FrdsJpsfm64zJpy8uL7ADD7YO4bx",
	})
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_CopyFile failed\n    Error: %v", err)
	}

	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "copy-test-src", minio.RemoveObjectOptions{})
	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "copy-test-dst", minio.RemoveObjectOptions{})

	err = s.CreateFile("copy-test-src", []byte("copy-test-contents"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_CopyFile failed\n    Error: %v", err)
	}

	err = s.CopyFile("copy-test-src", "copy-test-dst")
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_CopyFile failed\n    Error: %v", err)
	}

	obj, err := s.client.GetObject(context.TODO(), s.config.Bucket, "copy-test-dst", minio.GetObjectOptions{})
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_CopyFile failed\n    Error: %v", err)
	}

	data, err := io.ReadAll(obj)
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_CopyFile failed\n    Error:")
	}

	if string(data) != "copy-test-contents" {
		t.Fatalf("\nMinioObjectStorage_CopyFile failed\n    Error: file was corrupted")
	}

	t.Log("\nMinioObjectStorage_CopyFile succeeded")
}

func TestMinioObjectStorage_MergeFiles(t *testing.T) {
	s, err := CreateMinioObjectStorage(config.StorageS3Config{
		Endpoint:  "localhost:9000",
		Bucket:    "gigo-test",
		AccessKey: "gigo-tests",
		SecretKey: "jDX2FrdsJpsfm64zJpy8uL7ADD7YO4bx",
	})
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_MergeFiles failed\n    Error: %v", err)
	}

	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "merge-test-1", minio.RemoveObjectOptions{})
	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "merge-test-2", minio.RemoveObjectOptions{})
	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "merge-test-3", minio.RemoveObjectOptions{})
	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "merge-test-dst", minio.RemoveObjectOptions{})
	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "merge-test-dst-large", minio.RemoveObjectOptions{})

	err = s.CreateFile("merge-test-1", []byte("this is"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_MergeFiles failed\n    Error: %v", err)
	}
	err = s.CreateFile("merge-test-2", []byte(" a test for merging "))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_MergeFiles failed\n    Error: %v", err)
	}
	err = s.CreateFile("merge-test-3", []byte("multiple files using the compose api\nfiles:\n1\n2\n3"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_MergeFiles failed\n    Error: %v", err)
	}

	err = s.MergeFiles("merge-test-dst", []string{"merge-test-1", "merge-test-2", "merge-test-3"}, true)
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_MergeFiles failed\n    Error: %v", err)
	}

	file, err := s.GetFile("merge-test-dst")
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_MergeFiles failed\n    Error: %v", err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_MergeFiles failed\n    Error:")
	}

	_ = file.Close()

	if string(data) != "this is a test for merging multiple files using the compose api\nfiles:\n1\n2\n3" {
		t.Fatalf("\nMinioObjectStorage_MergeFiles failed\n    Error: file was corrupted")
	}

	hasher := sha3.New256()
	paths := make([]string, 0)

	for i := 0; i < 5; i++ {
		buf := make([]byte, 1024*1024*5)
		if i == 4 {
			buf = make([]byte, 1024)
		}
		_, err = rand.Read(buf)
		if err != nil {
			t.Fatalf("\nMinioObjectStorage_MergeFiles failed\n    Error: %v", err)
		}

		hasher.Write(buf)

		path := fmt.Sprintf("merge-test-src-large-%d", i)

		defer s.client.RemoveObject(context.TODO(), s.config.Bucket, path, minio.RemoveObjectOptions{})

		err = s.CreateFile(path, buf)
		if err != nil {
			t.Fatalf("\nMinioObjectStorage_MergeFiles failed\n    Error: %v", err)
		}

		paths = append(paths, path)
	}

	buff := make([]byte, 32)
	hasher.Sum(buff[:0])
	controlHash := hex.EncodeToString(buff)

	err = s.MergeFiles("merge-test-dst-large", paths, false)
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_MergeFiles failed\n    Error: %v", err)
	}

	file, err = s.GetFile("merge-test-dst-large")
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_MergeFiles failed\n    Error: %v", err)
	}

	data, err = io.ReadAll(file)
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_MergeFiles failed\n    Error:")
	}

	resultHash, err := hashData(data)
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_MergeFiles failed\n    Error: %v", err)
	}

	if controlHash != resultHash {
		t.Fatalf("\nMinioObjectStorage_MergeFiles failed\n    Error: file corrupted")
	}

	t.Log("\nMinioObjectStorage_MergeFiles succeeded")
}

func TestMinioObjectStorage_Exists(t *testing.T) {
	s, err := CreateMinioObjectStorage(config.StorageS3Config{
		Endpoint:  "localhost:9000",
		Bucket:    "gigo-test",
		AccessKey: "gigo-tests",
		SecretKey: "jDX2FrdsJpsfm64zJpy8uL7ADD7YO4bx",
	})
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_Exists failed\n    Error: %v", err)
	}

	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "exists-test", minio.RemoveObjectOptions{})
	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "exists-test-dir/test", minio.RemoveObjectOptions{})

	err = s.CreateFile("exists-test", []byte("exists-test-contents"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_Exists failed\n    Error: %v", err)
	}

	exists, pathType, err := s.Exists("exists-test")
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_Exists failed\n    Error: %v", err)
	}

	if !exists {
		t.Fatalf("\nMinioObjectStorage_Exists failed\n    Error: file was not found")
	}

	if pathType != "file" {
		t.Fatalf("\nMinioObjectStorage_Exists failed\n    Error: incorrect path type")
	}

	exists, pathType, err = s.Exists("exists-test-no-exist")
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_Exists failed\n    Error: %v", err)
	}

	if exists {
		t.Fatalf("\nMinioObjectStorage_Exists failed\n    Error: non-existent file returned as existing")
	}

	t.Log("\nMinioObjectStorage_Exists succeeded")
}

func TestMinioObjectStorage_ListDir(t *testing.T) {
	s, err := CreateMinioObjectStorage(config.StorageS3Config{
		Endpoint:  "localhost:9000",
		Bucket:    "gigo-test",
		AccessKey: "gigo-tests",
		SecretKey: "jDX2FrdsJpsfm64zJpy8uL7ADD7YO4bx",
	})
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_ListDir failed\n    Error: %v", err)
	}

	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "dir/list-test-1", minio.RemoveObjectOptions{})
	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "dir/list-test-2", minio.RemoveObjectOptions{})
	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "dir/nested/list-test-3", minio.RemoveObjectOptions{})
	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "dir/nested/list-test-4", minio.RemoveObjectOptions{})

	err = s.CreateFile("dir/list-test-1", []byte("list-test-"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_ListDir failed\n    Error: %v", err)
	}
	err = s.CreateFile("dir/list-test-2", []byte("list-test-"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_ListDir failed\n    Error: %v", err)
	}
	err = s.CreateFile("dir/nested/list-test-3", []byte("list-test-"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_ListDir failed\n    Error: %v", err)
	}
	err = s.CreateFile("dir/nested/list-test-4", []byte("list-test-"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_ListDir failed\n    Error: %v", err)
	}

	contents, err := s.ListDir("dir", false)
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_ListDir failed\n    Error: %v", err)
	}

	if len(contents) != 3 {
		fmt.Println(contents)
		t.Fatalf("\nMinioObjectStorage_ListDir failed\n    Error: incorrect number of files returned")
	}

	sort.Strings(contents)

	if !reflect.DeepEqual(contents, []string{"dir/list-test-1", "dir/list-test-2", "dir/nested/"}) {
		t.Fatalf("\nMinioObjectStorage_ListDir failed\n    Error: incorrect contents returned")
	}

	contents, err = s.ListDir("dir", true)
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_ListDir failed\n    Error: %v", err)
	}

	if len(contents) != 4 {
		t.Fatalf("\nMinioObjectStorage_ListDir failed\n    Error: incorrect number of files returned")
	}

	sort.Strings(contents)

	if !reflect.DeepEqual(contents, []string{"dir/list-test-1", "dir/list-test-2", "dir/nested/list-test-3", "dir/nested/list-test-4"}) {
		t.Fatalf("\nMinioObjectStorage_ListDir failed\n    Error: incorrect contents returned")
	}

	t.Log("\nMinioObjectStorage_ListDir succeeded")
}

func TestMinioObjectStorage_DeleteDir(t *testing.T) {
	s, err := CreateMinioObjectStorage(config.StorageS3Config{
		Endpoint:  "localhost:9000",
		Bucket:    "gigo-test",
		AccessKey: "gigo-tests",
		SecretKey: "jDX2FrdsJpsfm64zJpy8uL7ADD7YO4bx",
	})
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_DeleteDir failed\n    Error: %v", err)
	}

	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "dir/del-test-1", minio.RemoveObjectOptions{})
	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "dir/del-test-2", minio.RemoveObjectOptions{})
	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "dir/nested/del-test-3", minio.RemoveObjectOptions{})
	defer s.client.RemoveObject(context.TODO(), s.config.Bucket, "dir/nested/del-test-4", minio.RemoveObjectOptions{})

	err = s.CreateFile("dir/del-test-1", []byte("del-test-"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_DeleteDir failed\n    Error: %v", err)
	}
	err = s.CreateFile("dir/del-test-2", []byte("del-test-"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_DeleteDir failed\n    Error: %v", err)
	}
	err = s.CreateFile("dir/nested/del-test-3", []byte("del-test-"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_DeleteDir failed\n    Error: %v", err)
	}
	err = s.CreateFile("dir/nested/del-test-4", []byte("del-test-"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_DeleteDir failed\n    Error: %v", err)
	}

	err = s.DeleteDir("dir", false)
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_DeleteDir failed\n    Error: %v", err)
	}

	contents, err := s.ListDir("dir", true)
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_DeleteDir failed\n    Error: %v", err)
	}

	if len(contents) != 2 {
		fmt.Println(contents)
		t.Fatalf("\nMinioObjectStorage_DeleteDir failed\n    Error: incorrect number of files returned")
	}

	sort.Strings(contents)

	if !reflect.DeepEqual(contents, []string{"dir/nested/del-test-3", "dir/nested/del-test-4"}) {
		t.Fatalf("\nMinioObjectStorage_DeleteDir failed\n    Error: incorrect contents returned")
	}

	err = s.CreateFile("dir/del-test-1", []byte("del-test-"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_DeleteDir failed\n    Error: %v", err)
	}
	err = s.CreateFile("dir/del-test-2", []byte("del-test-"))
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_DeleteDir failed\n    Error: %v", err)
	}

	err = s.DeleteDir("dir", true)
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_DeleteDir failed\n    Error: %v", err)
	}

	contents, err = s.ListDir("dir", true)
	if err != nil {
		t.Fatalf("\nMinioObjectStorage_DeleteDir failed\n    Error: %v", err)
	}

	if len(contents) != 0 {
		t.Fatalf("\nMinioObjectStorage_DeleteDir failed\n    Error: incorrect number of files returned")
	}

	t.Log("\nMinioObjectStorage_DeleteDir succeeded")
}
