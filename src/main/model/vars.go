package model

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound

type BitBool bool

//func (b BitBool) Value() (driver.Value, error) {
//	if b {
//		return []byte{1}, nil
//	} else {
//		return []byte{0}, nil
//	}
//}
//
//func (b *BitBool) Scan(src interface{}) error {
//	v, ok := src.([]byte)
//	if !ok {
//		return errors.New("bad []byte type assertion")
//	}
//	*b = v[0] == 1
//	return nil
//}
