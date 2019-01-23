package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/docker/distribution/manifest"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/libtrust"
	"github.com/nokia/docker-registry-client/registry"
	"github.com/opencontainers/go-digest"
	"io/ioutil"
	"os"
	"log"
)

func main() {
	url      := "http://localhost:5000"
	username := "" // anonymous
	password := "" // anonymous
	repository := "isuru/test"
	hub, err := registry.New(url, username, password)
	if err != nil {
		panic(err)
	}
	fr, err := os.Open("test.tar.gz")
	if fr != nil {
		defer fr.Close()
	}
	if err != nil {
		panic(err)
	}
	fileBytes, err := ioutil.ReadAll(fr)
	if err != nil {
		panic(err)
	}
	//digest := digest.NewDigestFromEncoded("sha256", hex.EncodeToString(fileBytes))
	h := sha256.New()
	h.Write(fileBytes)
	sha256sum := hex.EncodeToString(h.Sum(nil))
	digest := digest.NewDigestFromEncoded("sha256", sha256sum)
	log.Printf("Digest: %s", digest)
	exists, err := hub.HasBlob(repository, digest)
	if err != nil {
		panic(err)
	}
	if !exists {
		// read stream of files
		err = hub.UploadBlob(repository, digest, bytes.NewReader(fileBytes), nil)
		if err != nil {
			panic(err)
		}
		log.Printf("successfully uploaded the layer")
	} else {
		log.Printf("already exists")
	}

	manifest := &schema1.Manifest{
		Name: repository,
		Versioned: manifest.Versioned{
			SchemaVersion: 1,
		//	MediaType: schema1.MediaTypeManifest,
		},
		Tag: "latest",
		Architecture: "amd64",
		FSLayers: [] schema1.FSLayer{
			schema1.FSLayer{BlobSum: digest},
		},
		History: [] schema1.History{
			schema1.History{},
		},
	}

	key, err := libtrust.GenerateECP256PrivateKey()
	if err != nil {
		panic(err)
	}

	signedManifest, err := schema1.Sign(manifest, key)
	if err != nil {
		panic(err)
	}
	err = hub.PutManifest(repository, "latest", signedManifest)
	if err != nil {
		panic(err.Error())
	}
}