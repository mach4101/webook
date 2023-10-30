.PHYON: docker
docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -buildvcs=false -tags=k8s -o webook .
	@docker rmi -f mach4101/webook:v0.0.1
	@docker build -t mach4101/webook:v0.0.1 .
