package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"umn-technology/constants"
)

func WorkerPost(url, authorization string, requestData interface{}) ([]byte, error) {
	bodyreq, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Err Worker Post - json.Marshal : ", err.Error())
		return nil, err
	}

	bodyreqtostr := string(bodyreq)
	fmt.Println("Request :> ", bodyreqtostr)

	httpreq, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyreq))
	if err != nil {
		fmt.Println("Err Worker Post - http.NewRequest ", url, " :", err.Error())
		return nil, err
	}
	defer httpreq.Body.Close()

	httpreq.Close = constants.TRUE_VALUE
	httpreq.Header.Add("Content-Type", "application/json")
	httpreq.Header.Set("Connection", "close")
	if authorization != constants.EMPTY_VALUE {
		httpreq.Header.Set("Authorization", authorization)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: constants.TRUE_VALUE},
	}

	client := &http.Client{Transport: tr}
	defer client.CloseIdleConnections()

	fmt.Println("URL :>", url)
	fmt.Println("Request :>", requestData)

	resp, err := client.Do(httpreq)
	if err != nil {
		fmt.Println("Err Worker Post - client.Do :", err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	byteResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Err Worker Post - ioutil.ReadAll :", err.Error())
		return nil, err
	}

	bodyString := string(byteResponse)
	fmt.Println("Response :>", bodyString)

	return byteResponse, nil
}

func WorkerGet(url, authorization string, requestData interface{}) ([]byte, error) {
	bodyreq, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Err Worker Post - json.Marshal : ", err.Error())
		return nil, err
	}

	bodyreqtostr := string(bodyreq)
	fmt.Println("Request :>", bodyreqtostr)

	httpreq, err := http.NewRequest("GET", url, bytes.NewBuffer(bodyreq))
	if err != nil {
		fmt.Println("Err Worker Post - http.NewRequest ", url, " :", err.Error())
		return nil, err
	}
	defer httpreq.Body.Close()

	httpreq.Close = constants.TRUE_VALUE
	httpreq.Header.Add("Content-Type", "application/json")
	httpreq.Header.Set("Connection", "close")
	if authorization != constants.EMPTY_VALUE {
		httpreq.Header.Set("Authorization", authorization)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: constants.TRUE_VALUE},
	}

	client := &http.Client{Transport: tr}
	defer client.CloseIdleConnections()

	fmt.Println("URL :>", url)

	resp, err := client.Do(httpreq)
	if err != nil {
		fmt.Println("Err Worker Post - client.Do :", err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	byteResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Err Worker Post - ioutil.ReadAll :", err.Error())
		return nil, err
	}

	bodyString := string(byteResponse)
	fmt.Println("Response :>", bodyString)

	return byteResponse, nil
}
