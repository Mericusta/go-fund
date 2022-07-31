package global

var (
	FundNameCodeMap map[string]string
	FundCodeNameMap map[string]string
)

func init() {
	FundNameCodeMap = make(map[string]string)
	FundCodeNameMap = make(map[string]string)

	FundNameCodeMap["招商国证生物医药指数(LOF)A"] = "161726"
	FundNameCodeMap["招商中证白酒指数(LOF)A"] = "161725"
}
