package ldc

import (
	"testing"
)

// go test -v -run="Compile" -args f=ESP32/dht-test.zip b=esp32 c=esp32duino-virtual
func TestCompile(t *testing.T)  {
	compileInfo := getCompileInfo()
	compileFile(compileInfo)
}

func TestDownload(t *testing.T) {
	compileInfo := getCompileInfo()
	downloadCompiledFile(compileInfo, "test.bin")
}