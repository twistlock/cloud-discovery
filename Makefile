.PHONY: build

build:
	# Find all of the open source LICENSE files in the source code and include them in a /licenses directory inside the container image.
	# Inclusion of the license files may be a requirement for some to distribute the cloud-discovery container image.
	rm -rf licenses; mkdir -p licenses; find . | grep LICENSE | while IFS= read -r pathname; do filename=$$(echo $$pathname | tr "\.\/" "_" | sed "s/^__//"); cp $$pathname licenses/$$filename; done
	docker build -t twistlock/cloud-discovery .
