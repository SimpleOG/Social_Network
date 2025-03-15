package CustomMediaDB

type MediaInterface interface {
}

type Media struct {
}

func NewMedia() MediaInterface {
	return &Media{}
}
