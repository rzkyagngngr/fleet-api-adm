package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type warehouseDetail struct {
	ID                 int       `json:"id"`
	WarehouseID        int       `json:"warehouse_id"`
	WarehouseCodeD     string    `json:"warehouse_code_d"`
	WarehouseNameD     string    `json:"warehouse_name_d"`
	WarehouseDType     string    `json:"werehouse_d_type"`
	WarehouseCapacityM int       `json:"warehouse_capacity_d_m3"`
	WarehouseCapacityT int       `json:"warehouse_capacity_d_ton"`
	XPosition          int       `json:"x_position"`
	YPosition          int       `json:"y_position"`
	WSize              int       `json:"w_size"`
	HSize              int       `json:"h_size"`
	Status             int       `json:"status"`
	CreationDate       time.Time `json:"creation_date"`
	CreationBy         string    `json:"creation_by"`
	LastUpdatedDate    time.Time `json:"last_updated_date"`
	LastUpdatedBy      string    `json:"last_updated_by"`
}

func registerStaticHealthRoutes(r *gin.Engine) {
	r.GET("/health/tiny", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "ok",
			"data": gin.H{
				"id": 1,
			},
		})
	})

	r.GET("/health/low", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "customers retrieved successfully",
			"data": []gin.H{
				{
					"id":                         13,
					"branch_code":                1001,
					"branch_name":                "Surabaya",
					"terminal_code":              101,
					"terminal_name":              "SURABAYA",
					"customer_code":              "CUST000013",
					"customer_name":              "PT. BABEH SEBLAK UPDATE LAGI",
					"customer_type":              "INTERNAL",
					"profit_center":              "2610",
					"customer_country":           "Indonesia",
					"customer_address":           "jaktim",
					"city":                       "jakarta",
					"contact_person":             "yaya",
					"phone_number":               "085876541192",
					"email_address":              "BABEH@ilcs.co.id",
					"fax_number":                 "13221",
					"tax_id_number":              "123992",
					"tax_id_16_digit":            "111111111111111",
					"tax_branch_code":            "6178",
					"national_id_number":         "001",
					"business_license_date":      "2000-10-20T00:00:00Z",
					"tax_id_document_upload":     "-",
					"registered_taxpayer_name":   "BABEH",
					"registered_taxpayer_address": "BABELAN",
					"business_type":              "logistik",
					"business_entity_type":       "PT",
					"bank_code":                  "BCA",
					"bank_account_idr":           "8080808081",
					"foreign_currency_account":   "IDR",
					"program_name":               "Master Customer",
					"status":                     1,
					"internal_notes":             "PLENGER",
					"creation_date":              "2026-05-06T11:00:59.792114Z",
					"creation_by":                "rizky.nugroho@mbg.co.id",
					"last_updated_date":          "2026-05-06T15:51:25.11242Z",
					"last_updated_by":            "ryan.hasbie@mbg.co.id",
				},
				{
					"id":                         15,
					"branch_code":                1001,
					"branch_name":                "",
					"terminal_code":              101,
					"terminal_name":              "SURABAYA",
					"customer_code":              "CUST000015",
					"customer_name":              "PT BERSINAR BERSAMA BABEH",
					"customer_type":              "INTERNAL",
					"profit_center":              "PC001",
					"customer_country":           "Indonesia",
					"customer_address":           "BEKASEA",
					"city":                       "BEKASEA",
					"contact_person":             "0123123123123",
					"phone_number":               "012301231023",
					"email_address":              "babeh@gmail.com",
					"fax_number":                 "0123123123",
					"tax_id_number":              "12312312",
					"tax_id_16_digit":            "1231231",
					"tax_branch_code":            "13123",
					"national_id_number":         "09331231231",
					"business_license_date":      "2026-02-07T00:00:00Z",
					"tax_id_document_upload":     "-",
					"registered_taxpayer_name":   "12312",
					"registered_taxpayer_address": "sndgnsldngl",
					"business_type":              "GENERAL CARGO",
					"business_entity_type":       "PT",
					"bank_code":                  "MANDIRI",
					"bank_account_idr":           "123123123123",
					"foreign_currency_account":   "IDR",
					"program_name":               "Master Customer",
					"status":                     1,
					"internal_notes":             nil,
					"creation_date":              "2026-05-07T14:02:29.351506Z",
					"creation_by":                "rizky.nugroho@mbg.co.id",
					"last_updated_date":          "2026-05-07T14:02:29.351506Z",
					"last_updated_by":            "rizky.nugroho@mbg.co.id",
				},
				{
					"id":                         10,
					"branch_code":                1001,
					"branch_name":                "Surabaya",
					"terminal_code":              101,
					"terminal_name":              "SURABAYA",
					"customer_code":              "CUST000010",
					"customer_name":              "PT. RYNEXUS UPDATE NEW",
					"customer_type":              "CORPORATE",
					"profit_center":              "11012701",
					"customer_country":           "Indonesia",
					"customer_address":           "Karawang",
					"city":                       "Karawang",
					"contact_person":             "Ryan Hasbie",
					"phone_number":               "08123456",
					"email_address":              "rynexus@email.com",
					"fax_number":                 "41357",
					"tax_id_number":              "123456789",
					"tax_id_16_digit":            "123456789",
					"tax_branch_code":            "KRW01",
					"national_id_number":         "123456789",
					"business_license_date":      "2027-10-10T00:00:00Z",
					"tax_id_document_upload":     "-",
					"registered_taxpayer_name":   "PT. RYNEXUS",
					"registered_taxpayer_address": "KARAWANG",
					"business_type":              "IT Solutions",
					"business_entity_type":       "PT",
					"bank_code":                  "BCA",
					"bank_account_idr":           "447",
					"foreign_currency_account":   "IDR",
					"program_name":               "Master Customer",
					"status":                     1,
					"internal_notes":             nil,
					"creation_date":              "2026-04-21T22:43:20.000282Z",
					"creation_by":                "ryan.hasbie@mbg.co.id",
					"last_updated_date":          "2026-05-06T14:15:23.532492Z",
					"last_updated_by":            "ryan.hasbie@mbg.co.id",
				},
			},
			"meta": gin.H{
				"page":               1,
				"limit":              10,
				"total_items":        3,
				"total_pages":        1,
				"max_download_limit": 1000,
			},
		})
	})

	r.GET("/health/high", func(c *gin.Context) {
		details := make([]warehouseDetail, 0, 100)
		createdAt, _ := time.Parse(time.RFC3339Nano, "2026-05-11T19:02:03.745846Z")
		for i := 1; i <= 100; i++ {
			row := (i - 1) / 10
			col := (i - 1) % 10
			details = append(details, warehouseDetail{
				ID:                 i + 17,
				WarehouseID:        1,
				WarehouseCodeD:     fmt.Sprintf("WH-001-%03d", i),
				WarehouseNameD:     fmt.Sprintf("Zone %03d", i),
				WarehouseDType:     "AREA",
				WarehouseCapacityM: 1000 + (i * 10),
				WarehouseCapacityT: 500 + (i * 5),
				XPosition:          10 + (col * 120),
				YPosition:          20 + (row * 80),
				WSize:              100,
				HSize:              60,
				Status:             1,
				CreationDate:       createdAt,
				CreationBy:         "ryan.hasbie@mbg.co.id",
				LastUpdatedDate:    createdAt,
				LastUpdatedBy:      "ryan.hasbie@mbg.co.id",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "warehouse detail retrieved successfully",
			"data": gin.H{
				"id":                1,
				"branch_code":       1001,
				"branch_name":       "Surabaya",
				"terminal_code":     101,
				"terminal_name":     "SURABAYA",
				"warehouse_code":    "WH-001",
				"warehouse_name":    "Warehouse Utama",
				"warehouse_type":    "Dry Storage",
				"warehouse_capacity": 5000,
				"status":            1,
				"creation_date":     "2026-04-19T04:39:56.606768Z",
				"creation_by":       "ryan.hasbie@mbg.co.id",
				"last_updated_date": "2026-05-11T19:02:03.745846Z",
				"last_updated_by":   "ryan.hasbie@mbg.co.id",
				"details":           details,
			},
		})
	})
}
