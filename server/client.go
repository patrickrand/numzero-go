package server

import (
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path"
	"path/filepath"

	"github.com/emicklei/go-restful"
	"github.com/nkcraddock/gooby"
)

type StaticContentHandler struct {
	notFound    []byte
	contentRoot string
}

func RegisterStaticContent(container *restful.Container, root string) *StaticContentHandler {
	ws := new(restful.WebService)
	h := new(StaticContentHandler)

	if root == "" {
		log.Println("Hosting static content in memory")
		notFound, _ := gooby.Asset("404.html")
		h.notFound = notFound
		ws.Route(ws.GET("/{path:*}").To(h.serveBinData))
	} else {
		cur, _ := os.Getwd()
		h.contentRoot = path.Join(cur, root)
		log.Println("Hosting static content at", h.contentRoot)
		notFound, _ := ioutil.ReadFile(path.Join(h.contentRoot, "404.html"))
		h.notFound = notFound
		ws.Route(ws.GET("/{path:*}").To(h.serveFileSystem))
	}

	container.Add(ws)

	return h
}

func (h *StaticContentHandler) serveBinData(req *restful.Request, res *restful.Response) {
	filePath := req.PathParameter("path")

	if filePath == "" {
		filePath = "index.html"
	}

	if data, err := gooby.Asset(filePath); err == nil {
		mimetype := mime.TypeByExtension(filepath.Ext(filePath))
		res.AddHeader("Content-Type", mimetype)
		res.Write(data)
	} else {
		log.Println("NOT FOUND:", filePath)
		res.AddHeader("Content-Type", "text/html")
		res.Write(h.notFound)
	}
}

func (h *StaticContentHandler) serveFileSystem(req *restful.Request, res *restful.Response) {
	filePath := req.PathParameter("path")

	if filePath == "" {
		filePath = "index.html"
	}

	filePath = path.Join(h.contentRoot, filePath)

	if _, err := os.Stat(filePath); err == nil {
		data, err := ioutil.ReadFile(filePath)
		if err == nil {
			mimetype := mime.TypeByExtension(filepath.Ext(filePath))
			res.AddHeader("Content-Type", mimetype)
			res.Write(data)
			return
		}
	}

	log.Println("NOT FOUND:", filePath)
	res.AddHeader("Content-Type", "text/html")
	res.Write(h.notFound)
}
