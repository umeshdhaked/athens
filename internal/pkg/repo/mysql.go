package repo

import (
	"context"
	"errors"

	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/utils"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/fastbiztech/hastinapura/pkg/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MysqlRepository struct {
	db *gorm.DB
}

// Create creates a new item in the database
func (m *MysqlRepository) Create(ctx context.Context, model models.IModel) error {
	query := m.db.Create(model)

	return GetError(ctx, query)
}

// Create creates a new item in the database
func (m *MysqlRepository) CreateInBatches(ctx context.Context, models interface{}, batch int) error {
	query := m.db.CreateInBatches(models, batch)

	return GetError(ctx, query)
}

// Update multiple attributes with `struct`, will only update those changed & non-blank fields
func (m *MysqlRepository) Update(ctx *gin.Context, model models.IModel) error {
	// This check is very important. If the model is empty, it'll end up updating the whole table.
	if utils.IsEmpty(model) {
		logger.GetLogger().Error("empty model for update")
		return errors.New("empty model for update")
	}

	query := m.db.Model(model).Updates(model)

	return GetError(ctx, query)
}

// Delete deletes the given model
func (m *MysqlRepository) Delete(ctx *gin.Context, model models.IModel) error {
	query := m.db.Delete(model)

	return GetError(ctx, query)
}

// FindByID returns model based by on primary key
func (m *MysqlRepository) FindByID(ctx *gin.Context, id interface{}, model models.IModel) error {
	query := m.db.First(model, "id = ?", id)

	return GetError(ctx, query)
}

// Find finds a model based on given conditions
func (m *MysqlRepository) Find(ctx *gin.Context, model models.IModel, condition map[string]interface{}) error {
	query := m.db.Where(condition).First(model)

	return GetError(ctx, query)
}

// FindMultiple will fetch multiple entities based on the condition
func (m *MysqlRepository) FindMultiple(ctx *gin.Context, models interface{}, condition map[string]interface{}) error {
	query := m.db.Where(condition).Find(models)

	return GetError(ctx, query)
}

// todo: implement pagination thing (offset)
func (m *MysqlRepository) FindMultiplePagination(ctx *gin.Context,
	models interface{},
	condition map[string]interface{},
	pagination dtos.Pagination) error {

	query := m.db

	// remove empty conditions
	for k, v := range condition {
		if utils.IsEmpty(v) {
			delete(condition, k)
		}
	}

	if !utils.IsEmpty(pagination.From) {
		query = query.Where("created_at >= ?", pagination.From)
	}

	if !utils.IsEmpty(pagination.To) {
		query = query.Where("created_at <= ?", pagination.To)
	}

	query = query.Find(models, condition)

	return GetError(ctx, query)
}
