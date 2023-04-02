#!/bin/bash

sudo make build
version='1.0.0'

repo_name="w3ifun/stackup_bundler_go"
image_name="${repo_name}:${version}"
latest="${repo_name}:latest"
sudo docker rmi -f $(docker images -q ${repo_name})
sudo docker build --pull -t "$latest" -t "$image_name" .
sudo docker push "$repo_name" -a