# Instance setup for Stevedore

Stevedore runs on an instance at GCE.  The instance needs to be configured with
service account **read/write** access (the default is readonly) access to GCS,
in order for `stevedore` to be able to push new Docker images to Google
Container Registry (GCR).

The following steps are required to setup the instance after it is created. The
latest GCE-provided debian image for this setup.

#### Install the latest Docker build (v1.5 at time of writing)

    $ curl -sSL https://test.docker.com/ | sh
    $ /etc/init.d/docker start

#### Add an entry to /etc/hosts for the gitlab server

    $ echo "IP_ADDRESS github.com" >> /etc/hosts

#### Tell git (which stevedore shells out to) to ignore SSL cert warnings

    $ git config --global http.sslVerify false

#### Ensure that images can be pulled from the GCR registry

    $ gcloud components update
    $ gcloud preview docker pull gcr.io/YOUR_PROJECT/dcarney-actually-test:7052a3b

