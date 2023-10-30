package now_time

import "time"

type NowType func() time.Time

var Get = func() time.Time {
	return time.Now()
}
