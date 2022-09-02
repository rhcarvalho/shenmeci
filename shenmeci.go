package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"

	"github.com/rhcarvalho/shenmeci/internal/segmentation"
	"github.com/rhcarvalho/shenmeci/internal/shenmeci"
)

//go:embed static
var static embed.FS

func main() {
	shenmeci.LoadConfig()
	shenmeci.ValidateConfig()
	config := shenmeci.GlobalConfig

	static, err := fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}

	// Test whether we can listen on the provided Host and Port.
	// If the Host:Port is already in use, we can exit before wasting more resources.
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Http.Host, config.Http.Port))
	if err != nil {
		if err.(*net.OpError).Err.Error() == "address already in use" {
			os.Exit(0)
		}
		log.Fatal(err)
	}
	ln.Close()

	vocabulary := shenmeci.LoadCEDICT()
	shenmeci.LoadDB()
	defer shenmeci.CloseDB()
	shenmeci.Serve(&shenmeci.HTTPConfig{
		Host:      config.Http.Host,
		Port:      config.Http.Port,
		Static:    static,
		Segmenter: segmentation.NewSegmenter(vocabulary),
	})
}
