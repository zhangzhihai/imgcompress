package main

import (
	"bytes"
	//"errors"
	"flag"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"
	"ulucu.github.com/log.v1"
)

type config struct {
	filepath string
	dist     string
	level    int
	acy      int
}

var (
	cfg config
)

func init() {
	//定义日志类型
	log.Std = log.New(os.Stderr, "", log.Ldate|log.Ltime)
	log.Std.Level = 1

	flag.StringVar(&cfg.filepath, "filepath", "", "filepath")
	flag.StringVar(&cfg.dist, "dist", "", "dist")
	flag.IntVar(&cfg.acy, "acy", 40, "0-100")
	flag.IntVar(&cfg.level, "level", 1, "email list")
}

func CuttingJpeg(Reader io.Reader, w int, h int, acy int) (b *bytes.Buffer, err error) {
	log.Infof("w:%d,h:%d,acy:%d", w, h, acy)
	origin, _, err := image.Decode(Reader)
	if err != nil {
		return nil, err
	}
	x := uint(w)
	y := uint(h)

	canvas := resize.Resize(x, y, origin, resize.Lanczos3)

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, canvas, &jpeg.Options{acy})
	return buf, nil
}

//读取目录下的文件列表
func WalkDir(dirPth, suffix string) (files []string, err error) {

	files = make([]string, 0, 30)

	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录

		if err != nil { //忽略错误
			return err
		}

		if fi.IsDir() { // 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})

	log.Info(files)

	for _, v := range files {
		Imgsave(v)
	}
	return files, err

}

func Imgsave(file string) {
	log.Infof("%s\n", file)
	pathlist := strings.Split(file, ".")
	sum := len(pathlist)
	ext := pathlist[sum-1]

	fi, err := os.Open(file)
	if err != nil {
		log.Error(err)
		return
	}
	out_buf, err := CuttingJpeg(fi, 0, 0, cfg.acy)
	fi.Close()

	if err != nil {
		log.Error(err)
		return
	}

	err = os.Rename(file, fmt.Sprintf("%s.%s", file, ext))
	if err != nil {
		log.Error(err)
		return
	}
	w, _ := os.Create(file)

	io.WriteString(w, string(out_buf.Bytes()))
	w.Close()
	return
}

func main() {

	flag.Parse()
	log.Std.Level = cfg.level

	log.Info(cfg)

	if cfg.dist == "" {
		Imgsave(cfg.filepath)
	} else if cfg.filepath == "" {
		WalkDir(cfg.dist, "jpg")
	} else {
		log.Error("flag err")
	}

}
