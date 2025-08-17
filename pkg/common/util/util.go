package util

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"go.uber.org/zap"
)

func Elapsed(proc string) func() {
	start := time.Now()
	return func() {
		d := time.Since(start)
		if ms := d.Milliseconds(); ms < 10000 {
			zap.S().Infof("%s took %vms", proc, ms)
		} else {
			zap.S().Infof("%s took %v", proc, d)
		}
	}
}

func StructToMap(v interface{}) (m map[string]interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return
	}
	return
}

func GetMethod(fullMethod string) string {
	methods := strings.Split(fullMethod, "/")
	return methods[len(methods)-1]
}

// Minify returns json string with all whitespaces are removed.
func Minify(str string) ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := json.Compact(buff, []byte(str)); err != nil {
		return nil, err
	}

	return io.ReadAll(buff)
}
