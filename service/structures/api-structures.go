package structures

type ChatPreviewResponse struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	IsGroupChat bool   `json:"isGroupChat"`
}
