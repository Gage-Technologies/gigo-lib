package agentsdk

import (
	"context"
	"fmt"
	"github.com/gage-technologies/gigo-lib/db/models"
	"golang.org/x/xerrors"
	"io"
	"net/http"
	"os"
)

func (c *AgentClient) WorkspaceInitializationStepCompleted(ctx context.Context, state models.WorkspaceInitState) error {
	stateReq := PostWorkspaceInitStateCompleted{State: state}
	res, err := c.Request(ctx, http.MethodPost, "/internal/v1/ws/init-state", stateReq)
	if err != nil {
		return xerrors.Errorf("init state post request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return readBodyAsError(res)
	}

	return nil
}

func (c *AgentClient) WorkspaceInitializationFailure(ctx context.Context, req PostWorkspaceInitFailure) error {
	res, err := c.Request(ctx, http.MethodPost, "/internal/v1/ws/init-failure", req)
	if err != nil {
		return xerrors.Errorf("agent init failure post request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return readBodyAsError(res)
	}

	return nil
}

func (c *AgentClient) WorkspaceGetExtension(ctx context.Context, extPath string) error {
	res, err := c.Request(ctx, http.MethodGet, "/internal/v1/ws/ext", nil)
	if err != nil {
		return xerrors.Errorf("agent get extension request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return readBodyAsError(res)
	}

	// write body to file
	f, err := os.Create(extPath)
	if err != nil {
		return xerrors.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(f, res.Body)
	if err != nil {
		return xerrors.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (c *AgentClient) WorkspaceGetCtExtension(ctx context.Context, extPath string) error {
	res, err := c.Request(ctx, http.MethodGet, "/internal/v1/ws/ext/ct", nil)
	if err != nil {
		return xerrors.Errorf("agent get extension request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return readBodyAsError(res)
	}

	// write body to file
	f, err := os.Create(extPath)
	if err != nil {
		return xerrors.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(f, res.Body)
	if err != nil {
		return xerrors.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (c *AgentClient) WorkspaceGetThemeExtension(ctx context.Context, extPath string) error {
	res, err := c.Request(ctx, http.MethodGet, "/internal/v1/ws/ext/theme", nil)
	if err != nil {
		return xerrors.Errorf("agent get extension request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return readBodyAsError(res)
	}

	// write body to file
	f, err := os.Create(extPath)
	if err != nil {
		return xerrors.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(f, res.Body)
	if err != nil {
		return xerrors.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (c *AgentClient) WorkspaceGetHolidayThemeExtension(ctx context.Context, extPath string, holiday int) error {
	res, err := c.Request(ctx, http.MethodGet, fmt.Sprintf("/internal/v1/ws/ext/holiday-theme?holiday=%d", int(holiday)), nil)
	if err != nil {
		return xerrors.Errorf("agent get extension request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return readBodyAsError(res)
	}

	// write body to file
	f, err := os.Create(extPath)
	if err != nil {
		return xerrors.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(f, res.Body)
	if err != nil {
		return xerrors.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (c *AgentClient) WorkspaceGetOpenVsxExtension(ctx context.Context, extId, version, vscVersion, extPath string) error {
	// override version if it is empty
	if version == "" {
		version = "latest"
	}

	// format url
	url := fmt.Sprintf("/internal/v1/ws/ext/open-vsx-cache?ext=%s&version=%s", extId, version)
	if vscVersion != "" {
		url = fmt.Sprintf("%s&vscVersion=%s", url, vscVersion)
	}

	res, err := c.Request(ctx, http.MethodGet, url, nil)
	if err != nil {
		return xerrors.Errorf("agent get extension request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return readBodyAsError(res)
	}

	// write body to file
	f, err := os.Create(extPath)
	if err != nil {
		return xerrors.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(f, res.Body)
	if err != nil {
		return xerrors.Errorf("failed to write file: %w", err)
	}

	return nil
}
