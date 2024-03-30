package logic

type StatusCode int32

const (
	STATUS_ON_SALE StatusCode = iota
	STATUS_SOLD_OUT
	STATUS_NOT_START
	STATUS_ENDED
)

func (s StatusCode) String() string {
	switch s {
	case STATUS_ON_SALE:
		return "on_sale"
	case STATUS_SOLD_OUT:
		return "sold_out"
	case STATUS_NOT_START:
		return "not_start"
	case STATUS_ENDED:
		return "ended"
	default:
		return "UNKNOWN"
	}
}
