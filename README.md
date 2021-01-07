# maria

This is a frontend to [emily](https://github.com/antigravities/emily).

## Install

### Natively
```bash
git clone git@github.com:antigravities/maria.git
cd maria
go build
EMILY=http://localhost:3200 ./maria
#EMILY=http://localhost:3200 PORT=8000 ./maria #for a custom port
```

### Docker
```bash
# Create a network for maria and emily to share
docker network create maria-emily
# Start an emily container
docker run -it --name emily --network maria-emily -d antigravities/emily:1.1
# Start a maria container
docker run -it -p 3201:3201 -e EMILY=http://emily:3200/ --name maria --network maria-emily -d antigravities/maria:1.2
```

## License

GNU AGPL v3. See LICENSE for details.