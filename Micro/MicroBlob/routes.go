package main

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UploadFile handler to upload a file
func (bm *BlobMicro) UploadFile(ctx *gin.Context) {
	var err error
	var name string
	var url string
	var buffer []any
	data := map[string]any{}

	err = ctx.BindJSON(&data)

	name = data["name"].(string)
	buffer = data["buffer"].([]any)
	var bytes []byte = make([]byte, len(buffer))
	for index, value := range buffer {
		bytes[index] = byte(value.(float64))
	}
	if err == nil {

		url, err = bm.client.UpLoadBlob(name, bytes)

		if err == nil {
			ctx.IndentedJSON(http.StatusOK, gin.H{"url": url})
			return
		}

	}

	ctx.AbortWithError(http.StatusInternalServerError, err)
}

// LoadFile handler to load a file
func (bm *BlobMicro) LoadFile(ctx *gin.Context) {
	var err error
	var name string = ""
	if name = ctx.Params[0].Value; name != "" {
		img, err := bm.client.LoadBlob(name)
		if err == nil {
			ctx.Writer.Header().Set("Content-Type", *img.ContentType)
			image, err := io.ReadAll(img.Body)
			if err == nil {
				ctx.Writer.Write(image)
				ctx.Done()
				return
			}

		}

	}
	ctx.AbortWithError(http.StatusConflict, err)
}
