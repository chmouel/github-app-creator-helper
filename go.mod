module github.com/chmouel/github-app-manifest-svc

go 1.16

require (
	github.com/google/go-github/v35 v35.1.0
	github.com/labstack/echo/v4 v4.2.2
	github.com/openshift/api v3.9.1-0.20190924102528-32369d4db2ad+incompatible
	github.com/openshift/client-go v0.0.0-20200116152001-92a2713fa240
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go v0.19.0
)

replace (
	k8s.io/apimachinery => github.com/openshift/kubernetes-apimachinery v0.0.0-20191211181342-5a804e65bdc1
	k8s.io/client-go => github.com/openshift/kubernetes-client-go v0.0.0-20191211181558-5dcabadb2b45
)
