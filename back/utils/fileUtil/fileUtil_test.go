package fileUtil

import (
	"os"
	"testing"
)

func TestGetFileMD5(t *testing.T) {
	file, err := os.OpenFile("D:/blueOcean/auto-music/img/huge.jpg", os.O_RDONLY, 0666)
	if err != nil {
		t.Error(err)
	}
	md5, err := GetFileMD5(file)
	if err != nil {
		t.Error(err)
	}
	t.Log(md5)
}
