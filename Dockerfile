FROM       scratch
ADD        ./build-server /build-server
ENTRYPOINT ["/build-server"]
