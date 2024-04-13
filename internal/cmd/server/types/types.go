package types

type ShutdownReq struct {
	Msg string
}

type HelloReq struct {
	Content string `json:"content"`
}

type ShutdownRsp struct {
	Reason string
}
