package main

import (
	"log"
	"net/http"

	"core-gitlab.corp.zulily.com/core/build/api"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {

	serverCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			container := restful.NewContainer()
			r := api.NewRepoResource()
			r.Register(container)

			config := swagger.Config{
				WebServices:    container.RegisteredWebServices(),
				WebServicesUrl: "localhost:8080",
				ApiPath:        "/apidocs.json",

				SwaggerPath:     "/apidocs/",
				SwaggerFilePath: "/home/sreed/git/3rdparty/swagger-ui/dist",
			}
			swagger.RegisterSwaggerService(config, container)

			server := &http.Server{Addr: ":8080", Handler: container}
			log.Fatal(server.ListenAndServe())
		},
	}

	serverCmd.Flags().String("data", "/data", "Path to store repos and other data")
	viper.BindPFlag("data", serverCmd.Flags().Lookup("data"))

	serverCmd.Execute()
}
