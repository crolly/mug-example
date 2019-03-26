# Example using mug - microservices understand golang
<p align="center"><img src="./logo.svg" width="480px" /></p>

This repository is an example for the mug generator following the Getting Started guide.
Two resources are created (Course and User).

To quickly test, I installed dynamodb-local, spun up the two tables and ran `make debug`. This creates the binaries in the `debug` subdirectory.
Afterwards a quick `sam local start-api --docker-network lambda-local` gives me a little local lambda magic, I can test against. The definitions for the API can be found in the `template.yml`.