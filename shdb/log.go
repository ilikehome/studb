package shdb

import "os"

type log struct{
	f *os.File
}

type logLine struct{
	seq int64
	len int64
	content []byte
}

func main() {
}
