package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"sync"

	"golang.org/x/crypto/acme/autocert"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
)

type bucket struct {
	name string
	*storage.Client
	cert []byte
	mu   sync.RWMutex
}

func newBucket(ctx context.Context, name, project string) (*bucket, error) {
	c, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not create storage client")
	}
	return &bucket{Client: c, name: name}, nil
}

func (b *bucket) Get(ctx context.Context, key string) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.cert != nil {
		return b.cert, nil
	}
	r, err := b.Bucket(b.name).Object(key).NewReader(ctx)
	if err == storage.ErrObjectNotExist {
		return nil, autocert.ErrCacheMiss
	} else if err != nil {
		return nil, errors.Wrap(err, "could not create reader on bucket")
	}
	return ioutil.ReadAll(r)
}

func (b *bucket) Put(ctx context.Context, key string, data []byte) error {
	w := b.Bucket(b.name).Object(key).NewWriter(ctx)
	if _, err := io.Copy(w, bytes.NewBuffer(data)); err != nil {
		return errors.Wrap(err, "could not write to object")
	}
	b.mu.Lock()
	b.cert = data
	b.mu.Unlock()
	return nil
}

func (b *bucket) Delete(ctx context.Context, key string) error {
	if err := b.Bucket(b.name).Object(key).Delete(ctx); err == storage.ErrObjectNotExist {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "could not delete object")
	}
	return nil
}
