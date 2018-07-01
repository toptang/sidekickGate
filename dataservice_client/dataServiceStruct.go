package dataservice_client

type DataRspCommon struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// login
type DataReqLogin struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// add user
type DataReqAddUser DataReqLogin

// delete user
type DataReqDelUser struct {
	Name string `json:"name"`
}

// update user password
type DataReqUpdateUserPassword struct {
	Name        string `json:"name"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
