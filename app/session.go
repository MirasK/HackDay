/*
	This file define session and describe them
*/

package app

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

// GetSesVal get value
// 	f, e := ioutil.ReadFile("sess/" + sid)
// 	for i, v := range arr
// 	arr := strings.Split(string(f), "\n")
// 	kv := strings.Split(v, "=")
// 	if kv[0] == key {
// 		return kv[1], nil
// 	}
func GetSesVal(key, sid string) (string, error) {
	f, e := ioutil.ReadFile("sess/" + sid)
	if e != nil {
		return "", e
	}

	arr := strings.Split(string(f), "\n")
	if len(arr) > 1 {
		for i, v := range arr {
			if i < len(arr)-1 {
				kv := strings.Split(v, "=")
				if kv[0] == key {
					return kv[1], nil
				}
			}
		}
	}
	return "", nil
}

// SetSesVal set value to session
// 	f, e := os.OpenFile("sess/"+sid, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
// 	defer f.Close()
// 	if e != nil {
// 		return e
// 	}
// 	f.WriteString(key + "=" + value + "\n")
// 	return nil
func SetSesVal(key, value, sid string) error {
	r, _ := GetSesVal(key, sid)
	if r != "" {
		WriteLog("value is exist")
		return errors.New("value is exist")
	}
	f, e := os.OpenFile("sess/"+sid, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer f.Close()
	if e != nil {
		return e
	}
	f.WriteString(key + "=" + value + "\n")
	return nil
}

// Delete ..
// func (sc SesControl) Delete(key string) error {
// 	delete(sc.Value, key)
// 	return nil
// }
