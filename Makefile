default: run

run:
	go run main.go

fmt:
	go fmt ./...

test:
	go test ./...

api:
	curl http://localhost:2023/-/play/compile -X POST --data '{"body":"a=5 + 5"}' -H "content-type:application/json"
	curl http://localhost:2023/-/play/fmt?body="a=1" -X POST -H "content-type:application/json"

image:
	docker build . -t kcllang/kcl-playground
	docker push kcllang/kcl-playground
