package location_updater

type LocationEvent struct {
	Latitude  float32 `json:"lat,string"`
	Longitude float32 `json:"lng,string"`
	Id        string  `json:"id"`
	UnixTime  int     `json:"time,string"`
	TraceID   string  `json:"traceID"`
	SpanID    string  `json:"spanID"`
}
