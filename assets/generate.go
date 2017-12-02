package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strconv"
)

// A go generate script to inject image strings strings.
// simply add //go:generate go run path/to/this/file/generate.go
// to your file, then run go generate, and below it will be
// const icon = "..."
func main() {

	// Check the args for the filename.
	fmt.Println(filepath.Abs("."))

	fmt.Println(os.Args)
	if len(os.Args) != 2 {
		log.Fatal("Usage: generate.go [image.png]")
		return
	}

	// GOFILE and GOLINE are set by go generate.
	goFile := os.Getenv("GOFILE")
	goLineString := os.Getenv("GOLINE")
	goLine := -1
	if goLineString != "" {
		line, err := strconv.Atoi(goLineString)
		if err != nil {
			log.Fatal("Invalid goLine")
			return
		}
		goLine = line
	}

	imageString, err := getImageString(os.Args[1])
	if err != nil {
		log.Fatal(err)
		return
	}

	parentDir, err := getDirectory()
	if err != nil {
		log.Fatal(err)
		return
	}

	file, err := os.Open(parentDir + string(os.PathSeparator) + goFile)
	if err != nil {
		log.Fatal(err)
		return
	}

	readBytes, err := ioutil.ReadAll(bufio.NewReader(file))
	if err != nil {
		log.Fatal(err)
		return
	}

	byteBuf := new(bytes.Buffer)
	writer := bufio.NewWriter(byteBuf)
	reader := bufio.NewReader(bytes.NewReader(readBytes))

	readWriter := bufio.NewReadWriter(reader, writer)

	lineNumber := 0
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal("Error reading file", err)
				return
			}
		}
		lineNumber++
		if lineNumber == goLine+1 {
			readWriter.WriteString(fmt.Sprintf("const icon = %q\n", imageString))
		} else {
			readWriter.Write(line)
			readWriter.Write([]byte{'\n'})
		}
	}

	readWriter.Flush()
	ioutil.WriteFile(file.Name(), byteBuf.Bytes(), os.FileMode(0644))

}

// Convert the file to base64
func getImageString(imageFileName string) (string, error) {

	imgFile, err := os.Open(imageFileName)
	if err != nil {
		return "", fmt.Errorf("unable to open image file. error: %q", err)
	}

	defer imgFile.Close()

	imageBody, err := ioutil.ReadAll(imgFile)
	if err != nil {
		return "", fmt.Errorf("unable to read image file. error: %q", err)
	}

	_, format, err := image.Decode(bytes.NewBuffer(imageBody))
	if err != nil {
		return "", fmt.Errorf("unable to decode image. error: %q", err)
	}

	data := mime.TypeByExtension("." + format)
	img64 := base64.StdEncoding.EncodeToString(imageBody)
	return "data:" + data + ";base64," + img64, nil

}

func getDirectory() (string, error) {

	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("no working directory O_O. Error: %q", err)
	}

	return filepath.Abs(workingDir)

}
