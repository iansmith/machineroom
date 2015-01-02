package shared

//
// THIS IS SHARED BY THE CLIENT AND SERVER
//
type ApiPayload struct {
	Username string
	Password string

	//it's a bad idea to use anything that has marshalling issues to/from
	//json.  sticking to string, int64, float64 and bool is advised.
}
