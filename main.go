package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	fmt.Println("main enter")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-quit
	fmt.Printf("main exit, %v\n", sig)
	fmt.Println("main exit")

	mux := http.NewServeMux()
	mux.HandleFunc("/exportFile", func(w http.ResponseWriter, r *http.Request) {

	})
	mux.HandleFunc("/importFile", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		fmt.Println(readRow(file))
	})
}

func exporter(writer io.Writer, bytes []byte) error {
	_, err := writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func readRow(reader io.Reader) string {
	file, err := excelize.OpenReader(reader)
	if err != nil {
		fmt.Printf("OpenReader err: %v\n", err)
		return ""
	}
	rows, err := file.GetRows("sheet1")
	if err != nil {
		fmt.Printf("GetRows err: %v\n", err)
		return ""
	}
	defer file.Close()

	return strings.Join(rows[0], " ")
}
