package xenon

import "context"

type FillOption = map[string]bool

type Entity struct {
	Ctx context.Context
}

type Repository struct {
	Ctx context.Context
}

type Service struct {
	Ctx context.Context
}