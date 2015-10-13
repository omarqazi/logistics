logistics : 
	go build
test :
	go test ./...
fmt :
	go fmt ./...
run : logistics
	./logistics
clean :
	rm -rf logistics
image : test clean
	docker build -t omarqazi/logistics .
deploy : image
	docker push omarqazi/logistics