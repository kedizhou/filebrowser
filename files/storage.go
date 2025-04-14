package files

import (
	"log"
	"strings"
)

// Link is the information needed to build a shareable link.
type ReadStatus struct {
	// Hash     string `json:"hash" storm:"id,index"`
	ID       uint   `storm:"id,increment" json:"id"`
	Path     string `json:"path" storm:"index"`
	FileName string `json:"filename" storm:"index"`
	UserID   uint   `json:"userID" storm:"index"`
	SetDate  string `json:"setdate" storm:"index"`
	// Expire       int64  `json:"expire"`
	// PasswordHash string `json:"password_hash,omitempty"`
	// // Token is a random value that will only be set when PasswordHash is set. It is
	// // URL-Safe and is used to download links in password-protected shares via a
	// // query arg.
	// Token string `json:"token,omitempty"`
}

type FileOwnerInfo struct {
	ID       uint   `storm:"id,increment" json:"id"`
	Path     string `json:"path" storm:"index"`
	OwnerID  uint64 `json:"userID" storm:"index"`
	Comments string `json:"comments" storm:"index"`
	SetDate  string `json:"setdate" storm:"index"`
}

// StorageBackend is the interface to implement for a share storage.
type fileReadStatus interface {
	// Delete(hash string) error
	SetFileReadStatus(*ReadStatus) error
	CleanFileReadStatus(*ReadStatus) error
	FoundFileReadStatus(*ReadStatus) (bool, error)
	ListFileReadStatus(*ReadStatus) ([]uint64, error)
	UpdateFileReadStatus(*ReadStatus, *ReadStatus) (bool, error)
}

type fileOwnerInfo interface {
	SetfileOwnerInfo(Path string, OwnerID uint64, Comments string) error
	GetfileOwnerInfo(Path string) ([]uint64, error)
}

// ReadStatus 定义了一个类型
type Storage struct {
	back fileReadStatus
}

// ReadStatus 定义了一个类型
type StorageFileOwnerInfo struct {
	back       fileOwnerInfo
	readStatus fileReadStatus
}

func NewReadStatus(back fileReadStatus) *Storage {
	return &Storage{back: back}
}

// func NewOwnerInfo(back fileOwnerInfo, fileReadStatus fileReadStatus) *StorageFileOwnerInfo {
func NewOwnerInfo(back fileOwnerInfo, fileReadStatus fileReadStatus) *StorageFileOwnerInfo {
	return &StorageFileOwnerInfo{back: back, readStatus: fileReadStatus}
}

func (s *StorageFileOwnerInfo) GetfileOwnerInfo(Path string) ([]uint64, error) {
	return s.back.GetfileOwnerInfo(Path)
}

func (s *StorageFileOwnerInfo) SetfileOwnerInfo(Path string, OwnerID uint64, Comments string) error {
	if strings.Contains(Comments, ":rename") {
		parts := strings.Split(Path, ",")
		src := strings.Replace(parts[0], "src:", "", 1)
		dst := strings.Replace(parts[1], "dst:", "", 1)
		// frsb := fileReadStatusBackend{}
		status_old := ReadStatus{
			Path: src,
		}
		status_new := ReadStatus{
			Path: dst,
		}
		log.Println("update src" + src + ",dst:" + dst)
		_, err := s.readStatus.UpdateFileReadStatus(&status_old, &status_new)
		if err != nil {
			return err
		}
	}
	return s.back.SetfileOwnerInfo(Path, OwnerID, Comments)
}

// Save saves the settings for the current instance.
func (s *Storage) SetFileReadStatus(set *ReadStatus) error {
	return s.back.SetFileReadStatus(set)
}

func (s *Storage) CleanFileReadStatus(set *ReadStatus) error {
	return s.back.CleanFileReadStatus(set)
}

// Save saves the settings for the current instance.
func (s *Storage) FoundFileReadStatus(set *ReadStatus) (bool, error) {
	return s.back.FoundFileReadStatus(set)
}

func (s *Storage) UpdateFileReadStatus(old *ReadStatus, new *ReadStatus) (bool, error) {
	return s.back.UpdateFileReadStatus(old, new)
}

// Save saves the settings for the current instance.
func (s *Storage) ListFileReadStatus(set *ReadStatus) ([]uint64, error) {
	return s.back.ListFileReadStatus(set)
}
