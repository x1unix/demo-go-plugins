package main

import "C"
import (
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/x1unix/demo-go-plugins/sources/reddit/extension"
)

//export NewDataSource
func NewDataSource(rawCfg []byte, out *uintptr) error {
	cfg := new(extension.Config)
	if err := json.Unmarshal(rawCfg, cfg); err != nil {
		return fmt.Errorf("invalid configuration format (%s)", err)
	}

	ds := extension.NewDataSource(*cfg)
	*out = uintptr(unsafe.Pointer(ds))
	return nil
}
func main() {}
