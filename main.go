package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"test_struct/parse"
)

var (
	fa   = flag.String("fa", "", "文件路径地址Annotations")
	fq   = flag.String("fq", "", "文件路径地址questions")
	dele = flag.String("dele", "", "文件路径地址wait to delete")
)

func main() {
	flag.Parse()

	recordInvalidImageIds := parse.Parse(*fa, *fq, parse.PictureArray)

	fmt.Println(recordInvalidImageIds)

	var waitDeleteFiles []string
	for _, e := range recordInvalidImageIds {
		waitDeleteFiles = append(waitDeleteFiles,
			path.Join(*dele, "test2015", fmt.Sprintf("%d.npy", e)),
			path.Join(*dele, "train2014", fmt.Sprintf("%d.npy", e)),
			path.Join(*dele, "val2014", fmt.Sprintf("%d.npy", e)),
		)
	}
	DeleteFiles(waitDeleteFiles)
}

func DeleteFiles(paths []string) {
	for _, e := range paths {
		if FileIsExisted(e) {
			fmt.Println(e)
			//err := os.Remove(e)
			//panic(err)
		}
	}
}

func FileIsExisted(filename string) bool {
	existed := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		existed = false
	}
	return existed
}
