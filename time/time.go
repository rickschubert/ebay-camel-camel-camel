package time

import "time"

type Minutes int64

type MilliSeconds int64

func (m Minutes) ToMs() MilliSeconds {
	return MilliSeconds(m * 60000)
}

func GetCurrentTime() MilliSeconds {
	return MilliSeconds(time.Now().UnixNano() / int64(time.Millisecond))
}
