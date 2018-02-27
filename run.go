// +build ignore

package main

import (
	"net/http"

	"github.com/not-quite-vacation/blog/blog"
)

func main() {
	http.ListenAndServe(":8080", http.FileServer(blog.FS(false)))
}
