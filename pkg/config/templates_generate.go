// Code generated by vfsgen; DO NOT EDIT.

// +build release

package config

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	pathpkg "path"
	"time"
)

// TemplateFolder statically implements the virtual filesystem provided to vfsgen.
var TemplateFolder = func() http.FileSystem {
	fs := vfsgen۰FS{
		"/": &vfsgen۰DirInfo{
			name:    "/",
			modTime: time.Date(2019, 1, 19, 14, 35, 44, 23013607, time.UTC),
		},
		"/default.test.tiogo.yaml": &vfsgen۰CompressedFileInfo{
			name:             "default.test.tiogo.yaml",
			modTime:          time.Date(2018, 12, 15, 15, 23, 48, 0, time.UTC),
			uncompressedSize: 340,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\xcf\x4d\x4b\xc3\x40\x10\xc6\xf1\xfb\x7e\x8a\x07\x22\xf6\x52\xf2\x52\xdf\xf7\xa6\x55\x51\xaa\x08\x09\x7a\x5f\x37\x93\xec\xea\x9a\x09\xbb\x63\x4b\xbe\xbd\x34\x14\xbd\x7b\x7d\xe0\xff\x1b\x26\xcb\xb0\xe6\x71\x82\x38\x82\x0d\x9e\x06\x41\x22\x2b\x9e\x07\x18\x1c\x3d\xbc\x3c\xdf\xdd\x3e\xd6\xe8\x7c\x20\x2c\xf2\x9e\x47\x47\x31\x9f\xcc\x57\x58\xc0\x0c\x2d\x78\x4b\x31\xfa\x96\x54\x96\x61\xe7\xc5\xcd\xcc\xb5\xb5\x94\xd2\x86\xa6\xe3\x86\x6c\x24\xd9\xd0\xa4\xd6\x33\xad\x15\x70\x63\x12\xbd\xd6\x4f\x1a\x4e\x64\xd4\x45\x11\xd8\x9a\xe0\x38\x89\xae\xca\xaa\xac\x14\xfe\x7a\x0d\xf3\x6e\xab\xd5\xc9\xb2\x77\xfe\xe2\xf2\x4a\x01\xbf\xa0\x46\x4b\xdd\xe9\xd9\xf9\xf2\xe3\x33\x94\xd5\x4a\xa9\x86\xe2\x96\xa2\xfe\x67\x0e\xd4\xcc\x72\xcf\xa1\xa5\xa8\x91\x17\x96\x87\xce\xf7\x45\xcb\x36\x32\x4b\xa1\xd4\x1b\xc5\xe4\x79\xd8\xfb\x8d\xe3\x5d\x43\xfb\x63\x90\xf8\x4d\x87\xe5\xf0\xe0\x3c\xfd\x04\x00\x00\xff\xff\x0e\xd9\x61\xb1\x54\x01\x00\x00"),
		},
		"/default.tiogo.yaml": &vfsgen۰CompressedFileInfo{
			name:             "default.tiogo.yaml",
			modTime:          time.Date(2019, 1, 19, 4, 40, 38, 838739600, time.UTC),
			uncompressedSize: 471,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x90\x41\x6b\xb3\x40\x10\x86\xef\xfb\x2b\xde\x8f\xef\x9a\x66\xb5\xd0\xcb\xdc\x6c\xa2\x41\x9a\x18\xd0\xb5\x3d\x6f\xd6\x29\x8a\xab\x5b\xdc\xb5\xd0\xfe\xfa\xa2\x04\x72\x28\xbd\x65\x98\xcb\x3b\x3c\xcc\x3c\xcc\xff\xbb\x94\x58\x1b\x6f\x49\x59\xe4\xc5\x81\xb0\x4f\xb3\xa4\x3e\x2a\xec\xce\x45\x96\x1f\xea\x32\x51\xf9\xb9\xb8\x42\xa7\xba\x52\xa8\x52\x45\x48\x8c\x61\xef\x5f\xf8\x0b\x7a\x6c\x50\xb1\x99\x38\x2c\x69\x98\x7d\xc0\x85\xe1\x39\x20\x38\xe8\x15\x83\xe2\x51\x5f\x2c\x6f\x3b\xf7\x4f\xdc\xc7\x7a\x11\x7f\x3d\x91\x00\x9e\xb5\xe7\xba\x3c\x12\xda\x10\x3e\x48\x4a\xeb\x8c\xb6\xad\xf3\x81\xe2\x28\x8e\x62\x81\x9b\x2b\x41\xf7\xf1\x46\xf7\x8f\x02\x37\x65\x82\xef\xe3\x8d\x5f\x87\x3b\x6d\x5a\xce\x9c\x6d\x78\x22\x6c\xcd\x92\xa4\xb1\x1d\x8f\x41\x0a\x60\xcf\xef\x7a\xb6\x41\x75\x03\x7f\xbb\x91\x09\x0f\xd1\x53\x14\x21\xad\x94\x10\x15\x4f\x9f\x3c\xfd\xf2\xf1\x24\xa5\xb1\x6e\x6e\xb6\xe1\xfa\x02\xe3\x86\x3f\x0e\xf9\x75\x85\x14\x3f\x01\x00\x00\xff\xff\x20\xc6\x1d\x64\xd7\x01\x00\x00"),
		},
		"/template": &vfsgen۰DirInfo{
			name:    "template",
			modTime: time.Date(2019, 1, 15, 3, 6, 34, 52840687, time.UTC),
		},
		"/template/client": &vfsgen۰DirInfo{
			name:    "client",
			modTime: time.Date(2019, 1, 15, 3, 5, 53, 568285188, time.UTC),
		},
		"/template/client/table.tmpl": &vfsgen۰CompressedFileInfo{
			name:             "table.tmpl",
			modTime:          time.Date(2019, 1, 8, 4, 46, 51, 0, time.UTC),
			uncompressedSize: 3130,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x96\x5f\x6f\xda\x30\x10\xc0\xdf\xf9\x14\x27\x0f\xfa\x47\x25\x71\xd9\xd6\x97\x4a\x3c\x4c\x42\x9a\x78\xd9\x53\x1f\x91\xae\x59\x71\x42\xd6\x36\x89\x12\x8b\x69\xb2\xf3\xdd\xa7\xb3\x43\xa1\xc5\x0e\x21\x50\x07\x39\xc1\xf7\xc7\x67\x72\xbf\x3b\x38\x87\x95\x94\xc5\x3d\xe7\x45\x24\xf3\xf2\xcf\x73\xf8\x94\xbf\xf2\x2a\x8f\xe5\xdf\xa8\x14\x5c\x46\x51\xc2\xbf\x14\xd3\x65\x5a\x15\x2f\xd1\xbf\x8b\x78\x3a\x5f\x97\xa9\xbc\x90\xd3\x24\x2f\x56\xa2\x1c\xdd\xfe\x18\x28\xb5\x14\x71\x9a\x09\x60\x3f\xcd\xda\x43\xf4\xfb\x45\xb0\xba\x1e\xdc\xf4\x1a\x27\x98\x81\x63\xe0\xde\x8a\xcb\x0c\x91\x14\x11\x69\xa2\x2f\x1a\x34\x3d\x9b\x15\xb3\xe0\x34\xe3\x80\x8f\xa0\x39\x20\x2c\x34\x5c\xbe\xcd\x60\x56\xe0\x12\x51\xbb\xcc\x34\x5c\xa1\x36\xf3\xb5\xd9\xc7\xcc\xe6\x02\x44\x4e\x37\xf7\x6e\x0b\xc4\x31\xe8\x05\x22\x72\x0d\x21\xcd\xe4\x06\xcd\x0a\x3d\x7a\xcc\x34\x19\xd8\x93\x6f\xb4\x3a\xfd\x24\x1d\xc6\x07\xb3\xa0\xc7\xb8\x19\x28\x05\x45\x99\x66\x32\x06\xa6\x61\x14\x7c\xaf\xf4\x28\x98\xdc\xd2\xfc\xed\xae\xd2\xc0\x80\xcd\x67\x0c\xd8\xaf\xe8\x55\x30\x60\x33\x51\x3d\x95\x69\x21\xd3\x3c\x63\x40\x39\xd6\x77\xd3\x32\xca\x12\x01\xc3\xe7\x31\x0c\xd7\x70\x3f\x85\x10\x82\xba\x1e\xd0\xa9\xda\x03\x5a\x64\x0c\xc2\xf9\x0c\x42\x0a\x08\xc2\x9d\x78\x28\x1c\xa5\x02\x10\xd9\xd2\xf8\xea\x1d\x9a\xc8\x96\x75\x3d\xd8\xe1\xea\x61\x95\x66\x49\x4f\xac\x7a\xa9\x5b\x6c\xd0\x09\x90\x4b\xdd\xe4\xa1\x41\xe6\x0a\xaf\x1b\x5e\x2c\x52\x1e\x75\x42\xc3\x02\xb3\x45\xe6\x11\x5c\xb8\x6c\xbd\xbf\xbb\x2c\x42\x4e\x75\x02\xa5\x41\x03\xb7\x88\x8c\xbd\xea\xfb\x63\x83\x4c\x47\xf5\x33\xe0\xe0\xc2\x00\x6c\xda\xc1\x28\xf8\x7a\x57\xe9\x16\x0e\x8e\xcd\xb5\xdd\xf4\x4f\x28\xff\x93\x03\x00\x80\x7e\x17\xca\xa7\x20\xe0\x4c\xfc\xa6\xa1\x50\xfa\x57\x27\xb6\x15\xff\xf8\x1c\x87\xae\x04\xf9\x88\x92\x17\x32\x8f\xc3\x43\x5d\xca\x8f\xa1\xdb\x61\x87\xfe\xe5\x03\xd5\xe9\xf0\x70\x67\xf3\xa3\xec\x8e\xf0\x70\xcf\xf3\xc2\xee\x76\xd8\xa5\x1b\x7a\xca\x41\xd7\xb7\xdc\x79\xec\x39\x3c\x06\xe2\x6e\x44\x79\x3b\x99\xa5\xb9\xa9\x32\x13\x23\x98\xb4\x16\x99\x46\x62\x58\x0c\xcf\xd9\x8c\xdb\x0f\xe0\xab\x53\x5b\x89\x24\x89\x34\x12\x5b\x28\x1a\xf9\x91\x47\xa7\xa2\x36\x5c\x27\x54\xd7\xe8\x66\x4a\x1b\x3d\xec\x56\xb7\xe1\x5a\x5a\xb9\xdc\xc8\x65\x5b\xf5\x53\x2a\x8d\x21\xcb\xe5\x26\xb0\xba\x56\xea\xfc\x41\x31\xd6\x7c\xc8\xbd\xad\xa0\x27\xff\x0b\x69\x7f\x27\x76\x93\xff\x01\x00\x00\xff\xff\x7d\xd8\xe2\xa8\x3a\x0c\x00\x00"),
		},
		"/template/cmd": &vfsgen۰DirInfo{
			name:    "cmd",
			modTime: time.Date(2019, 1, 19, 14, 44, 42, 473217881, time.UTC),
		},
		"/template/cmd/server": &vfsgen۰DirInfo{
			name:    "server",
			modTime: time.Date(2019, 1, 19, 14, 44, 42, 473217881, time.UTC),
		},
		"/template/cmd/server/server.tmpl": &vfsgen۰FileInfo{
			name:    "server.tmpl",
			modTime: time.Date(2019, 1, 17, 14, 30, 22, 480579891, time.UTC),
			content: []byte(""),
		},
		"/template/cmd/tiogo.tmpl": &vfsgen۰CompressedFileInfo{
			name:             "tiogo.tmpl",
			modTime:          time.Date(2019, 1, 19, 14, 20, 47, 136626632, time.UTC),
			uncompressedSize: 1268,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x94\x41\x4f\xe3\x3a\x10\xc7\xef\xfe\x14\xf3\xd0\x93\xb8\x34\xad\xde\xe3\xbd\x0b\xd2\x1e\xd8\x22\x58\x0e\x10\xa4\x14\xf6\x10\xf5\xe0\x24\x93\xd4\x5a\xc7\x13\xec\x71\xa0\xaa\xfa\xdd\x57\x8e\x9b\xa6\x42\xbb\x0b\x39\xc4\x19\xcf\xfc\xfe\x19\x8f\xc7\xde\xed\x2a\xac\x95\x41\x38\x7b\x72\xb2\xc1\x33\x48\xf6\x7b\x71\xce\x8a\xe6\x0d\x9d\x83\x72\x20\x0d\x28\xc3\x68\x6b\x59\x62\xf8\x22\x58\xa1\x91\x85\xc6\xf9\x5d\x0a\x57\x8f\x77\xe0\x9d\x32\x0d\xdc\xd2\x5c\x88\x1b\x65\x2a\x68\xc9\x86\xc0\x9a\x6c\x2b\x59\x91\x01\xc9\x97\x02\x00\x60\xc3\xdc\xb9\xcb\xc5\xa2\x51\xbc\xf1\xc5\xbc\xa4\x76\xf1\xba\x41\x8b\xca\xfd\xf0\x96\x17\xac\xa8\xa1\x85\x10\x43\x1e\x91\x60\x45\x90\x2f\xd3\xfb\xfb\xab\x87\xeb\x35\xe4\xd9\xd3\xd7\xc9\xb8\x5a\xae\xee\xd2\x07\x98\xcf\xe7\x6b\xc8\xd3\xc7\x60\x64\x6b\x21\x96\xd4\xb6\xd2\x54\x2e\x0a\xf4\x2d\xc4\x67\x9c\x86\x9a\xec\xe9\x02\x9e\xbd\x36\x68\x65\xa1\xb4\xe2\x2d\xdc\x4b\x23\x1b\x6c\xd1\x30\xe4\x15\xd6\xd2\x6b\x9e\x41\x29\x0d\x14\x08\xd4\x2a\x66\xac\xd6\x83\xae\x43\xdb\xa3\x7d\xaf\xab\xa9\x94\x1a\x3a\x4b\x6f\x5b\x90\xa6\x82\x6f\xab\xd5\xe3\x18\xaa\x8c\x63\x69\x4a\x1c\xf0\xd2\x7d\x90\xd6\x92\x0c\x4b\x65\xd0\x42\x86\xa5\xb7\x21\xb7\xbc\xf6\xec\x2d\xfe\x15\x13\x78\xc5\x42\x76\xdd\x1f\x14\xbe\xc7\x80\xac\x94\x26\xc8\x4c\xb4\xc8\x7c\x91\x94\x91\x1a\x8b\x14\x47\x80\x0d\xea\x6e\x06\x2e\x32\x6e\x06\xb2\x41\xc3\xe3\x98\x34\x96\x7c\xe7\xa2\x3f\x4c\x3a\x87\xc1\x89\x6f\x1d\x59\x4e\xde\x99\xbd\xd7\x21\xc8\x3b\xb4\xf1\x7d\xc4\x59\xda\x06\x47\x35\x71\x52\xce\x31\x0b\xc7\xd2\xf2\x0c\x1c\x53\x27\x0e\xe5\x1a\x5d\xb9\x21\x83\x90\x40\x5c\x0d\x74\x5a\x1a\x17\xd6\x34\x95\xe4\xa3\xc8\x5b\x4d\x85\xd4\x90\x76\xa1\x35\x0f\xba\xcf\x68\x0b\x72\x8a\xb7\x23\x9c\x24\x4e\x69\x34\x3c\x03\x48\xe2\x56\x65\xc8\xa0\xa9\x69\x94\x69\x16\xe4\xb9\xf3\x0c\x1a\x7b\xd4\x90\x0f\xc3\x3f\xeb\x23\xf9\xe2\x15\x06\x10\x92\x97\xcf\x90\xff\x4e\x64\x38\x33\xb3\xe1\xb3\xff\x0c\x79\x91\x1c\x5a\x74\x52\xa8\xb0\xf0\xcd\x0c\x8e\xcf\x07\x0a\xff\x4d\x24\x5b\x59\xe2\x09\x99\xc6\x68\x26\xc8\x56\xd7\xe9\xd3\x6a\xe8\x67\xa6\xa0\x06\xb5\xd2\x78\x50\xf8\x7f\x52\x18\xec\x2f\x17\xa7\xff\x76\xc0\x1b\x84\xc3\x8f\xfb\xb1\xc8\x87\x14\x8c\x6f\xd1\xaa\x52\x6a\xbd\x3d\x1e\xb6\xb5\x10\x37\x64\xe3\x0d\x12\x7a\x31\xee\xc7\xdf\xc3\x3d\xd0\xb7\xc3\xd4\xb1\x3b\x7f\xe1\x8a\xed\xfa\x1b\xc6\x09\x21\x76\x3b\x34\xd5\x7e\xff\x33\x00\x00\xff\xff\x3b\x97\x83\xb2\xf4\x04\x00\x00"),
		},
		"/template/cmd/vm": &vfsgen۰DirInfo{
			name:    "vm",
			modTime: time.Date(2019, 1, 19, 14, 44, 35, 977143187, time.UTC),
		},
		"/template/cmd/vm/vm.tmpl": &vfsgen۰CompressedFileInfo{
			name:             "vm.tmpl",
			modTime:          time.Date(2019, 1, 19, 14, 13, 34, 143606479, time.UTC),
			uncompressedSize: 3453,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x56\xdd\x4f\x23\x37\x10\x7f\xdf\xbf\x62\xb4\x54\xca\x43\xb3\x9b\x04\xee\x4e\xea\x49\x91\x4a\xef\xb8\x13\x95\x48\x50\x03\xf4\xaa\xc0\x83\x63\x4f\x36\x86\x5d\x7b\xf1\x47\x12\x8a\xf8\xdf\x2b\xdb\xbb\x21\xd9\x7c\x94\x5e\xdb\x07\xc8\x7a\xbe\x3c\x33\xfe\xcd\xcf\x7e\x7e\x66\x38\xe5\x02\x21\x9e\x17\xd7\x9a\x64\x18\x43\xf2\xf2\x12\xb5\x0c\x97\x69\x26\x5b\xc0\x35\x10\x01\x5c\x18\x54\x53\x42\xd1\x7d\x49\xb8\x42\x41\x26\x39\xa6\xe7\x43\x38\xbd\x3c\x07\xab\xb9\xc8\xe0\xab\x4c\xa3\xe8\x0b\x17\x0c\x0a\xa9\x9c\xe1\x54\xaa\x82\x18\x2e\x05\x10\xf3\x31\x02\x00\x98\x19\x53\xea\x8f\x9d\x4e\xc6\xcd\xcc\x4e\x52\x2a\x8b\xce\x62\x86\x0a\xb9\x7e\xb0\xca\x74\x0c\x97\x99\xec\x44\x91\xcf\x23\x78\x18\x2e\x61\x5e\xc0\x78\x74\xfd\xcb\xa7\xe1\xc5\xc5\xe9\xe0\xf3\x1d\x8c\x4f\x3f\x5d\x9d\x0f\x07\x90\xa6\xe9\x1d\x8c\x87\x97\x6e\x31\xba\x8b\xa2\x91\x9d\x24\x54\x16\x05\x11\x4c\x57\xfb\x61\x5e\xb6\x41\x53\x22\x04\x2a\xdd\x06\x92\xa1\x30\xf5\x6f\x92\x29\x69\x4b\x1d\xf4\x4e\xa8\x35\x3a\x25\x2e\x4b\xa9\x4c\xd2\x58\xce\x6d\xee\x8c\xac\x46\x15\xfe\xaf\xdc\x0d\x51\x19\xd6\xd1\xa2\xe8\x94\xba\x92\xc3\xfe\x39\xd7\xa6\x0d\x54\x21\x31\xd8\x06\x86\x39\xba\x5f\xfd\x24\x68\x1b\x14\x6a\x3a\x43\x66\x73\x74\x21\xb2\x36\x58\x61\x48\x16\x45\x37\x17\x30\x2c\x5d\x88\xaa\x86\x11\xe6\xe8\x43\x42\x21\x19\x9f\x72\x54\x95\x22\x49\x38\xeb\x8f\xad\xe0\x8f\x16\x81\xb3\xbb\x4a\x28\x48\x81\xfd\xb1\x36\x8a\x8b\xac\x96\x29\xcc\x70\xd9\x1f\x2b\xcc\x6c\x4e\x94\x2b\x49\xa1\xd6\x5c\x8a\xda\xe0\xfe\xd1\xe9\xef\x1f\x37\x54\xd1\xd0\x9a\xd2\x1a\xb8\x90\x0c\x57\x7b\x52\x3d\xf7\x49\x19\x30\x0e\x02\x20\xbd\x8d\x06\x23\xc1\xf7\x1e\x34\x96\x44\x11\x83\x0c\xa6\x3c\x47\x0d\x63\x8e\x29\x64\x52\x32\x98\x4a\x05\x67\x4b\x8a\x39\xfc\x08\xa3\x32\xb7\xe2\xa1\x0d\x68\x68\xba\x4a\x42\x4b\xb1\x27\xf4\xaf\xa3\xe1\xa0\x11\xc9\x41\x32\x53\x1e\x5f\x0e\xa2\x0c\xee\x1f\xa1\x20\x82\x97\x36\x0f\xc2\xf4\xce\x37\x33\x9c\x47\x30\x39\x5b\x92\xa2\xcc\xeb\x62\x8e\x8e\x60\x54\x41\x03\xa8\x14\xcc\x52\x03\x02\xcd\x42\xaa\x07\x98\x10\x8d\x2c\x20\xc7\x61\xdb\x39\x2b\x74\x40\x70\xe7\x66\x73\xa3\x61\x42\xe8\x03\x6c\xcc\x42\x5a\x47\xad\x61\x4e\x73\x69\x59\x6a\x2a\x03\x87\x76\x52\xf2\xce\x51\x47\xa1\x96\x56\x51\xd4\x9d\x1a\x9a\xde\xf3\x07\x8f\xf6\x1d\xa2\x79\xf1\x37\x52\x0f\xb4\xa8\xde\xfe\xd4\xa3\x1c\x88\x1f\x42\x6d\x48\x9e\x23\x83\x82\xd0\x19\x17\xa8\xc1\xcc\x88\x01\x5c\x22\xb5\x06\x61\x26\xb5\x2f\xa5\x2e\x56\xaf\x57\xfa\x9f\x54\x18\x46\x6e\x2d\xed\x3d\x02\x5f\x02\x24\xc9\x82\x9b\x99\xb4\xd5\x38\xed\xb5\xf2\x5a\x88\x07\xb8\x80\xaf\xfe\x73\x40\x0a\x8c\x57\x1d\x08\x32\x57\xca\xb5\x08\xa6\x55\x4f\x3c\x7b\xf9\xef\x60\xb3\x23\x8f\x60\xbf\x6f\x8b\x7a\x98\x20\xee\x75\x6f\xd3\xee\x6d\x7a\x7c\x1b\x6f\xc7\xb0\xe2\x70\x94\xe6\x51\x85\x5c\x20\x59\x3f\x3a\x37\x79\x3e\xdd\x4a\xe9\xaa\xe1\xc6\xf1\x71\x6b\xbd\x80\x56\x38\x50\xe7\x71\xe5\x89\x48\x83\x9c\x56\x25\x3a\x74\xeb\xef\x3d\xb3\x9a\xd0\x36\x00\xb7\xae\x09\xa0\x6b\xd4\xbe\xa2\xc1\xfd\x6e\x81\x0e\x2b\xa6\x82\x78\xad\x18\x38\xfa\x23\x3e\xe0\x18\xf8\x73\xaf\x63\x5d\x67\xc5\x5a\xe7\x97\x40\x18\x73\x64\x86\xaf\x2d\xf1\x64\xa2\x25\x2c\x10\x28\x11\x75\x2a\xa1\x71\xdb\x88\xd8\x55\x30\x24\x09\xd1\xc9\x06\xe7\x43\x92\x68\xfe\x27\xf6\x7b\x1f\x0e\xf8\x3a\xd2\x3f\xec\xbb\x3d\xda\x5b\x2d\xde\x90\xbb\x46\xe0\x1c\x55\xa2\xac\xd8\xe9\xfc\x7a\xc3\x40\x92\x48\x81\x86\x17\xee\x8b\xb9\x92\x63\x81\x4b\x03\x1a\x1d\xef\xc1\xef\xc8\x04\x6a\x46\x9e\x1c\xbc\xbd\x55\xdc\x7b\xd7\xed\xae\xc1\xd4\xdf\x85\xdf\x01\xa4\x57\xbf\xd7\xae\x78\xd9\x1a\x63\x35\x35\x86\x64\x90\x24\xd6\x72\xd6\x8f\xbb\xbd\xe3\x93\xc4\xfd\x9d\x1c\xbf\x7b\x9f\x1c\xf7\x4e\x7e\x7a\xef\x72\xa4\xc4\x60\x26\xd5\x13\xc4\xc3\x85\x40\xe5\xd3\x26\x19\xc4\x5f\xb8\x20\x82\x62\x13\x44\x21\xae\xbf\x5f\xff\x75\xe4\xd0\x04\x38\xf3\xaf\x02\x3f\xb2\x21\xba\x1b\xcf\x1b\xf7\x46\xa8\x0c\xde\xde\xa4\xf0\xc0\x68\x74\x69\xe3\x11\x72\x40\xe5\x52\x96\x62\x8e\xca\x24\x46\x26\x94\xd0\x19\x46\x3b\xcd\xfd\xfb\x65\x75\xa0\xd7\x1a\xd5\xa6\x99\x7b\xd2\xec\x80\x9b\x17\x37\x26\xf6\x37\x24\xf9\x8a\x0d\xb1\x20\x3c\x87\x58\x21\x71\xca\x9f\x99\x2c\x08\x17\xa9\xc9\x59\xbc\x23\x50\x63\x82\x5d\x16\xf5\x00\x7f\x8b\x37\x92\xf3\xd2\x14\x06\xd2\x00\x77\x57\x76\x81\xc2\x70\x91\xb5\x61\x62\x0d\x10\x50\x58\x70\xc1\x50\xc1\x62\xf6\xe4\x86\x99\x49\xd1\x32\x40\x49\x9e\x57\xec\x58\x93\xe3\xbd\xd5\x06\x5a\xd5\x6a\x3b\xa3\x3d\x54\xd6\xd0\x36\xea\x6f\xa4\xbd\xd7\xed\x6d\xd5\x6e\x90\x8f\xe7\x71\xeb\x6e\x63\x2e\xaa\x6b\x38\x40\x8b\x11\x3d\x83\x89\x24\x8a\xb9\x4b\x6c\xeb\x4a\x5e\xed\xbf\x49\x2e\x8d\xc4\x3f\x9f\xdd\xc0\xa0\x7a\xe0\x1c\x7d\x73\xa7\x57\x60\x31\x41\xa5\xfb\x71\xaf\x9b\xf6\xd2\x6e\xda\xed\xf4\x3e\xb4\xa1\xd7\x4d\x8f\xab\x45\x7c\x28\xbc\x46\xf3\xbf\xc5\x26\x65\x89\x82\xbd\x31\xfc\xc9\x7a\xf8\x77\x6f\x09\xaf\xb0\x90\xf3\xb7\x76\xe6\x9f\x87\x6f\x1c\x7d\x23\x7c\x14\x3d\x3f\xa3\x60\x2f\x2f\x7f\x05\x00\x00\xff\xff\x9f\x1d\xb9\x1d\x7d\x0d\x00\x00"),
		},
		"/vfsgen_templates.go": &vfsgen۰CompressedFileInfo{
			name:             "vfsgen_templates.go",
			modTime:          time.Date(2019, 1, 19, 14, 35, 44, 19013561, time.UTC),
			uncompressedSize: 1323,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\x53\x41\x6f\xf3\x36\x0c\x3d\x5b\xbf\x82\x13\xf0\x21\x0e\x96\x59\xf7\x00\xbd\xec\xfb\xda\xa2\xc0\xd0\x0e\x58\xb6\x1d\x0b\x45\xa6\x65\x2d\xb2\x64\x48\x74\x83\x60\xe8\x7f\x1f\x28\x3b\x29\xe2\xf9\x64\x89\xe4\xe3\x7b\x8f\x94\x52\xf0\xb7\xf3\x1e\x2c\x06\x4c\x9a\x10\x34\x6c\x08\x87\xd1\x6b\xc2\xfc\x7e\xbd\x6d\x6c\xdc\xc0\xd9\x51\x0f\xda\x7b\x88\x1d\x50\x8f\xd0\x39\x8f\x19\xa6\xd0\x62\x02\xea\x5d\x86\x2e\xfa\x16\x93\x50\x0a\x0e\x7c\x74\x19\x02\x1a\xcc\x59\xa7\x0b\x1c\xd1\xe8\x29\x33\xfc\xd1\x05\xbe\x18\x53\xb4\x49\x0f\x60\x74\x80\x23\x42\x9a\x02\x74\x29\x0e\xa0\xc3\xe5\xdc\x63\x42\x88\xe1\xab\xcb\x25\x13\x72\xa8\x65\xf0\x41\x5f\x20\x44\x82\x5e\x7f\x30\x5e\x42\xaf\xc9\x7d\xe0\xd2\x1e\x36\x8d\x32\x31\x74\xce\xaa\xab\x0e\xb5\x69\x00\xfe\xcc\x2e\x58\xf8\xe8\xb2\xc5\x00\x67\x04\x93\x70\x96\x9b\x49\x93\x33\x60\x63\x69\xc5\x0d\x8a\x50\x13\x03\x61\xa0\x7c\x55\x7b\x33\x05\x70\x38\x62\xdb\x62\xdb\xc0\x4d\x67\x1b\x03\xce\x65\xc7\xc9\xf9\x16\x48\xdb\xdc\x88\x51\x9b\x93\xb6\x08\x83\x76\x41\x08\xa5\x6c\xdc\xdf\x6c\xb6\xb1\x48\x9e\xf9\xbc\xdf\xc0\x1b\x1b\x85\x70\xc3\x18\x13\x41\x2d\x2a\x69\x1d\xf5\xd3\xb1\x31\x71\x50\xb9\x9f\x92\x89\xf1\x37\x35\xd7\xc8\x55\xd4\xa5\x69\xcc\x18\x94\x8f\x36\x4d\x99\xa3\x01\x49\xf5\x44\x23\xff\xc7\x72\x93\x29\xb9\x60\xb3\x14\x5b\x21\xba\x29\x98\xc2\xac\xde\xc2\xbf\xa2\x8a\x13\x8d\x13\x3d\x39\x8f\x41\x0f\x08\xfb\x07\x90\xe3\xc9\xae\x9d\xbc\xdb\x08\x29\x44\xa5\x14\xbc\xbe\x1d\x1e\xf7\xf0\xd2\xb1\xa9\xb7\x29\xb2\x63\x2f\x3f\x1e\x97\x9d\x81\xe4\x6c\x4f\xbf\x18\xef\xcc\x09\xe2\x94\xc0\x9c\xdb\x7a\xcb\xc6\xb9\x90\x5d\x8b\xec\xf1\xba\x93\xa8\xcc\xb9\xdd\xc1\x3b\x53\x89\xb9\x79\x46\xe2\x1a\x51\x5d\x13\x9e\xe6\x69\x33\xd1\xa5\x94\xf9\xcc\xea\x9b\x97\xd0\xc5\xae\x96\xcf\x33\x59\x1e\x7c\x61\x75\x8e\xe9\x54\x0e\xa5\x76\xcf\x34\xf6\xdf\xb2\xdc\xf1\xcf\x76\x56\xf3\xbd\x47\x73\x02\xc7\x6a\x36\xa9\x08\x0a\x5c\xb1\xf0\x64\x59\x2b\xa2\x0b\xd8\x8e\xb7\x13\x74\xfb\xcf\x94\xe9\x6b\x25\x47\x4d\x7d\x6e\x44\xe5\x3a\x58\xbc\x6f\xbe\xc7\x40\xda\x85\x5c\x17\x79\x92\x5c\xb4\x71\xb1\x59\x96\x49\xac\x15\x3e\x80\x6c\x94\x14\xd5\x7a\x44\x7c\xdf\x28\x09\x3f\xc3\x7d\x40\x54\x9f\xff\x33\xe2\x0f\xa4\xe2\xc2\x8a\xf4\x7e\xf3\x2d\x6f\x0a\xf1\x19\xa2\xbc\x00\xc6\x28\x01\xb9\x83\x7b\x2a\xbb\x55\xa7\xad\x10\xa2\xc2\x54\x86\x30\x2f\x65\xb3\x38\x8e\x35\x2f\x5e\xf3\xc3\xa5\xfa\x1e\x62\xbb\xbb\x66\xbe\x8d\xe4\x62\xc8\x2c\xf8\x8a\xb7\x07\xfe\xee\x7b\xec\x44\x55\xfd\x3e\xbf\xa3\xd7\x39\x65\x19\xb7\xe4\xc8\xaf\xfc\xda\x0e\xda\xe6\x52\x2a\x13\x7a\xd4\x19\x4b\xe8\x2f\x9d\x9c\x3e\xfa\xa5\x4a\x1e\xee\x68\x70\xc6\x27\x4f\xdc\x75\xc0\x02\x7e\x7a\x80\xe0\x7c\x71\x7f\x31\xee\x49\x93\xf6\x3e\xd4\x98\xd2\xb6\x38\xfa\x29\xfe\x0b\x00\x00\xff\xff\xe6\x28\xa2\xb7\x2b\x05\x00\x00"),
		},
	}
	fs["/"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/default.test.tiogo.yaml"].(os.FileInfo),
		fs["/default.tiogo.yaml"].(os.FileInfo),
		fs["/template"].(os.FileInfo),
		fs["/vfsgen_templates.go"].(os.FileInfo),
	}
	fs["/template"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/template/client"].(os.FileInfo),
		fs["/template/cmd"].(os.FileInfo),
	}
	fs["/template/client"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/template/client/table.tmpl"].(os.FileInfo),
	}
	fs["/template/cmd"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/template/cmd/server"].(os.FileInfo),
		fs["/template/cmd/tiogo.tmpl"].(os.FileInfo),
		fs["/template/cmd/vm"].(os.FileInfo),
	}
	fs["/template/cmd/server"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/template/cmd/server/server.tmpl"].(os.FileInfo),
	}
	fs["/template/cmd/vm"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/template/cmd/vm/vm.tmpl"].(os.FileInfo),
	}

	return fs
}()

type vfsgen۰FS map[string]interface{}

func (fs vfsgen۰FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {
	case *vfsgen۰CompressedFileInfo:
		gr, err := gzip.NewReader(bytes.NewReader(f.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &vfsgen۰CompressedFile{
			vfsgen۰CompressedFileInfo: f,
			gr:                        gr,
		}, nil
	case *vfsgen۰FileInfo:
		return &vfsgen۰File{
			vfsgen۰FileInfo: f,
			Reader:          bytes.NewReader(f.content),
		}, nil
	case *vfsgen۰DirInfo:
		return &vfsgen۰Dir{
			vfsgen۰DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// vfsgen۰CompressedFileInfo is a static definition of a gzip compressed file.
type vfsgen۰CompressedFileInfo struct {
	name              string
	modTime           time.Time
	compressedContent []byte
	uncompressedSize  int64
}

func (f *vfsgen۰CompressedFileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰CompressedFileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰CompressedFileInfo) GzipBytes() []byte {
	return f.compressedContent
}

func (f *vfsgen۰CompressedFileInfo) Name() string       { return f.name }
func (f *vfsgen۰CompressedFileInfo) Size() int64        { return f.uncompressedSize }
func (f *vfsgen۰CompressedFileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰CompressedFileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰CompressedFileInfo) IsDir() bool        { return false }
func (f *vfsgen۰CompressedFileInfo) Sys() interface{}   { return nil }

// vfsgen۰CompressedFile is an opened compressedFile instance.
type vfsgen۰CompressedFile struct {
	*vfsgen۰CompressedFileInfo
	gr      *gzip.Reader
	grPos   int64 // Actual gr uncompressed position.
	seekPos int64 // Seek uncompressed position.
}

func (f *vfsgen۰CompressedFile) Read(p []byte) (n int, err error) {
	if f.grPos > f.seekPos {
		// Rewind to beginning.
		err = f.gr.Reset(bytes.NewReader(f.compressedContent))
		if err != nil {
			return 0, err
		}
		f.grPos = 0
	}
	if f.grPos < f.seekPos {
		// Fast-forward.
		_, err = io.CopyN(ioutil.Discard, f.gr, f.seekPos-f.grPos)
		if err != nil {
			return 0, err
		}
		f.grPos = f.seekPos
	}
	n, err = f.gr.Read(p)
	f.grPos += int64(n)
	f.seekPos = f.grPos
	return n, err
}
func (f *vfsgen۰CompressedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.seekPos = 0 + offset
	case io.SeekCurrent:
		f.seekPos += offset
	case io.SeekEnd:
		f.seekPos = f.uncompressedSize + offset
	default:
		panic(fmt.Errorf("invalid whence value: %v", whence))
	}
	return f.seekPos, nil
}
func (f *vfsgen۰CompressedFile) Close() error {
	return f.gr.Close()
}

// vfsgen۰FileInfo is a static definition of an uncompressed file (because it's not worth gzip compressing).
type vfsgen۰FileInfo struct {
	name    string
	modTime time.Time
	content []byte
}

func (f *vfsgen۰FileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰FileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰FileInfo) NotWorthGzipCompressing() {}

func (f *vfsgen۰FileInfo) Name() string       { return f.name }
func (f *vfsgen۰FileInfo) Size() int64        { return int64(len(f.content)) }
func (f *vfsgen۰FileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰FileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰FileInfo) IsDir() bool        { return false }
func (f *vfsgen۰FileInfo) Sys() interface{}   { return nil }

// vfsgen۰File is an opened file instance.
type vfsgen۰File struct {
	*vfsgen۰FileInfo
	*bytes.Reader
}

func (f *vfsgen۰File) Close() error {
	return nil
}

// vfsgen۰DirInfo is a static definition of a directory.
type vfsgen۰DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *vfsgen۰DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *vfsgen۰DirInfo) Close() error               { return nil }
func (d *vfsgen۰DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *vfsgen۰DirInfo) Name() string       { return d.name }
func (d *vfsgen۰DirInfo) Size() int64        { return 0 }
func (d *vfsgen۰DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *vfsgen۰DirInfo) ModTime() time.Time { return d.modTime }
func (d *vfsgen۰DirInfo) IsDir() bool        { return true }
func (d *vfsgen۰DirInfo) Sys() interface{}   { return nil }

// vfsgen۰Dir is an opened dir instance.
type vfsgen۰Dir struct {
	*vfsgen۰DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *vfsgen۰Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *vfsgen۰Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}