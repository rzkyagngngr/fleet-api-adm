package helper

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DownloadQuery struct {
	IsDownload bool   `json:"is_download"`
	Type       string `json:"type"`
	RangeStart int    `json:"range_start"`
	RangeEnd   int    `json:"range_end"`
}

type SortQuery struct {
	By    string `json:"by"`
	Order string `json:"order"`
}

type PaginationQuery struct {
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
	Search   string            `json:"search"`
	Filters  map[string]string `json:"filters"`
	Sort     SortQuery         `json:"sort"`
	Download DownloadQuery     `json:"download"`
}

type PaginationMeta struct {
	Page             int   `json:"page"`
	Limit            int   `json:"limit"`
	TotalItems       int64 `json:"total_items"`
	TotalPages       int   `json:"total_pages"`
	MaxDownloadLimit int   `json:"max_download_limit"`
}

type MetaResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type NativePaginationConfig struct {
	TableName         string
	SelectColumns     []string
	SearchColumns     []string
	FilterableColumns map[string]string
	SortableColumns   map[string]string
	DefaultSortBy     string
	DefaultSortOrder  string
	MaxLimit          int
	MaxDownloadLimit  int
}

func MetaSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}, meta interface{}) {
	c.JSON(statusCode, MetaResponse{Success: true, Message: message, Data: data, Meta: meta})
}

func GetDynamicPaginatedNativeData(db *gorm.DB, config NativePaginationConfig, param PaginationQuery, result interface{}) (PaginationMeta, error) {
	meta := PaginationMeta{
		Page:             1,
		Limit:            10,
		MaxDownloadLimit: config.MaxDownloadLimit,
	}

	if len(config.SelectColumns) == 0 {
		return meta, errors.New("select columns are required")
	}

	if config.TableName == "" {
		return meta, errors.New("table name is required")
	}

	if err := validateSlicePointer(result); err != nil {
		return meta, err
	}

	page := maxInt(param.Page, 1)
	limit := maxInt(param.Limit, 1)

	if config.MaxLimit > 0 && limit > config.MaxLimit {
		limit = config.MaxLimit
	}

	baseQuery := fmt.Sprintf(" FROM %s WHERE 1=1", config.TableName)
	args := make([]interface{}, 0)
	searchValue := strings.TrimSpace(param.Search)

	if searchValue != "" && len(config.SearchColumns) > 0 {
		searchParts := make([]string, 0, len(config.SearchColumns))
		for _, column := range config.SearchColumns {
			searchParts = append(searchParts, fmt.Sprintf("UPPER(CAST(%s AS TEXT)) LIKE UPPER(?)", column))
			args = append(args, "%"+searchValue+"%")
		}
		baseQuery += " AND (" + strings.Join(searchParts, " OR ") + ")"
	}

	for key, value := range param.Filters {
		column, ok := config.FilterableColumns[key]
		filterValue := strings.TrimSpace(value)
		if !ok || filterValue == "" {
			continue
		}

		baseQuery += fmt.Sprintf(" AND UPPER(CAST(%s AS TEXT)) LIKE UPPER(?)", column)
		args = append(args, "%"+filterValue+"%")
	}

	countQuery := "SELECT COUNT(1)" + baseQuery
	if err := db.Raw(countQuery, args...).Scan(&meta.TotalItems).Error; err != nil {
		return meta, err
	}

	sortColumn := config.SortableColumns[param.Sort.By]
	if sortColumn == "" {
		sortColumn = config.SortableColumns[config.DefaultSortBy]
	}
	if sortColumn == "" {
		sortColumn = config.DefaultSortBy
	}

	sortOrder := strings.ToUpper(strings.TrimSpace(param.Sort.Order))
	if sortOrder != "DESC" {
		sortOrder = "ASC"
	}
	if strings.EqualFold(config.DefaultSortOrder, "DESC") && strings.TrimSpace(param.Sort.Order) == "" {
		sortOrder = "DESC"
	}

	offset := (page - 1) * limit
	effectiveLimit := limit

	if param.Download.IsDownload {
		switch strings.ToLower(strings.TrimSpace(param.Download.Type)) {
		case "range":
			rangeStart := maxInt(param.Download.RangeStart, 1)
			rangeEnd := maxInt(param.Download.RangeEnd, rangeStart)
			offset = (rangeStart - 1) * limit
			effectiveLimit = (rangeEnd - rangeStart + 1) * limit
		default:
			offset = 0
			effectiveLimit = int(meta.TotalItems)
		}

		if config.MaxDownloadLimit > 0 && effectiveLimit > config.MaxDownloadLimit {
			effectiveLimit = config.MaxDownloadLimit
		}
	}

	selectQuery := "SELECT " + strings.Join(config.SelectColumns, ", ")
	dataQuery := fmt.Sprintf("%s%s ORDER BY %s %s", selectQuery, baseQuery, sortColumn, sortOrder)

	if effectiveLimit > 0 {
		dataQuery += " LIMIT ? OFFSET ?"
		args = append(args, effectiveLimit, offset)
	}

	if err := db.Raw(dataQuery, args...).Scan(result).Error; err != nil {
		return meta, err
	}

	meta.Page = page
	meta.Limit = limit
	if limit > 0 {
		meta.TotalPages = int(math.Ceil(float64(meta.TotalItems) / float64(limit)))
	}

	return meta, nil
}

func validateSlicePointer(value interface{}) error {
	t := reflect.TypeOf(value)
	if t == nil || t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Slice {
		return errors.New("result must be a pointer to slice")
	}

	return nil
}

func maxInt(value int, fallback int) int {
	if value < fallback {
		return fallback
	}

	return value
}

func NewPaginationMeta(totalItems int64, page int, limit int) PaginationMeta {
	totalPages := 0
	if limit > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(limit)))
	}
	return PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
