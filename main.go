package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/dchest/captcha"
)

var (
	flagLen   = flag.Int("count", 2, "for every number, gen how many sample, 1-30")
	flagValue = flag.Int("value", -1, "the value for sample command, 0-9999")
)

func usage() {
	fmt.Println("usage: captchagen [flags] ACTION")
	fmt.Println("       ACTION: sample or genbatch")
	flag.PrintDefaults()
}

func genData(count int) {
	for j := 1; j <= count; j++ {
		for i := 0; i < 10000; i++ {
			path := fmt.Sprintf("out%s%02d%s%02d%s",
				string(filepath.Separator), i/100, string(filepath.Separator), i%100, string(filepath.Separator))
			fileName := fmt.Sprintf("%04d-%02d.png", i, j)
			os.MkdirAll(path, 0777)
			createOneCaptcha(i, path+fileName)
			fmt.Println(path + fileName + " generated")
		}
	}
}

func createOneCaptcha(value int, fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer f.Close()
	var d []byte
	if value == -1 {
		d = captcha.RandomDigits(4)
	} else {
		d = make([]byte, 4)
		d[0] = (byte)((value / 1000) % 10)
		d[1] = (byte)((value / 100) % 10)
		d[2] = (byte)((value / 10) % 10)
		d[3] = (byte)(value % 10)
	}
	var w io.WriterTo
	w = captcha.NewImage(fileName, d, 180, 60)
	_, err = w.WriteTo(f)
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func main() {
	flag.Parse()
	action := flag.Arg(0)
	if action == "sample" {
		createOneCaptcha(*flagValue, "sample.png")
		fmt.Println("gen a sample captcha to sample.png")
	} else if action == "genbatch" {
		if *flagLen <= 0 || *flagLen > 30 {
			usage()
			os.Exit(1)
		}
		fmt.Println("genbatch, count:", *flagLen)
		genData(*flagLen)
	} else {
		usage()
		os.Exit(1)
	}
}
