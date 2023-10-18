package git

import (
	"context"
	"errors"
	"fmt"
	"github.com/gage-technologies/gigo-lib/utils"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"path/filepath"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/server"
	"net/http"
)

func GitPush(directory string, requestContext context.Context, upr *packp.UploadPackRequest, w http.ResponseWriter) error {
	hashDir, err := utils.HashData([]byte(directory))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to hash git directory, err: %v", err))
	}

	ep, err := transport.NewEndpoint("/")
	if err != nil {
		return err
	}
	bfs := osfs.New(filepath.Join(hashDir[:3], hashDir[3:6], hashDir+".git"))
	ld := server.NewFilesystemLoader(bfs)
	svr := server.NewServer(ld)
	sess, err := svr.NewUploadPackSession(ep, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create upload pack session, err: %v", err))
	}
	res, err := sess.UploadPack(requestContext, upr)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to upload session pack, err: %v", err))
	}
	err = res.Encode(w)
	if err != nil {
		return err
	}
	return nil
}

func GitClone(directory string, requestContext context.Context, upr *packp.ReferenceUpdateRequest, w http.ResponseWriter) error {
	hashDir, err := utils.HashData([]byte(directory))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to hash git directory, err: %v", err))
	}

	ep, err := transport.NewEndpoint("/")
	if err != nil {
		return errors.New(fmt.Sprintf("failed to connect to endpoint, err: %v", err))
	}
	bfs := osfs.New(filepath.Join(hashDir[:3], hashDir[3:6], hashDir+".git"))
	ld := server.NewFilesystemLoader(bfs)
	svr := server.NewServer(ld)
	sess, err := svr.NewReceivePackSession(ep, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to receive pack session, err: %v", err))
	}
	res, err := sess.ReceivePack(requestContext, upr)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to receive pack, err: %v", err))
	}

	err = res.Encode(w)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to encode, err: %v", err))
	}

	return nil
}

func GitPull(directory string, requestContext context.Context, upr *packp.ReferenceUpdateRequest, w http.ResponseWriter) error {
	hashDir, err := utils.HashData([]byte(directory))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to hash git directory, err: %v", err))
	}

	ep, err := transport.NewEndpoint("/")
	if err != nil {
		return err
	}
	bfs := osfs.New(filepath.Join(hashDir[:3], hashDir[3:6], hashDir+".git"))
	ld := server.NewFilesystemLoader(bfs)
	svr := server.NewServer(ld)
	sess, err := svr.NewReceivePackSession(ep, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create upload pack session, err: %v", err))
	}
	res, err := sess.ReceivePack(requestContext, upr)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to upload session pack, err: %v", err))
	}

	err = res.Encode(w)
	if err != nil {
		return err
	}
	return nil
}

func GitFork(directory string, requestContext context.Context, upr *packp.ReferenceUpdateRequest, w http.ResponseWriter) error {
	hashDir, err := utils.HashData([]byte(directory))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to hash git directory, err: %v", err))
	}

	ep, err := transport.NewEndpoint("/")
	if err != nil {
		return err
	}
	bfs := osfs.New(filepath.Join(hashDir[:3], hashDir[3:6], hashDir+".git"))
	ld := server.NewFilesystemLoader(bfs)
	svr := server.NewServer(ld)
	sess, err := svr.NewReceivePackSession(ep, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create upload pack session, err: %v", err))
	}
	res, err := sess.ReceivePack(requestContext, upr)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to upload session pack, err: %v", err))
	}

	err = res.Encode(w)
	if err != nil {
		return err
	}
	return nil
}
