package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/filebrowser/filebrowser/v2/files"
)

// var Set = withUser(func(_ http.ResponseWriter, _ *http.Request, d *data) (int, error) {
// 	err := d.store.Users.Delete(d.raw.(uint))
// 	if err != nil {
// 		return errToStatus(err), err
// 	}

// 	return http.StatusOK, nil
// })

// var cleanFileReadStatus = withUser(func(_ http.ResponseWriter, _ *http.Request, d *data) (int, error) {
// 	err = d.store.FileReadStatus.Set()
// 	// err := d.store.Users.Delete(d.raw.(uint))
// 	if err != nil {
// 		return errToStatus(err), err
// 	}

// 	return http.StatusOK, nil
// })

type RequestData struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func setFileReadStatus() handleFunc {
	return withUser(func(_ http.ResponseWriter, r *http.Request, d *data) (int, error) {
		if r.Method != http.MethodPost {
			return http.StatusForbidden, errors.New("invalid request method")
		}
		currentTime := time.Now()
		var data RequestData
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {
			return http.StatusBadRequest, errors.New("failed to parse JSON body")
		}
		var path, err = url.QueryUnescape(data.Path)
		if err != nil {
			log.Println("parse path:"+data.Path+", error:", err)
			path = data.Path
		}
		readStatus := files.ReadStatus{
			Path:     path,
			FileName: data.Name,
			UserID:   d.user.ID,
			SetDate:  currentTime.Format(time.RFC3339),
		}
		d.store.FileReadStatus.SetFileReadStatus(&readStatus)
		// log.Println("debug file read status.")
		return http.StatusOK, nil
	})
}

func getFileReadStatus() handleFunc {
	return withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		if r.Method != http.MethodPost {
			return http.StatusForbidden, errors.New("invalid request method")
		}
		currentTime := time.Now()
		var data RequestData
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {
			return http.StatusBadRequest, errors.New("failed to parse JSON body")
		}
		var path, err = url.QueryUnescape(data.Path)
		if err != nil {
			log.Println("parse path:"+data.Path+", error:", err)
			path = data.Path
		}
		readStatus := files.ReadStatus{
			Path:     path,
			FileName: data.Name,
			UserID:   d.user.ID,
			SetDate:  currentTime.Format(time.RFC3339),
		}
		userIDs, _ := d.store.FileReadStatus.ListFileReadStatus(&readStatus)
		var usernames []string
		for _, userID := range userIDs {
			var user, err = d.store.Users.GetById(userID)
			if err != nil {
				log.Println("fetch username through userid error: " + err.Error())
				return 1, err
			}
			// log.Println("convert id to username: " + user.Username)
			usernames = append(usernames, user.Username)
		}
		// username
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strings.Join(usernames, ",")))
		// return http.StatusOK, nil
		return 0, nil

	})
}

func getOwner() handleFunc {
	return withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		if r.Method != http.MethodPost {
			return http.StatusForbidden, errors.New("invalid request method")
		}
		var data RequestData
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {
			return http.StatusBadRequest, errors.New("failed to parse JSON body")
		}
		var path, err = url.QueryUnescape(data.Path)
		if err != nil {
			log.Println("parse path:"+data.Path+", error:", err)
			path = data.Path
		}
		userIDs, _ := d.store.FilesOwnerInfo.GetfileOwnerInfo(path)
		var usernames []string
		for _, userID := range userIDs {
			var user, err = d.store.Users.GetById(userID)
			if err != nil {
				log.Println("fetch username through userid error: " + err.Error())
				return 1, err
			}
			// log.Println("convert id to username: " + user.Username)
			checkRepeat := 0
			for _, v := range usernames {
				if v == user.Username {
					checkRepeat = 1
				}
			}
			if checkRepeat == 0 {
				usernames = append(usernames, user.Username)
			}
		}
		// username
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strings.Join(usernames, ",")))
		// return http.StatusOK, nil
		return 0, nil

	})
}
