package shift

type ShiftInput struct {
	Type      string `json:"type" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}