package cmd

import (
	"fmt"
	"net/http"
)

func StartWeb() {
	fmt.Println("[Web] 启动 Web UI： http://localhost:8000")
	http.Handle("/", http.FileServer(http.Dir("web/ui")))
	http.ListenAndServe(":8000", nil)
}
