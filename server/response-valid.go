package server

//// responseValid request validity data type
//type responseValid struct {
//	valid bool
//	cause string
//}
//
//// setCause sets the cause in response
//func (responseValid *responseValid) setCause(cause string) {
//	responseValid.cause = cause
//}
//
//// setValid sets the validity of the response
//func (responseValid *responseValid) setValid(flag bool) {
//	responseValid.valid = flag
//}
//
//// exportJSON returns the validity of the response in JSON format
//func (responseValid *responseValid) exportJSON() []byte {
//	dataJSON, err := json.Marshal(struct {
//		Valid bool   `json:"valid"`
//		Cause string `json:"cause,omitempty"`
//	}{
//		Valid: responseValid.valid,
//		Cause: responseValid.cause,
//	})
//	if err != nil {
//		log.Println(err)
//	}
//	return dataJSON
//}
