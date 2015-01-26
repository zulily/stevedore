package main

import (
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {

	serverCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			shutdown := make(chan bool)
			hostport := startWebServer(shutdown)
			startBuilder(shutdown, hostport)
			<-shutdown
		},
	}

	serverCmd.Flags().String("data", "/data", "Path to store repos and other data")
	viper.BindPFlag("data", serverCmd.Flags().Lookup("data"))

	serverCmd.Execute()
}

func startWebServer(shutdown chan bool) (hostport string) {
	container := restful.NewContainer()
	for _, resource := range NewAPIResources() {
		resource.Register(container)
	}

	port := "8080"

	config := swagger.Config{
		WebServices:    container.RegisteredWebServices(),
		WebServicesUrl: "localhost:" + port,
		ApiPath:        "/apidocs.json",

		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "/home/sreed/git/3rdparty/swagger-ui/dist",
	}
	swagger.RegisterSwaggerService(config, container)

	server := &http.Server{Addr: ":" + port, Handler: container}

	go func() {
		err := server.ListenAndServe()
		shutdown <- true
		log.Fatal(err)
	}()

	return "localhost:" + port
}

func startBuilder(shutdown chan bool, hostport string) {
	go func() {
		for {
			log.Println("Checking repos...")

			time.Sleep(5 * time.Second)
		}
		shutdown <- true
	}()
}
