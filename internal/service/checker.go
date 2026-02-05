package service

import (
	entities "URL_checker/internal/repo/dto"
	"context"
)

type Checker interface {
	Check(ctx context.Context, target entities.Targets) (entities.Checks, error)
}

type Checkerr struct {
	repo []entities.Checks
}

func (c *Checkerr) Check(ctx context.Context, target entities.Targets) (entities.Checks, error) {
	return entities.Checks{}, nil
}
