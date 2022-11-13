package server

// response data type to answer the request in the JSON format
//type response struct {
//	nrp.Simple
//id     int
//body   interface{}
//client *client
//}

// setID set id of response
//func (response *response) setID(id int) {
//	response.id = id
//}

// setBody set body of response
//func (response *response) setBody(body []byte) {
//	err := json.Unmarshal(body, &response.body)
//	if err != nil {
//		log.Println(err)
//	}
//}

// setClient set a link to the client in response
//func (response *response) setClient(client *client) {
//	response.client = client
//}

// exportJSON returns a response in JSON format
//func (response *response) exportJSON() []byte {
//	jsonData, err := json.Marshal(struct {
//		Id   int         `json:"id"`
//		Body interface{} `json:"body"`
//	}{
//		Id:   response.id,
//		Body: response.body,
//	})
//	if err != nil {
//		log.Println(err)
//	}
//	return jsonData
//}

// write send data to websocket channel of client
//func (response *response) write(data []byte) {
//	response.client.send <- data
//}
