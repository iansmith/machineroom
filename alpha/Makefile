alpha: main.go ../lib/consul/consul.go
	docker run \
	-v $(PWD):/go/src/github.com/igneous-systems/alpha \
	-v $(PWD)/../lib:/go/src/github.com/igneous-systems/lib \
	-w /go/src/github.com/igneous-systems/alpha \
	gotooling:0.0.1 \
	/go/src/github.com/igneous-systems/alpha/go.sh build \
		-o /go/src/github.com/igneous-systems/alpha/alpha \
		/go/src/github.com/igneous-systems/alpha/main.go

#this depends on you having everything set right in external because we want
#to be able to try things out without committing them
hackygodep:
	docker run \
	-v $(PWD):/go/src/github.com/igneous-systems/alpha \
	-v $(PWD)/../external:/external \
	-v $(PWD)/../lib:/go/src/github.com/igneous-systems/lib \
	-e "GOPATH=/go:/external" \
	-w /go/src/github.com/igneous-systems/alpha \
	gotooling:0.0.1 \
	/go/src/github.com/igneous-systems/alpha/inside-godep.sh
