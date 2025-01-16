package v1

type CreateUserEventRequest struct {
	UserID      string `json:"userID" valid:"required,stringlength(6|30)"`
	EventType   uint   `json:"eventType" valid:"required"`
	EventDetail string `json:"eventDetail" valid:"required,stringlength(2|1000)"`
	Duration    int64  `json:"duration"`
	Referer     string `json:"referer"`
	Platform    string `json:"platform"  valid:"required,stringlength(6|30)"`
}
