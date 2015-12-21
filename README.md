## Stevedore 

Stevedore is a helpful aid(e) in building and deploying Docker images to container registries, assuming your code is hosted in a git repository. 
It encompasses: 
 - A build script to build and tag images
 - A push script to build, tag, and then push images to a repository 
 - A standalone server to poll git repositories for changes, then build, tag, and push images to a repository 

## Using it

The easiest way to use the build & push scripts is to clone the source of this repository, navigate to the individual stevedore-build/ and stevedore-push/ directories, and then `go install .` on each.
That puts them in your $GOBIN (provided it's set), and you should be able to call them directly from there.

## Build script 

Inside the stevedore-build folder is the code for building images locally.  
Available command-line parameters: 
 - `-verbose`: enables verbose output
 - `-i [PATTERN]`: only build Dockerfiles with file extensions that match the regex [PATTERN], e.g. `stevedore-build -i dev`
 - `-registry-base [REGISTRY]`: tag builds with the appropriate registry prefix, e.g. `stevedore-build -registry-base docker.io/mygroup` 

When building, Stevedore automatically tags your built image with [REGISTRY]:[GIT_SHA].

## Push script

The push script (the stevedore-push folder) contains code for building as well as pushing.
It accepts all of the same arguments that the build script does. 

In addition to building and tagging with the GIT_SHA of the repo's current HEAD, it will also build and tag a secondary image with `latest`.

## Standalone server 

TODO: document this.
