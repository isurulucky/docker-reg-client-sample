package main

import (
	"github.com/nokia/docker-registry-client/registry"
	"io/ioutil"
	"log"
)

func main()  {
	url      := "http://localhost:5000"
	username := "" // anonymous
	password := "" // anonymous
	repository := "isuru/test"
	hub, err := registry.New(url, username, password)
	if err != nil {
		panic(err)
	}
	tags, err := hub.Tags(repository)
	if err != nil {
		panic(err)
	}
	log.Printf("tags: %+v", tags)
	manifestv1, err := hub.Manifest(repository, "latest")
	if err != nil {
		panic(err)
	}
	log.Printf("v1 manifest: %s", manifestv1)
	//manifestv2, err := hub.ManifestV2(repository, "latest")
	//if err != nil {
	//	panic(err)
	//}
	//log.Printf("v2 manifest: %s", manifestv2)
	digest, err := hub.ManifestDigest(repository, "latest")
	if err != nil {
		panic(err)
	}
	log.Printf("digest: %s", digest)
	// get the fs layers
	for _, element := range manifestv1.References() {
		reader, err := hub.DownloadBlob(repository, element.Digest)
		if reader != nil {
			defer reader.Close()
		}
		if err != nil {
			panic(err)
		}
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		// write
		log.Println("Going to write to file... ")
		ioutil.WriteFile(element.Digest.Hex() + ".tar.gz", bytes, 0644)
	}
}
