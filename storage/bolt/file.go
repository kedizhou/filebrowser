package bolt

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/asdine/storm/v3"
	"github.com/filebrowser/filebrowser/v2/files"
)

type fileReadStatusBackend struct {
	db *storm.DB
}

type fileOwnerInfoBackend struct {
	db *storm.DB
}

func (s fileOwnerInfoBackend) GetfileOwnerInfo(Path string) ([]uint64, error) {
	var ownerInfo []files.FileOwnerInfo
	err := s.db.Find("Path", Path, &ownerInfo)
	if err != nil {
		// return false, nil
		if err.Error() != "not found" {
			return []uint64{0}, errors.New("query database found error: " + err.Error() + ", path: " + Path)
		}
	}
	var result []uint64
	for _, status := range ownerInfo {
		if string(status.Comments) == "tusPatchHandler:" || string(status.Comments) == "resourcePostHandler:create" {
			result = append(result, uint64(status.OwnerID))
		}
	}
	var strResult []string
	for _, userID := range result {
		strResult = append(strResult, strconv.FormatUint(uint64(userID), 10))
	}
	log.Println("query the file path: " + Path + ", the ower is (uid?): " + strings.Join(strResult, ","))
	return result, nil
}

func (s fileOwnerInfoBackend) SetfileOwnerInfo(Path string, OwnerID uint64, Comments string) error {
	log.Printf("set file onwner info: " + Path + ", comments: " + Comments)
	// currentTime := time.Now()
	fileOwnerInfo := files.FileOwnerInfo{
		Path:     Path,
		OwnerID:  OwnerID,
		Comments: Comments,
	}
	s.db.Save(&fileOwnerInfo)
	return errors.New("o")
}

func (s fileReadStatusBackend) SetFileReadStatus(l *files.ReadStatus) error {
	r, err := s.FoundFileReadStatus(l)
	if r {
		log.Printf("had set read status, ignore the path: " + l.Path + " ,userid:" + strconv.FormatUint(uint64(l.UserID), 10))
		return err
	} else if err != nil && err.Error() == "not found" {

	} else if err != nil && !r {
		log.Printf("mark file read status had error." + err.Error())
		return err
	}
	log.Printf("mark file read status, path: " + l.Path + ", filename:" + l.FileName)
	return s.db.Save(l)
}

func (s fileReadStatusBackend) CleanFileReadStatus(l *files.ReadStatus) error {
	return s.db.DeleteStruct(l)
}

func (s fileReadStatusBackend) FoundFileReadStatus(l *files.ReadStatus) (bool, error) {
	var filesStatus []files.ReadStatus
	// err := s.db.One("FileName", l.FileName, &status)
	err := s.db.Find("FileName", l.FileName, &filesStatus)

	if err != nil {
		// return false, nil
		if err.Error() != "not found" {
			return false, errors.New("query database found error: " + err.Error() + ", filename: " + l.FileName)
		}
	}

	// var result []uint64
	for _, status := range filesStatus {
		if status.FileName == l.FileName && status.Path == l.Path && status.UserID == l.UserID {
			return true, nil
		}
	}

	return false, errors.New("not found") // 没有匹配成功
}

func (s fileReadStatusBackend) UpdateFileReadStatus(old *files.ReadStatus, new *files.ReadStatus) (bool, error) {
	var filesStatus []files.ReadStatus
	err := s.db.Find("Path", old.Path, &filesStatus)
	if err != nil {
		// return false, nil
		if err.Error() != "not found" {
			return false, errors.New("query database found error: " + err.Error() + ", path: " + old.Path)
		}
	}
	// var result []uint64
	for _, status := range filesStatus {
		status.Path = new.Path
		if err := s.db.Update(&status); err != nil {
			log.Panicln("update file read status by path:" + old.Path + ", err" + err.Error())
			return false, errors.New("update file read status by path:" + old.Path + ", err" + err.Error())
		}
	}

	return true, nil
}

func (s fileReadStatusBackend) ListFileReadStatus(l *files.ReadStatus) ([]uint64, error) {
	var statuses []files.ReadStatus
	err := s.db.Find("Path", l.Path, &statuses)
	if err != nil {
		log.Println("to list file read status record: " + err.Error() + ", for file path: " + l.Path)
		return []uint64{}, err
	}
	var result []uint64
	for _, config := range statuses {
		if config.Path == l.Path { // 再根据 "path" 过滤
			result = append(result, uint64(config.UserID))
		}
	}

	var strResult []string
	for _, userID := range result {
		strResult = append(strResult, strconv.FormatUint(uint64(userID), 10))
	}
	log.Println("query the file path: " + l.Path + ",how many people to read this file(uid?): " + strings.Join(strResult, ","))
	// return strings.Join(strResult, ","), nil
	return result, nil
}
