logistics : 
	go build
test :
	go test . ./datastore ./auth
fmt :
	go fmt ./...
run : logistics
	./logistics
clean :
	rm -rf logistics
image : clean
	docker build -t omarqazi/logistics .
deploy : image
	docker push omarqazi/logistics