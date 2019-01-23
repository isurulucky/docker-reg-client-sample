# docker-reg-client-sample

### docker-registry-file-uploader.go

This uploads a file to an unauthenticated docker registry running in localhost, port 5000.
Along with the file, it uploads a manifest which is a descriptor of the file system layer that was uploaded. 

### docker-registry-client.go

This code downloads the manifest and the relevant layers which are included the manifest. 

#### Note

This code uses the docker registry client implementation at https://github.com/nokia/docker-registry-client. 

