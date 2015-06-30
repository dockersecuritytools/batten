# batten - Docker Audit Toolkit

## About

`batten` down the hatches!  `batten` is an auditing framework that
contains some tools to help audit and harden your Docker deployments.

Identify potential security issues, and harden your existing Docker 
containers using a configurable policy.

## Docker Image
The supplied Dockerfile will help you create and run a Docker image.
Build a Docker image by running following command from the source root directory:

```docker build -t batten .```

After you create the Docker image, you can run it by supplying volume mount to the Docker socket file:

```docker run -v /var/run/docker.sock:/var/run/docker.sock batten```


