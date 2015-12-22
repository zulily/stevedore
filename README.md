## Stevedore

Stevedore is a helpful aid(e) in building and deploying Docker images to container registries, assuming your code is hosted in a git repository.
It encompasses:
 - A build script to build and tag images
 - A push script to build, tag, and then push images to a repository

## Get it

`go get github.com/zulily/stevedore/...`

## Use it

```
> ~/mygroup/myproject/ $ ls
Dockerfile	Makefile	README.md	src/
> ~/mygroup/myproject/ $ stevedore-build
2015/12/22 09:52:25 > git config --get remote.origin.url
2015/12/22 09:52:25 > git rev-parse --show-toplevel
2015/12/22 09:52:25 > git rev-parse HEAD
2015/12/22 09:52:25 Building docker.io/mygroup/myproject:asdf000
...
```

#### Build script

Inside the stevedore-build folder is the code for building images locally.
Available command-line parameters:
 - `-verbose`: enables verbose output
 - `-i [PATTERN]`: only build Dockerfiles with file names that match the regex [PATTERN], e.g. `stevedore-build -i dev`.
 - `-registry-base [REGISTRY]`: tag builds with the appropriate registry prefix, e.g. `stevedore-build -registry-base docker.io/mygroup`

If the -i flag is omitted, you may optionally pass a list of Dockerfile names explicitly for stevedore to build, e.g.:  
`stevedore-build Dockerfile.dev Dockerfile.staging`  
This will only build the 'dev' and 'staging' Docker images for your project.

If neither the -i flag or Dockerfiles are provided, stevedore will build all images in the current working directory.

When building, Stevedore automatically tags your built image with {REGISTRY}/{GROUP}/{PROJECT}-{DOCKERFILE-SUFFIX}:{GIT_SHA}.

#### Push script

The push script (the stevedore-push folder) contains code for building as well as pushing.
It accepts all of the same arguments that the build script does.

After building and tagging, it pushes all images it builds to the tagged registry.

In addition to building and tagging with the GIT_SHA of the repo's current HEAD, it will also build and tag a secondary image per Dockerfile with `latest`.

