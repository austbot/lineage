# Lineage
Enforce docker image ancestry policies.


# Contributing
Instal Dep package manager for go.
`https://github.com/golang/dep`
Quickly you can install it with 
```
go get -u github.com/golang/dep/cmd/dep
dep ensure
go run Main.go scan-file -dockerfile DockerfileTest -whitelist testwhitelist.txt
```

