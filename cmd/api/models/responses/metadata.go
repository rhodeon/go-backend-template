package responses

import "math"

type Metadata interface {
	metadata()
}

type Pagination struct {
	CurrentPage  uint `json:"current_page,omitempty" doc:"Current page number" example:"1"`
	PageSize     uint `json:"page_size,omitempty" doc:"Number of items returned" example:"10" maximum:"50"`
	TotalPages   uint `json:"total_pages,omitempty" doc:"Total number of matching pages"`
	TotalRecords uint `json:"total_records,omitempty" doc:"Total number of items across all matching pages"`
}

func (p *Pagination) metadata() {}

// CalculatePagination generates a metadata from the given page, limit
// and total number of records.
func CalculatePagination(page uint, pageSize uint, totalRecords uint) *Pagination {
	// Return an empty metadata if no records are found
	if totalRecords == 0 {
		return nil
	}

	return &Pagination{
		CurrentPage:  page,
		PageSize:     pageSize,
		TotalPages:   uint(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
