#!/bin/bash

# Define base structure
base="internal"
folders=("model" "repository" "service" "handler")
files_model=("model.go")
files_common=("behaviour.go" "implementation.go" "concrete.go")

# Create folders and files
mkdir -p "$base"

for folder in "${folders[@]}"; do
  mkdir -p "$base/$folder"
  if [ "$folder" == "model" ]; then
    for file in "${files_model[@]}"; do
      touch "$base/$folder/$file"
    done
  else
    for file in "${files_common[@]}"; do
      touch "$base/$folder/$file"
    done
  fi
done

echo "Folder and file structure created successfully."
