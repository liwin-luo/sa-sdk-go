package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

func DoRequest(url, args string, to time.Duration) error {
	var resp *http.Response

	data := bytes.NewBufferString(args)

	req, _ := http.NewRequest("POST", url, data)

	client := &http.Client{Timeout: to}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		data, _ := json.Marshal(struct {
			StatusCode int
			Body       string
		}{resp.StatusCode, string(body)})
		return errors.New(string(data))
	}
	return nil
}

func gzipData(data string) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err := zw.Write([]byte(data))
	if err != nil {
		zw.Close()
		return nil, err
	}
	zw.Close()

	return buf.Bytes(), nil
}

func encodeData(data string) (string, error) {
	gdata, err := gzipData(data)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(gdata)
	return encoded, nil
}

func GeneratePostDataList(data string) (string, error) {
	return encodeData(data)
}

func GeneratePostData(data string) (string, error) {
	return encodeData(data)
}

func NowMs() int64 {
	return time.Now().UnixNano() / 1000000
}

func NowS() int64 {
	return time.Now().Unix()
}

func DeepCopy(value map[string]interface{}) map[string]interface{} {
	ncopy := deepCopy(value)
	if nmap, ok := ncopy.(map[string]interface{}); ok {
		return nmap
	}

	return nil
}

func deepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = deepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = deepCopy(v)
		}

		return newSlice
	}

	return value
}

func ConvertMap(data interface{}) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})
	valueV := reflect.ValueOf(data)
	valueT := reflect.TypeOf(data)
	switch reflect.TypeOf(data).Kind() {
	case reflect.Ptr:
		valueV = valueV.Elem()
		valueT = valueT.Elem()
		fallthrough
	case reflect.Struct:
		for i := 0; i < valueT.NumField(); i++ {
			if name, ok := getTag(valueT.Field(i)); ok && valueV.Field(i).CanInterface() {
				res["$"+name] = valueV.Field(i).Interface()
			}
		}
	case reflect.Map:
	}
	return res, nil
}

func getTag(field reflect.StructField) (string, bool) {
	var tag string
	if tag = field.Tag.Get("sd"); len(tag) != 0 {
		return tag, true
	}
	return "", false
}
