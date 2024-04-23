package paste

import (
	"github.com/artshirshov/gastebin/internal/model/visibility"
)

type ResponseDto struct {
	Hash       string         `json:"hash"`
	Title      string         `json:"title"`
	Visibility visibility.Dto `json:"visibility"`
	Content    string         `json:"content"`
	CreatedAt  string         `json:"created_at"`
	ExpiredAt  string         `json:"expired_at"`
}

type ResponseDtos []*ResponseDto
