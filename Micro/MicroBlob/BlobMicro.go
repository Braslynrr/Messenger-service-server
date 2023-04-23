package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

type config struct {
	port int
	env  string
}

type BlobMicro struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	ContainerName string
	client        *AzureClient
}

// Run sets up and runs the microserver
func (bm *BlobMicro) Run() (err error) {
	router := gin.Default()
	router.Use(cors.Default())
	micro := router.Group(bm.config.env)
	bm.infoLog.Println("Starting MicroService ...")
	micro.GET("/LoadBlob/:name", bm.LoadFile)
	micro.POST("/UpLoadBlob", bm.UploadFile)
	err = router.Run(fmt.Sprintf(":%d", bm.config.port))
	return
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 5000, "Server port to listen on")
	flag.StringVar(&cfg.env, "microblob", "blob", "url to blob microservice")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	accountName := os.Getenv("AZUREACCOUNTNAME")
	accountKey := os.Getenv("AZUREACCOUNTKEY")

	bm := &BlobMicro{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		client:   &AzureClient{Url: fmt.Sprintf("https://%s.blob.core.windows.net", accountName), connectionString: fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=core.windows.net", accountName, accountKey)},
	}

	bm.Run()

}
