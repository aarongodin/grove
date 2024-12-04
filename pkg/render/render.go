package render

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/aarongodin/grove/pkg/module"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Renderable interface {
	Render(m *module.Module, workDir string) ([]*Output, error)
}

type Output struct {
	Path    string
	Content string
}

func Render(m *module.Module, debug bool, logger *zerolog.Logger) error {
	if logger == nil {
		logger = &log.Logger
	}

	log.Info().Msg("starting module render")

	workDir, err := os.MkdirTemp("", "grove")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	logger.Debug().Str("workDir", workDir).Msg("created temp dir")
	defer func() {
		if !debug {
			if err := os.RemoveAll(workDir); err != nil {
				logger.Fatal().Err(err).Msg("failed to remove temp dir")
			}
		} else {
			logger.Info().Msg("debug mode: ignoring cleanup of temp dir")
		}
	}()

	log.Info().Msgf("> go mod init %s", m.Name)
	goModInit := exec.Command("go", "mod", "init", m.Name)
	goModInit.Dir = workDir
	if err := goModInit.Run(); err != nil {
		return fmt.Errorf("failed to init go module: %w", err)
	}

	renderables := make([]Renderable, 0)
	if len(m.RelationalDatabases) != 0 {
		renderables = append(renderables, NewSQLCRenderer(*logger))
	}

	outs := make([]*Output, 0)
	for _, r := range renderables {
		result, err := r.Render(m, workDir)
		if err != nil {
			return err
		}
		outs = append(outs, result...)
	}

	// TODO: process all the outs to the file system

	if err := os.CopyFS(m.TargetDir, os.DirFS(workDir)); err != nil {
		return fmt.Errorf("failed to write file output to target dir: %w", err)
	}

	return nil
}
