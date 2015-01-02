
beta: main.go
	docker run \
	-v $(PWD):/go/src/github.com/igneous-systems/beta \
	-v $(PWD)/../lib:/go/src/github.com/igneous-systems/lib \
	gotooling:0.0.1 \
	go build -o /go/src/github.com/igneous-systems/beta/beta \
		/go/src/github.com/igneous-systems/beta/main.go

#this depends on you having everything set right in external because we want
#to be able to try things out without committing them
hackygodep:
	docker run \
	-v $(PWD):/go/src/github.com/igneous-systems/beta \
	-v $(PWD)/../external:/external \
	-v $(PWD)/../lib:/go/src/github.com/igneous-systems/lib \
	-e "GOPATH=/go:/external" \
	-w /go/src/github.com/igneous-systems/beta \
	gotooling:0.0.1 \
	/go/src/github.com/igneous-systems/beta/inside-godep.sh


## run the gopherjs compiler
static/client.js: client/clientmain.go
	docker run \
	-v $(PWD):/go/src/github.com/igneous-systems/beta \
	-w /go/src/github.com/igneous-systems/beta \
	gotooling:0.0.1 \
	/go/src/github.com/igneous-systems/beta/gopherjs.sh build \
	-o static/client.js -m \
	github.com/igneous-systems/beta/client

