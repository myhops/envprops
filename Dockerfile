# Dockerfile for pre-compiled Go binary using Chainguard static base image
# syntax=docker/dockerfile:1

# ARG TARGETPLATFORM is automatically provided by buildx and used to select the correct base image.
# For example, 'linux/amd64' or 'linux/arm64'.
ARG TARGETPLATFORM

# Use the Chainguard static base image for the target platform.
# It runs as a non-root user "nonroot" (uid 65532) by default.
# Choose the variant that matches your Go binary's linkage (glibc or musl).
# For standard CGO_ENABLED=0 Go builds on Linux, 'latest-glibc' is typical.
FROM --platform=$TARGETPLATFORM cgr.dev/chainguard/static:latest-glibc

# If you were specifically building against musl (e.g., if your previous target was Alpine for the binary itself),
# you might use: cgr.dev/chainguard/static:latest-musl

# The 'static' image's default user is 'nonroot' and WORKDIR is '/'.
# We'll set a specific WORKDIR for clarity.
WORKDIR /app

# Copy the pre-compiled binary from the Docker build context.
# GoReleaser places the binary named in 'builds.binary' (e.g., "f12" in your case)
# into the context for the corresponding platform.
# This binary should be statically compiled (CGO_ENABLED=0).
COPY f12 .

# Expose port if your application listens on one (e.g., 8080)
# EXPOSE 8080

# Command to run the application.
# The entrypoint path is relative to the WORKDIR /app
ENTRYPOINT ["/app/f12"]