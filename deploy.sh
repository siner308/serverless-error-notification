#!/bin/bash

# deploy
cd terraform
terraform apply -parallelism=30 -auto-approve
cd ..
rm function.zip
rm main
