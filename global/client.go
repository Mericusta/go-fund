package global

import (
	"github.com/go-resty/resty/v2"
)

var (
	HTTPClient *resty.Client = resty.New()
)
