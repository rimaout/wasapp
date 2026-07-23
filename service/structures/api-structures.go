package structures

type ChatPrevieResponse struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	IsGroupChat bool   `json:"isGroupChat"`
}
