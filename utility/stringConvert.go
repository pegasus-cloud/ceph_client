package utility

import "unsafe"

// StringHeader ...
type StringHeader struct {
	Data unsafe.Pointer
	Len  int
}

// SliceHeader ...
type SliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// StrToBytes ...
func StrToBytes(s string) []byte {
	header := (*StringHeader)(unsafe.Pointer(&s))
	bytesHeader := &SliceHeader{
		Data: header.Data,
		Len:  header.Len,
		Cap:  header.Len,
	}
	return *(*[]byte)(unsafe.Pointer(bytesHeader))
}

// BytesToStr ....
func BytesToStr(b []byte) string {
	header := (*SliceHeader)(unsafe.Pointer(&b))
	strHeader := &StringHeader{
		Data: header.Data,
		Len:  header.Len,
	}
	return *(*string)(unsafe.Pointer(strHeader))
}
