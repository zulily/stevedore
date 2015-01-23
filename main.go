package main

import (
	"log"
	"net/http"

	"core-gitlab.corp.zulily.com/core/build/api"
	"core-gitlab.corp.zulily.com/core/build/Godeps/_workspace/src/github.com/emicklei/go-restful"
	"core-gitlab.corp.zulily.com/core/build/Godeps/_workspace/src/github.com/emicklei/go-restful/swagger"
)

func main() {
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
}
