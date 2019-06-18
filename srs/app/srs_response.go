package app

type SrsResponse struct {
	StreamId int
}

func NewSrsResponse(sid int) *SrsResponse {
	return &SrsResponse{
		StreamId:sid,
	}
}


