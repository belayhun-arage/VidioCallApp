package entity

type RTCTocken struct{
	User_Id       uint32 `json:"uid"`
	Channel_Name  string `json:"channel_name"`
	Role          uint32 `json:"role"`
}