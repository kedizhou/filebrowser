package bolt

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/asdine/storm/v3"

	fbErrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/users"
)

type usersBackend struct {
	db *storm.DB
}

func (st usersBackend) GetById(i uint64) (user *users.User, err error) {
	user = &users.User{}
	err = st.db.One("ID", i, user)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (st usersBackend) GetBy(i interface{}) (user *users.User, err error) {
	user = &users.User{}

	var arg string
	switch i.(type) {
	case uint:
		arg = "ID"
	case string:
		arg = "Username"
	default:
		return nil, fbErrors.ErrInvalidDataType
	}

	err = st.db.One(arg, i, user)

	if err != nil {
		if errors.Is(err, storm.ErrNotFound) {
			return nil, fbErrors.ErrNotExist
		}
		return nil, err
	}

	return
}

func (st usersBackend) Gets() ([]*users.User, error) {
	var allUsers []*users.User
	err := st.db.All(&allUsers)
	if errors.Is(err, storm.ErrNotFound) {
		return nil, fbErrors.ErrNotExist
	}

	if err != nil {
		return allUsers, err
	}

	return allUsers, err
}

func (st usersBackend) Update(user *users.User, fields ...string) error {
	if len(fields) == 0 {
		return st.Save(user)
	}

	for _, field := range fields {
		userField := reflect.ValueOf(user).Elem().FieldByName(field)
		if !userField.IsValid() {
			return fmt.Errorf("invalid field: %s", field)
		}
		val := userField.Interface()
		if err := st.db.UpdateField(user, field, val); err != nil {
			return err
		}
	}

	return nil
}

func (st usersBackend) Save(user *users.User) error {
	err := st.db.Save(user)
	if errors.Is(err, storm.ErrAlreadyExists) {
		return fbErrors.ErrExist
	}
	return err
}

func (st usersBackend) DeleteByID(id uint) error {
	log.Println("Delete by id: " + strconv.FormatUint(uint64(id), 10))
	return st.db.DeleteStruct(&users.User{ID: id})
}

func (st usersBackend) DeleteByUsername(username string) error {
	user, err := st.GetBy(username)
	if err != nil {
		return err
	}

	return st.db.DeleteStruct(user)
}
