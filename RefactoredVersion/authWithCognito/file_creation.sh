#!/bin/bash

# Create the base directories
mkdir -p internal/handler/auth internal/handler/user
mkdir -p internal/repo/auth internal/repo/user
mkdir -p internal/service/auth internal/service/user
mkdir -p internal/model

# Create handler/auth files
touch internal/handler/auth/concrete.go
touch internal/handler/auth/behaviour.go
touch internal/handler/auth/implementation.go

# Create handler/user files
touch internal/handler/user/concrete.go
touch internal/handler/user/behaviour.go
touch internal/handler/user/implementation.go

# Create repo/auth files
touch internal/repo/auth/concrete.go
touch internal/repo/auth/behaviour.go
touch internal/repo/auth/implementation.go

# Create repo/user files
touch internal/repo/user/concrete.go
touch internal/repo/user/behaviour.go
touch internal/repo/user/implementation.go

# Create service/auth files
touch internal/service/auth/concrete.go
touch internal/service/auth/behaviour.go
touch internal/service/auth/implementation.go

# Create service/user files
touch internal/service/user/concrete.go
touch internal/service/user/behaviour.go
touch internal/service/user/implementation.go

# Create model/model.go file
touch internal/model/model.go

echo "Files created successfully!"
