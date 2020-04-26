package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

const tmpDir string = "/tmp"

func main() {
	http.HandleFunc("/image", imageHandler)

	log.Fatal(http.ListenAndServe(":80", nil))
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	if http.MethodPost != r.Method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Forbidden method")
		return
	}

	r.ParseMultipartForm(32 << 20)
	image, handler, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "image is required")
		return
	}
	defer image.Close()

	targetType := r.FormValue("type")
	if targetType == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "type is required")
		return
	}

	if 0 == strings.Index(targetType, ".") {
		targetType = strings.TrimPrefix(targetType, ".")
	}

	unique := make([]byte, 4)
	_, err = rand.Read(unique)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failure happened when generate filename")
		return
	}

	file, err := os.OpenFile(tmpDir+"/"+fmt.Sprintf("%X", unique)+"-"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failure happened when creating file")
		return
	}

	_, err = io.Copy(file, image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failure happened when writing file")
		return
	}

	targetPath := strings.TrimSuffix(file.Name(), path.Ext(file.Name())) + "." + targetType
	targetName := path.Base(targetPath)
	cmd := exec.Command("convert", file.Name(), targetPath)
	err = cmd.Run()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failure happened when convert file")
		return
	}

	targetFile, err := os.Open(targetPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failure happened when return file")
		return
	}

	fileHeader := make([]byte, 512)
	targetFile.Read(fileHeader)
	fileContentType := http.DetectContentType(fileHeader)

	fileStat, _ := targetFile.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)

	w.Header().Set("Content-Disposition", "attachment; filename="+targetName)
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Set("Content-Length", fileSize)

	targetFile.Seek(0, 0)
	io.Copy(w, targetFile)
	return
}
