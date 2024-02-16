package storage

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestCreateFileSystemStorage(t *testing.T) {
	storage, err := CreateFileSystemStorage("/tmp/gigo-fs-test")
	if err != nil {
		t.Fatalf("\nCreateFileSystemStorage failed\n    Error: %v", err)
	}

	if storage == nil {
		t.Fatalf("\nCreateFileSystemStorage failed\n    Error: storage was returned nil")
	}

	t.Log("\nCreateFileSystemStorage succeeded")
}

func TestFileSystemStorage_GetFile(t *testing.T) {
	storage, err := CreateFileSystemStorage("/tmp/gigo-fs-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_GetFile failed\n    Error: %v", err)
	}

	defer os.Remove("/tmp/gigo-fs-test/get-test")

	f, err := os.Create("/tmp/gigo-fs-test/get-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_GetFile failed\n    Error: %v", err)
	}

	defer f.Close()

	_, err = f.Write([]byte("get-file-test-contents"))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_GetFile failed\n    Error: %v", err)
	}

	_ = f.Close()

	file, _, err := storage.GetFile("get-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_GetFile failed\n    Error: %v", err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("\nFileSystemStorage_GetFile failed\n    Error: %v", err)
	}

	if string(data) != "get-file-test-contents" {
		t.Fatalf("\nFileSystemStorage_GetFile failed\n    Error: file corrupted")
	}

	file, _, err = storage.GetFile("get-test-no-exist")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_GetFile failed\n    Error: %v", err)
	}

	if file != nil {
		t.Fatalf("\nFileSystemStorage_GetFile failed\n    Error: file was not nil")
	}

	t.Log("\nFileSystemStorage_GetFile succeeded")
}

func TestFileSystemStorage_CreateFile(t *testing.T) {
	storage, err := CreateFileSystemStorage("/tmp/gigo-fs-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_CreateFile failed\n    Error: %v", err)
	}

	defer os.Remove("/tmp/gigo-fs-test/create-test")

	err = storage.CreateFile("create-test", []byte("create-test-contents"))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_CreateFile failed\n    Error: %v", err)
	}

	data, err := os.ReadFile("/tmp/gigo-fs-test/create-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_CreateFile failed\n    Error: %v", err)
	}

	if string(data) != "create-test-contents" {
		t.Fatalf("\nFileSystemStorage_CreateFile failed\n    Error: file corrupted")
	}

	t.Log("\nFileSystemStorage_CreateFile succeeded")
}

func TestFileSystemStorage_CreateFileStreamed(t *testing.T) {
	storage, err := CreateFileSystemStorage("/tmp/gigo-fs-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_CreateFileStreamed failed\n    Error: %v", err)
	}

	defer os.Remove("/tmp/gigo-fs-test/create-test-streamed")

	buf := []byte("create-streamed-contents")
	err = storage.CreateFileStreamed("create-test-streamed", int64(len(buf)), io.NopCloser(bytes.NewBuffer(buf)))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_CreateFileStreamed failed\n    Error: %v", err)
	}

	data, err := os.ReadFile("/tmp/gigo-fs-test/create-test-streamed")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_CreateFileStreamed failed\n    Error: %v", err)
	}

	if string(data) != string(buf) {
		t.Fatalf("\nFileSystemStorage_CreateFileStreamed failed\n    Error: file corrupted")
	}

	t.Log("\nFileSystemStorage_CreateFileStreamed succeeded")
}

func TestFileSystemStorage_DeleteFile(t *testing.T) {
	storage, err := CreateFileSystemStorage("/tmp/gigo-fs-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_DeleteFile failed\n    Error: %v", err)
	}

	defer os.Remove("/tmp/gigo-fs-test/delete-test")

	err = storage.CreateFile("delete-test", []byte("delete-test-contents"))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_DeleteFile failed\n    Error: %v", err)
	}

	err = storage.DeleteFile("delete-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_DeleteFile failed\n    Error: %v", err)
	}

	_, err = os.Stat("/tmp/gigo-fs-test/delete-test")
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("\nFileSystemStorage_DeleteFile failed\n    Error: %v", err)
	}

	if err == nil {
		t.Fatalf("\nFileSystemStorage_DeleteFile failed\n    Error: file was not deleted")
	}

	t.Log("\nFileSystemStorage_DeleteFile succeeded")
}

func TestFileSystemStorage_MoveFile(t *testing.T) {
	storage, err := CreateFileSystemStorage("/tmp/gigo-fs-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_MoveFile failed\n    Error: %v", err)
	}

	defer os.Remove("/tmp/gigo-fs-test/move-test-src")
	defer os.Remove("/tmp/gigo-fs-test/move-test-dst")

	err = storage.CreateFile("move-test-src", []byte("move-test-contents"))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_MoveFile failed\n    Error: %v", err)
	}

	err = storage.MoveFile("move-test-src", "move-test-dst")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_MoveFile failed\n    Error: %v", err)
	}

	data, err := os.ReadFile("/tmp/gigo-fs-test/move-test-dst")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_MoveFile failed\n    Error: %v", err)
	}

	if string(data) != "move-test-contents" {
		t.Fatalf("\nFileSystemStorage_MoveFile failed\n    Error: file corrupted")
	}

	_, err = os.Stat("/tmp/gigo-fs-test/move-test-src")
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("\nFileSystemStorage_MoveFile failed\n    Error: %v", err)
	}

	if err == nil {
		t.Fatalf("\nFileSystemStorage_MoveFile failed\n    Error: file was not deleted")
	}

	t.Log("\nFileSystemStorage_MoveFile succeeded")
}

func TestFileSystemStorage_CopyFile(t *testing.T) {
	storage, err := CreateFileSystemStorage("/tmp/gigo-fs-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_CopyFile failed\n    Error: %v", err)
	}

	defer os.Remove("/tmp/gigo-fs-test/copy-test-src")
	defer os.Remove("/tmp/gigo-fs-test/copy-test-dst")

	err = storage.CreateFile("copy-test-src", []byte("copy-test-contents"))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_CopyFile failed\n    Error: %v", err)
	}

	err = storage.CopyFile("copy-test-src", "copy-test-dst")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_CopyFile failed\n    Error: %v", err)
	}

	data, err := os.ReadFile("/tmp/gigo-fs-test/copy-test-dst")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_CopyFile failed\n    Error: %v", err)
	}

	if string(data) != "copy-test-contents" {
		t.Fatalf("\nFileSystemStorage_CopyFile failed\n    Error: file corrupted")
	}

	t.Log("\nFileSystemStorage_CopyFile succeeded")
}

func TestFileSystemStorage_MergeFiles(t *testing.T) {
	storage, err := CreateFileSystemStorage("/tmp/gigo-fs-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_MergeFiles failed\n    Error: %v", err)
	}

	defer os.Remove("/tmp/gigo-fs-test/merge-test-1")
	defer os.Remove("/tmp/gigo-fs-test/merge-test-2")
	defer os.Remove("/tmp/gigo-fs-test/merge-test-3")
	defer os.Remove("/tmp/gigo-fs-test/merge-test-dst")

	err = storage.CreateFile("merge-test-1", []byte("this is a test"))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_MergeFiles failed\n    Error: %v", err)
	}
	err = storage.CreateFile("merge-test-2", []byte(" for merging multiple "))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_MergeFiles failed\n    Error: %v", err)
	}
	err = storage.CreateFile("merge-test-3", []byte("files together"))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_MergeFiles failed\n    Error: %v", err)
	}

	err = storage.MergeFiles("merge-test-dst", []string{"merge-test-1", "merge-test-2", "merge-test-3"}, false)
	if err != nil {
		t.Fatalf("\nFileSystemStorage_MergeFiles failed\n    Error: %v", err)
	}

	data, err := os.ReadFile("/tmp/gigo-fs-test/merge-test-dst")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_MergeFiles failed\n    Error: %v", err)
	}

	if string(data) != "this is a test for merging multiple files together" {
		t.Fatalf("\nFileSystemStorage_MergeFiles failed\n    Error: file corrupted")
	}

	t.Log("\nFileSystemStorage_MergeFiles succeeded")
}

func TestFileSystemStorage_Exists(t *testing.T) {
	storage, err := CreateFileSystemStorage("/tmp/gigo-fs-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_Exists failed\n    Error: %v", err)
	}

	defer os.Remove("/tmp/gigo-fs-test/exists-test")
	defer os.RemoveAll("/tmp/gigo-fs-test/exists-test-dir")
	defer os.Remove("/tmp/gigo-fs-test/exists-test-sym")

	err = storage.CreateFile("exists-test", []byte("exists-test-contents"))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_Exists failed\n    Error: %v", err)
	}

	err = os.Mkdir("/tmp/gigo-fs-test/exists-test-dir", 0755)
	if err != nil {
		t.Fatalf("\nFileSystemStorage_Exists failed\n    Error: %v", err)
	}

	err = os.Symlink("/tmp/gigo-fs-test/exists-test", "/tmp/gigo-fs-test/exists-test-sym")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_Exists failed\n    Error: %v", err)
	}

	exists, pathType, err := storage.Exists("exists-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_Exists failed\n    Error: %v", err)
	}

	if !exists {
		t.Fatalf("\nFileSystemStorage_Exists failed\n    Error: file not found")
	}

	if pathType != "file" {
		fmt.Println(pathType)
		t.Fatalf("\nFileSystemStorage_Exists failed\n    Error: path type not file")
	}

	exists, pathType, err = storage.Exists("exists-test-dir")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_Exists failed\n    Error: %v", err)
	}

	if !exists {
		t.Fatalf("\nFileSystemStorage_Exists failed\n    Error: file not found")
	}

	if pathType != "dir" {
		t.Fatalf("\nFileSystemStorage_Exists failed\n    Error: path type not dir")
	}

	exists, pathType, err = storage.Exists("exists-test-sym")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_Exists failed\n    Error: %v", err)
	}

	if !exists {
		t.Fatalf("\nFileSystemStorage_Exists failed\n    Error: file not found")
	}

	if pathType != "symlink" {
		fmt.Println(pathType)
		t.Fatalf("\nFileSystemStorage_Exists failed\n    Error: path type not symlink")
	}

	exists, pathType, err = storage.Exists("exists-test-fake")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_Exists failed\n    Error: %v", err)
	}

	if exists {
		t.Fatalf("\nFileSystemStorage_Exists failed\n    Error: non-existent file found")
	}

	t.Log("\nFileSystemStorage_Exists succeeded")
}

func TestFileSystemStorage_CreateDir(t *testing.T) {
	storage, err := CreateFileSystemStorage("/tmp/gigo-fs-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_CreateDir failed\n    Error: %v", err)
	}

	defer os.RemoveAll("/tmp/gigo-fs-test/create-dir-test")

	err = storage.CreateDir("create-dir-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_CreateDir failed\n    Error: %v", err)
	}

	exists, pathType, err := storage.Exists("create-dir-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_CreateDir failed\n    Error: %v", err)
	}

	if !exists {
		t.Fatalf("\nFileSystemStorage_CreateDir failed\n    Error: dir not found")
	}

	if pathType != "dir" {
		t.Fatalf("\nFileSystemStorage_CreateDir failed\n    Error: path type not dir")
	}

	t.Log("\nFileSystemStorage_CreateDir succeeded")
}

func TestFileSystemStorage_ListDir(t *testing.T) {
	storage, err := CreateFileSystemStorage("/tmp/gigo-fs-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_ListDir failed\n    Error: %v", err)
	}

	defer os.RemoveAll("/tmp/gigo-fs-test/list-dir-test")

	err = storage.CreateFile("list-dir-test/test-1", []byte("list-dir-test-contents"))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_ListDir failed\n    Error: %v", err)
	}
	err = storage.CreateFile("list-dir-test/test-2", []byte("list-dir-test-contents"))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_ListDir failed\n    Error: %v", err)
	}
	err = storage.CreateFile("list-dir-test/nested/test-3", []byte("list-dir-test-contents"))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_ListDir failed\n    Error: %v", err)
	}
	err = storage.CreateFile("list-dir-test/nested/test-4", []byte("list-dir-test-contents"))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_ListDir failed\n    Error: %v", err)
	}

	files, err := storage.ListDir("list-dir-test", false)
	if err != nil {
		t.Fatalf("\nFileSystemStorage_ListDir failed\n    Error: %v", err)
	}

	if len(files) != 3 {
		t.Fatalf("\nFileSystemStorage_ListDir failed\n    Error: wrong number of files returned")
	}

	sort.Strings(files)

	if !reflect.DeepEqual(files, []string{"list-dir-test/nested/", "list-dir-test/test-1", "list-dir-test/test-2"}) {
		fmt.Println(files)
		t.Fatalf("\nFileSystemStorage_ListDir failed\n    Error: wrong files returned")
	}

	files, err = storage.ListDir("list-dir-test", true)
	if err != nil {
		t.Fatalf("\nFileSystemStorage_ListDir failed\n    Error: %v", err)
	}

	if len(files) != 4 {
		t.Fatalf("\nFileSystemStorage_ListDir failed\n    Error: wrong number of files returned")
	}

	sort.Strings(files)

	if !reflect.DeepEqual(files, []string{"list-dir-test/nested/test-3", "list-dir-test/nested/test-4", "list-dir-test/test-1", "list-dir-test/test-2"}) {
		fmt.Println(files)
		t.Fatalf("\nFileSystemStorage_ListDir failed\n    Error: wrong files returned")
	}

	t.Log("\nFileSystemStorage_ListDir succeeded")
}

func TestFileSystemStorage_DeleteDir(t *testing.T) {
	storage, err := CreateFileSystemStorage("/tmp/gigo-fs-test")
	if err != nil {
		t.Fatalf("\nFileSystemStorage_DeleteDir failed\n    Error: %v", err)
	}

	defer os.RemoveAll("/tmp/gigo-fs-test/del-test-dir")

	err = storage.CreateFile("del-test-dir/test-1", []byte(""))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_DeleteDir failed\n    Error: %v", err)
	}
	err = storage.CreateFile("del-test-dir/test-2", []byte(""))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_DeleteDir failed\n    Error: %v", err)
	}
	err = storage.CreateFile("del-test-dir/nested/test-3", []byte(""))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_DeleteDir failed\n    Error: %v", err)
	}
	err = storage.CreateFile("del-test-dir/nested/test-4", []byte(""))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_DeleteDir failed\n    Error: %v", err)
	}

	err = storage.DeleteDir("del-test-dir", false)
	if err != nil {
		t.Fatalf("\nFileSystemStorage_DeleteDir failed\n    Error: %v", err)
	}

	files, err := storage.ListDir("del-test-dir", true)
	if err != nil {
		t.Fatalf("\nFileSystemStorage_DeleteDir failed\n    Error: %v", err)
	}

	if len(files) != 2 {
		fmt.Println(files)
		t.Fatalf("\nFileSystemStorage_DeleteDir failed\n    Error: wrong number of files returned")
	}

	sort.Strings(files)

	if !reflect.DeepEqual(files, []string{"del-test-dir/nested/test-3", "del-test-dir/nested/test-4"}) {
		fmt.Println(files)
		t.Fatalf("\nFileSystemStorage_DeleteDir failed\n    Error: wrong files returned")
	}

	err = storage.CreateFile("del-test-dir/test-1", []byte(""))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_DeleteDir failed\n    Error: %v", err)
	}
	err = storage.CreateFile("del-test-dir/test-2", []byte(""))
	if err != nil {
		t.Fatalf("\nFileSystemStorage_DeleteDir failed\n    Error: %v", err)
	}

	err = storage.DeleteDir("del-test-dir", true)
	if err != nil {
		t.Fatalf("\nFileSystemStorage_DeleteDir failed\n    Error: %v", err)
	}

	files, err = storage.ListDir("del-test-dir", true)
	if err != nil {
		t.Fatalf("\nFileSystemStorage_DeleteDir failed\n    Error: %v", err)
	}

	if len(files) != 0 {
		t.Fatalf("\nFileSystemStorage_DeleteDir failed\n    Error: wrong number of files returned")
	}

	t.Log("\nFileSystemStorage_DeleteDir succeeded")
}
