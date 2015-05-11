# Stevedore

> A simple CI server for building and publishing Docker images to a private registry

Stevedore automates the process of checking out new code, building it, and
creating a new Docker image.

## Overview

Stevedore listens for changes in a set of configured git repos. When a new
change is detected, Stevedore will:

* fetch the lastest commits on your repo's `master` branch
* call `make` to build the default `make` target (if your repo contains a `Makefile`)
* build a Docker image for each `Dockerfile` in the repo
* push the Docker image(s) to your configured Docker registry
* send success/failure notifications for each new image

## Details

In order for a repo to be `Stevedore`-compatible, it needs to adhere to the
following contraints:

* Code is available in a git repository.
* If this repo is not publicly-accessible/cloneable, the user account used to
  run the Stevedore process will need to be configured with access to the repo
(via `ssh` configs, or other configured auth).
* One or more `Dockerfile`s are present in the **root** of your repository.

## Image naming

Stevedore automatically names your Docker images based on the name of the
configured Docker registry, the name of the git repo, and the latest git commit
SHA.

In keeping with Docker's naming restrictions, only the characters `[a-z0-9-_.]`
are used.  Additionally, all forward-slash `/` characters after the git repo
hostname are replaced with the dash `-` character. This name is prepended with
the base URL of the container registry.

Stevedore will build an image for any Dockerfile named `Dockerfile` or with a
name that matches the glob `Dockerfile.*`. Any images built from a `Dockerfile` with an
extension will have that extension appended to the image name after an
additional dash `-` delimiter is added.

Stevedore does **not** tag any images as “latest”. We believe that code builds
and Docker images should be immutable artifacts that do not change over time,
and using "latest" as a build tag goes against this philosophy.  Instead, all
images are tagged instead with the first-seven characters from the HEAD git
commit SHA.

Some examples (using the [Google Container Registry
(gcr.io)](https://cloud.google.com/tools/container-registry/) as the Docker
registry:

| git repo url | Dockerfile | HEAD commit SHA | image name and tag |
| ----------------------------------- | ----------- | --------------- | ---------- |
| https://github.com/zulily/stevedore | Dockerfile | 2c1f4d8 | gcr.io/gce_project-name/zulily-stevedore:2c1f4d8 |
| https://github.com/zulily/stevedore | Dockerfile.bar | 1234567 | gcr.io/gce_project_name/zulily-stevedore-bar:1234567 |
| https://github.com/megacorp/biz-baz-buzz | Dockerfile | abcdef0 | gcr.io/gce_project_name/megacorp/biz-baz-buzz:abcdef0 |

## Notifications

Stevedore currently supports sending success/failure notifications via [Slack](https://slack.com/).  Other messaging platforms (HipChat, IRC, Campfire, etc.) can easily be added (pull requests welcome!).

An example success notification looks like:

![Sample Slack notification](https://github.com/zulily/stevedore/blob/master/slack.png)
