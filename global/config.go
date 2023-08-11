package global

import (
	"os"

	"github.com/go-resty/resty/v2"
)

var (
	PersonalDocumentPath string
	HTTPClient           *resty.Client
	FundNameCodeMap      map[string]string
	FundCodeNameMap      map[string]string
)

func init() {
	PersonalDocumentPath = os.Getenv("PD")
	HTTPClient = resty.New()

	FundNameCodeMap = make(map[string]string)
	FundCodeNameMap = make(map[string]string)

	FundNameCodeMap["招商国证生物医药指数(LOF)A"] = "161726"
	FundNameCodeMap["招商中证白酒指数(LOF)A"] = "161725"
}
