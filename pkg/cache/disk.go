package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Disk is a simple disk cache for responses that supports a crypto key and folder parameter
type Disk struct {
	UseCrypto   bool
	CacheKey    []byte
	CacheFolder string
}

// NewDisk creates an Disk for storing responses
func NewDisk(folder string, key string, crypto bool) (d *Disk) {
	d = new(Disk)
	d.UseCrypto = crypto
	d.CacheKey = []byte(key)
	d.CacheFolder = strings.TrimSuffix(folder, "/")
	return
}

// Fetch looks for the stored file and returns it.
func (d *Disk) Fetch(filename string) ([]byte, error) {
	if _, stat := os.Stat(filename); os.IsNotExist(stat) {
		// File doesn't exist return no error
		return nil, nil
	}

	bb, err := ioutil.ReadFile(filename)
	if err != nil {
		// Error! Failed to read file that exists!
		return bb, err
	}
	if d.UseCrypto {
		var dd []byte
		dd, err = Decrypt(bb, d.CacheKey)
		if err == nil {
			// Successfully decrypted!
			bb = dd
		} else {
			err = fmt.Errorf("cache failed to decrypt file '%s' : %s", filename, err)
		}
	}
	return bb, err
}

// Store will create a cache file with the bb bytes
func (d *Disk) Store(filename string, bb []byte) (err error) {
	filename = fmt.Sprintf("%s/%s", d.CacheFolder, filename)
	if d.UseCrypto && len(d.CacheKey) > 0 {
		bb, err = Encrypt(bb, d.CacheKey)
		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, bb, 0644)
	return err
}

// Clear will delete the cache file.
func (d *Disk) Clear(filename string) {
	_ = os.Remove(filename) // delete the cache file.
}
