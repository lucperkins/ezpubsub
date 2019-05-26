emulator-env-init:
	$(gcloud beta emulators pubsub env-init)

start-emulator:
	gcloud beta emulators pubsub start

build:
	go build -v ./...

test:
	go test -v ./...

