package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jessevdk/go-flags"
)

func GetAllfiles(fileFolder string) (files []string, err error) {
	err = filepath.Walk(fileFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// fmt.Println(info.Name())
			filepath := path
			filepath = "./" + strings.Replace(filepath, "\\", "/", -1)
			fmt.Println("init:", filepath)
			files = append(files, filepath)
		}
		return nil
	})
	return
}

func GetFileExt(filePath string) (fileName, ext string) {
	extensionIndex := strings.LastIndex(filePath, ".")
	if extensionIndex == -1 {
		fmt.Println("The file has no extension")
		return
	}
	ext = filePath[extensionIndex:]
	fileName = filePath[:extensionIndex]
	// fmt.Println("The extension of", filePath, "is", fileName, ext)
	// fmt.Printf("fileName: %s,ext is %s", fileName, ext)
	return
}

func ConvertSrtToTxt(wg *sync.WaitGroup, filePath, fileName string) (err error) {
	defer wg.Done()
	content_lines, err := GetSrtFileContent(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(content_lines)
	txtPath := fileName + ".txt"
	err = WriteContentToTxt(txtPath, content_lines)
	return
}

func WriteContentToTxt(filePath string, content_lines []string) (err error) {
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("start creating %v ... \n", filePath)
	defer f.Close()
	var buf bytes.Buffer
	every := 0
	for _, line := range content_lines {
		n, err := buf.WriteString(line)
		if err != nil {
			fmt.Println(err)
			return err
		}
		every += n
		if every > 100 {
			_, err = buf.WriteTo(f)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return err
			}
			every = 0
			fmt.Printf("writing %v... \n", filePath)
		}

	}
	_, err = buf.WriteTo(f)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	return
}

func GetSrtFileContent(filePath string) (lines []string, err error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 解析字幕内容
	subtitles := strings.Split(string(content), "\r\n\r\n")
	for _, subtitle := range subtitles {
		// fmt.Println(subtitle)
		blockLines := strings.Split(subtitle, "\r\n")

		if len(blockLines) >= 3 {
			// 第一行是序号，第二行是时间戳，第三行开始是文本
			text := strings.Join(blockLines[2:], "\r\n") + "\r\n"
			// fmt.Printf("%v", text)
			lines = append(lines, text)

		}
	}
	return
}

type Options struct {
	Folder string `short:"f" long:"folder" description:"Folder option"`
}

var opts Options

func main() {
	parser := flags.NewParser(&opts, flags.Default)

	// 解析命令行参数
	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	fmt.Println("Folder:", opts.Folder)
	// args := flag.Args()

	folder := opts.Folder
	files, err := GetAllfiles(folder)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%v", files)
	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, file := range files {
		filePath := file
		if fileName, ext := GetFileExt(filePath); ext != ".srt" {
			wg.Done()
			continue
		} else {

			// fmt.Println(fileName, ext)
			fmt.Println("filePath:", filePath)
			go ConvertSrtToTxt(&wg, filePath, fileName)
			// fmt.Println(fileName,ext)
		}

	}
	wg.Wait()
	fmt.Println("\nends...")

}
