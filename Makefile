default: run

run:
	go run main.go

fmt:
	go fmt ./...

api:
	curl http://localhost:2023/-/play/compile -X POST --data '{"body":"a=5 + 5"}' -H "content-type:application/json"
	curl http://localhost:2023/-/play/fmt?body="a=1" -X POST -H "content-type:application/json"
