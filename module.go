package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/input"
	"github.com/rs/zerolog/log"
)

type Module struct {
	TargetDir string
	Name      string
}

func (m *Module) PromptTargetDir() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get cwd: %w", err)
	}
	result, err := prompt.New().Ask("Target directory").Input(
		cwd,
		input.WithHelp(true),
		input.WithValidateFunc(func(i string) error {
			path, err := filepath.Abs(i)
			if err != nil {
				return err
			}
			s, err := os.Stat(path)
			if os.IsNotExist(err) {
				return nil
			} else if err != nil {
				return fmt.Errorf("failed to read directory: %w", err)
			}
			if !s.IsDir() {
				return errors.New("Provided path is not a directory")
			}
			return nil
		}),
	)
	if err != nil {
		if errors.Is(err, prompt.ErrUserQuit) {
			log.Fatal().Msg("goodbye")
		}
		return err
	}
	m.TargetDir, err = filepath.Abs(result)
	if err != nil {
		return err
	}
	return nil
}

func (m *Module) PromptName() error {
	result, err := prompt.New().Ask("Module name (e.g. github.com/...)").Input(
		"",
		input.WithHelp(true),
		input.WithValidateFunc(func(i string) error {
			if i == "" {
				return errors.New("Module name is required")
			}
			return nil
		}),
	)
	if err != nil {
		if errors.Is(err, prompt.ErrUserQuit) {
			log.Fatal().Msg("goodbye")
		}
		return err
	}
	m.Name = result
	return nil
}

func (m *Module) Render() error {
	log.Info().Msg("starting module render")

	// Create the directory if it doesn't exist
	_, err := os.Stat(m.TargetDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(m.TargetDir, 0766); err != nil {
			return err
		}
	} else if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	log.Info().Msgf("> go mod init %s", m.Name)
	goModInit := exec.Command("go", "mod", "init", m.Name)
	goModInit.Dir = m.TargetDir
	if err := goModInit.Run(); err != nil {
		return fmt.Errorf("failed to init go module: %w", err)
	}

	return nil
}
