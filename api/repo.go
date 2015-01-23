package api

import (
	"net/http"

	"core-gitlab.corp.zulily.com/core/build/Godeps/_workspace/src/github.com/emicklei/go-restful"
)

// Repo represents a git source code repository.
type Repo struct {
	URL        string
	LastCommit string
}

// RepoResource provides functions for storing and retrieving Repo metadata
// from persistent storage.
type RepoResource struct {
	repos map[string]Repo
}

// NewRepoResource creates a new RepoResource.
func NewRepoResource() RepoResource {
	return RepoResource{map[string]Repo{}}
}

// Register creates a restful.WebService and configures API routes for managing
// Repos.
func (r RepoResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/repos").
		Doc("Manage Repos").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{repo-id}").To(r.findRepo).
		Doc("get a repo").
		Operation("findRepo").
		Param(ws.PathParameter("repo-id", "repo id").DataType("string")).
		Writes(Repo{}))

	container.Add(ws)
}

func (r RepoResource) findRepo(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("repo-id")
	repo := r.repos[id]
	if len(repo.URL) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: Repo could not be found.")
		return
	}
	response.WriteEntity(repo)
}
