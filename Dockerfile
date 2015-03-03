FROM       debian:stable
ADD        ./stevedore /stevedore
ENTRYPOINT ["/stevedore"]
