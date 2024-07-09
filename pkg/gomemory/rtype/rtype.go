package rtype

import (
	"unsafe"
)

type ITab = uintptr

func GetITab(t any) ITab {
	return (*[2]uintptr)(unsafe.Pointer(&t))[0]
}
