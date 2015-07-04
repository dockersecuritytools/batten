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

## Running a Remote Check
Provide the '--server' flag to run a check on a remote Docker host.
Note that the remote host needs to be configured with TCP/TLS connection enabled.
In case you are using TLS you need to provide the certificates and key file as parameters
to batten command line:

```./batten --tlscacert=ca.pem --tlskey=key.pem --tlscert=cert.pem --server=tcp://<docker host>:<port> check```
