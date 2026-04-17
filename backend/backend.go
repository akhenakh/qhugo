package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/image/draw"
)

var hugoCmd *exec.Cmd

//export InitBackend
func InitBackend() {
}

//export ReadFileContent
func ReadFileContent(path *C.char) *C.char {
	b, err := os.ReadFile(C.GoString(path))
	if err != nil {
		return C.CString("")
	}
	return C.CString(string(b))
}

//export SaveFileContent
func SaveFileContent(path *C.char, content *C.char) int {
	err := os.WriteFile(C.GoString(path), []byte(C.GoString(content)), 0644)
	if err != nil {
		log.Println(err)
		return 0
	}
	return 1
}

//export FreeString
func FreeString(str *C.char) {
	C.free(unsafe.Pointer(str))
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close() // Closes the listener, freeing the port for Hugo
	return l.Addr().(*net.TCPAddr).Port, nil
}

//export StartHugo
func StartHugo(repoC *C.char) int {
	repo := C.GoString(repoC)
	if hugoCmd != nil && hugoCmd.Process != nil {
		hugoCmd.Process.Kill()
	}

	port, err := getFreePort()
	if err != nil {
		log.Println("Error finding free port, falling back to 1313:", err)
		port = 1313
	}

	// Launch hugo in background
	hugoCmd = exec.Command("hugo", "server", "-s", repo, "-p", fmt.Sprintf("%d", port), "-D")
	hugoCmd.Start()

	return port
}

//export StopHugo
func StopHugo() {
	if hugoCmd != nil && hugoCmd.Process != nil {
		hugoCmd.Process.Kill()
		hugoCmd = nil
	}
}

// CatmullRom resize operation to scale down to blog-acceptable size
func resizeImage(srcPath, dstPath string) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width > 1200 {
		ratio := float64(height) / float64(width)
		newWidth := 1200
		newHeight := int(float64(newWidth) * ratio)

		dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
		img = dst
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, img, &jpeg.Options{Quality: 85})
}

//export ProcessImage
func ProcessImage(srcC, repoC, docC *C.char) *C.char {
	src := C.GoString(srcC)
	repo := C.GoString(repoC)
	doc := C.GoString(docC)

	docName := filepath.Base(doc)
	docPrefix := strings.TrimSuffix(docName, filepath.Ext(docName))

	imgDir := filepath.Join(repo, "static", "img")
	os.MkdirAll(imgDir, 0755)

	srcName := filepath.Base(src)
	dstName := fmt.Sprintf("%s-%s", docPrefix, srcName)
	dstName = strings.TrimSuffix(dstName, filepath.Ext(dstName)) + ".jpg"

	dstPath := filepath.Join(imgDir, dstName)

	err := resizeImage(src, dstPath)
	if err != nil {
		return C.CString(fmt.Sprintf("Error: %v", err))
	}

	markdownLink := fmt.Sprintf("![%s](/img/%s)", srcName, dstName)
	return C.CString(markdownLink)
}

func main() {}
