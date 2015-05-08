FROM       scratch
ADD        ./stevedore ./stevedore
ENTRYPOINT ["./stevedore"]
