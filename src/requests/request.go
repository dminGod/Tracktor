package requests

import (
	"net/http"
	"io/ioutil"
)


// Makes an HTTP Call
func Make_http_call( url string) []byte {


	// var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req, err := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic( err )
	}

	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)




	return body

	//return []byte("{ context : 123232323222}")
}

