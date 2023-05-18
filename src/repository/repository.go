package repository

import (
	"context"

	"github.com/Edilberto-Vazquez/weahter-services/src/models"
)

type StationRepository interface {
	GetRecords(query models.FindRecords, ctx context.Context) ([]map[string]interface{}, error)
}
