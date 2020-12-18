package models

type VkIDModel struct {
	Response struct {
		Count int   `json:"count"`
		Items []int `json:"items"`
	} `json:"response"`
}

type VkUserModel struct {
	Response []struct {
		FirstName       string `json:"first_name"`
		ID              int    `json:"id"`
		LastName        string `json:"last_name"`
		CanAccessClosed bool   `json:"can_access_closed"`
		IsClosed        bool   `json:"is_closed"`
		Sex             int    `json:"sex"`
		Verified        int    `json:"verified"`
		Interests       string `json:"interests"`
	} `json:"response"`
}
