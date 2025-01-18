package repo

import (
	"context"
	"fmt"
)

type Repository interface {
	SavePoints(context.Context, string, int)
	GetPoints(context.Context, string) (int, error)
}

type Repo struct {
	pointsDB map[string]int
}

func NewRepository() Repository {
	return &Repo{
		pointsDB: make(map[string]int),
	}
}

func (r *Repo) SavePoints(ctx context.Context, receiptId string, points int) {
	r.pointsDB[receiptId] = points
}

func (r *Repo) GetPoints(ctx context.Context, receiptId string) (int, error) {
	points, ok := r.pointsDB[receiptId]
	if !ok {
		return -1, fmt.Errorf("No receipt found for that ID.")
	}
	return points, nil
}
