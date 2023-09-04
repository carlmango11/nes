install_dev:
	GOOS=js GOARCH=wasm go build -o frontend/public/nes.wasm backend/wasm/main.go

	if [ ! -d "frontend/node_modules" ]; then \
    	cd frontend; npm install --silent; \
    fi

build:
	docker buildx build --platform=linux/amd64 . --tag nes:latest
	#docker buildx build --platform=linux/amd64 . --tag 273011490881.dkr.ecr.eu-west-1.amazonaws.com/nes:latest