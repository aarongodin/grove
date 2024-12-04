package render

import (
	"github.com/aarongodin/grove/pkg/module"
	"github.com/rs/zerolog"
)

type sqlcRenderer struct {
	log zerolog.Logger
}

func (r sqlcRenderer) Render(m *module.Module, workDir string) ([]*Output, error) {
	return nil, nil
}

func NewSQLCRenderer(logger zerolog.Logger) Renderable {
	return sqlcRenderer{log: logger}
}
