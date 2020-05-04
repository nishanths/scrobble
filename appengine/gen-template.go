// Code generated by go-bindata. DO NOT EDIT.
// sources:
// appengine/template/content.html (543B)
// appengine/template/root.html (478B)
// appengine/template/u.html (466B)
// appengine/helpguide.md (2.774kB)

package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _appengineTemplateContentHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x92\x31\x4f\xc3\x30\x10\x85\xf7\xfe\x8a\xe3\xf6\xd6\xea\xc6\x60\x67\x29\x20\x36\x10\x0a\x03\xa3\xb1\xaf\xd8\xc2\xb1\xa3\xf8\x94\xaa\xaa\xf2\xdf\x51\xe2\x14\xa5\x05\x21\xb1\x24\x77\xf7\xde\x7d\xd2\x3b\x59\xde\xdc\x3d\xed\xea\xb7\xe7\x7b\x70\xdc\x84\x6a\x25\xcf\x3f\xd2\xb6\x5a\x01\x00\xc8\x86\x58\x83\x71\xba\xcb\xc4\x0a\x5f\xeb\x87\xf5\x2d\x82\x58\x8a\x51\x37\xa4\xb0\xf7\x74\x68\x53\xc7\x08\x26\x45\xa6\xc8\x0a\x0f\xde\xb2\x53\x96\x7a\x6f\x68\x3d\x35\x38\xef\xb1\xe7\x40\xd5\xe9\x04\x9b\x7a\xac\x60\x18\xa4\x28\xb3\xa2\x07\x1f\x3f\xa1\xa3\xa0\x30\xf3\x31\x50\x76\x44\x8c\xe0\x3a\xda\x2b\x14\x99\x35\x7b\x23\xf6\x29\xf2\xf4\xd9\x98\x9c\xf1\x1f\x8b\x26\x67\xd1\xa5\x74\xde\x93\xa2\xa4\x95\xef\xc9\x1e\x67\x8c\xf5\x3d\x78\xab\x50\xb7\xed\x4c\xfe\x1e\x9b\xa0\x73\x56\xf8\x92\x12\x2f\xa4\x1f\xfa\xc8\xf4\xf1\xe3\xda\x32\xba\xdc\x76\x4a\xfe\x58\x1c\x53\x76\xb7\xbd\xf4\x49\x61\x7d\xff\x07\x7d\x3e\xf1\x2f\xf4\x91\xbc\x2b\x2a\x0c\xc3\x15\xe1\x92\xba\x68\xe7\x52\x8a\x72\x03\x29\xca\x3b\xf8\x0a\x00\x00\xff\xff\x64\x3a\xd1\x4b\x1f\x02\x00\x00")

func appengineTemplateContentHtmlBytes() ([]byte, error) {
	return bindataRead(
		_appengineTemplateContentHtml,
		"appengine/template/content.html",
	)
}

func appengineTemplateContentHtml() (*asset, error) {
	bytes, err := appengineTemplateContentHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "appengine/template/content.html", size: 543, mode: os.FileMode(0664), modTime: time.Unix(1588593833, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x3e, 0x36, 0x53, 0x9c, 0xb7, 0xb8, 0x70, 0x33, 0xa3, 0x54, 0x4e, 0x71, 0x1c, 0xd8, 0xf9, 0x73, 0x49, 0xfd, 0x19, 0x1a, 0xcf, 0x4b, 0x8a, 0xc6, 0x72, 0x95, 0x8d, 0xbf, 0x15, 0xb3, 0xcd, 0xfc}}
	return a, nil
}

var _appengineTemplateRootHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x91\x3f\x6f\xf3\x20\x10\xc6\xf7\x7c\x8a\x7b\xd9\x63\xd6\x57\x2a\x78\xe8\xbf\xb5\x1d\xdc\xa1\x23\x81\x8b\x7c\xa9\x0d\x16\x77\x72\x14\x59\xfe\xee\x95\x43\xac\x3a\xdd\xca\x00\x3c\x1c\x3f\x9e\x47\x9c\xf9\xf7\xfc\xf6\xd4\x7c\xbe\xbf\x40\x2b\x7d\x57\xef\xcc\xba\xa0\x0b\xf5\x0e\x00\xc0\xf4\x28\x0e\x7c\xeb\x32\xa3\x58\xf5\xd1\xbc\xee\xff\x2b\xd0\xdb\x62\x74\x3d\x5a\x35\x12\x9e\x87\x94\x45\x81\x4f\x51\x30\x8a\x55\x67\x0a\xd2\xda\x80\x23\x79\xdc\x5f\x85\xba\x71\x42\xd2\x61\x3d\x4d\x55\xb3\x6c\xe6\xd9\xe8\x72\x52\xaa\x1d\xc5\x2f\xc8\xd8\x59\xc5\x72\xe9\x90\x5b\x44\x51\xd0\x66\x3c\x5a\xa5\x59\x9c\x90\xd7\xc7\x14\xe5\x3a\x55\x9e\x59\xfd\x01\xf4\xcc\x3a\x26\x21\xdc\x82\xec\x33\x0d\x52\xc4\x32\x7c\x8a\x2c\x70\x48\x49\x58\xb2\x1b\xc0\xc2\x34\x55\x8f\xab\x9c\xe7\x87\x82\xe9\x95\x33\xba\xfc\x98\x39\xa4\x70\xb9\xbd\x19\x68\x04\x0a\x56\xb9\x61\x50\xb5\xd1\x81\xc6\x3b\x33\xe0\xec\x7f\x52\x9d\xd6\x50\x3d\xc5\xea\xc4\x0b\xb0\xcd\x74\xcf\x04\x62\xd1\x39\x25\xf9\x75\xd3\xe8\x62\x6f\x74\x69\xe3\x77\x00\x00\x00\xff\xff\x47\x84\x67\x47\xde\x01\x00\x00")

func appengineTemplateRootHtmlBytes() ([]byte, error) {
	return bindataRead(
		_appengineTemplateRootHtml,
		"appengine/template/root.html",
	)
}

func appengineTemplateRootHtml() (*asset, error) {
	bytes, err := appengineTemplateRootHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "appengine/template/root.html", size: 478, mode: os.FileMode(0664), modTime: time.Unix(1588580441, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x17, 0x73, 0x75, 0x3d, 0xa2, 0x2c, 0xb3, 0xc8, 0xb6, 0xbf, 0x50, 0xd0, 0xf6, 0x1f, 0x47, 0x46, 0x3, 0xeb, 0x10, 0x56, 0x70, 0x44, 0x7b, 0x2a, 0x8b, 0x28, 0x87, 0x18, 0xd8, 0x66, 0x76, 0xa6}}
	return a, nil
}

var _appengineTemplateUHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x51\x3d\x6f\xf3\x20\x10\xde\xf3\x2b\xee\x65\x4f\x58\x5f\xa9\xe0\xa5\x1f\x6b\x3b\xa4\x43\x47\x0a\x97\xf8\x52\x1b\x2c\xee\xe2\x28\xb2\xfc\xdf\x2b\x07\x47\xb5\x3b\x95\x01\x38\xee\xf9\x12\x67\xfe\x3d\xbd\x3e\xee\x3f\xde\x9e\xa1\x96\xb6\xa9\x36\xe6\x7e\xa0\x0b\xd5\x06\x00\xc0\xb4\x28\x0e\x7c\xed\x32\xa3\x58\xf5\xbe\x7f\xd9\xfe\x57\xa0\x97\xcd\xe8\x5a\xb4\xaa\x27\xbc\x74\x29\x8b\x02\x9f\xa2\x60\x14\xab\x2e\x14\xa4\xb6\x01\x7b\xf2\xb8\xbd\x15\x6a\xe6\x09\x49\x83\xd5\x30\xec\xf6\xd3\x65\x1c\x8d\x2e\x2f\xa5\xdb\x50\xfc\x82\x8c\x8d\x55\x2c\xd7\x06\xb9\x46\x14\x05\x75\xc6\x83\x55\x9a\xc5\x09\x79\x7d\x48\x51\x6e\xdb\xce\x33\xdf\x65\xd9\x67\xea\xa4\x14\xd3\xf2\x29\xb2\xc0\xd9\xe5\x23\x83\x85\x61\xd8\x8d\xe3\x43\x41\xea\x25\xf4\x4f\x86\x9e\x59\xc7\x2e\xa7\x63\x46\xe6\xd9\xd4\xe8\xf2\x51\xe6\x33\x85\xeb\xac\x15\xa8\x07\x0a\x56\xb9\xae\x53\x95\xd1\x81\xfa\x55\x38\xe0\xec\x7f\x44\x4f\x4b\xcd\x13\x4f\x84\x55\xb0\x15\x27\x10\x8b\x3e\xff\x82\x19\x5d\xbc\x8d\x2e\xa3\xfb\x0e\x00\x00\xff\xff\x60\xd4\xfe\x41\xd2\x01\x00\x00")

func appengineTemplateUHtmlBytes() ([]byte, error) {
	return bindataRead(
		_appengineTemplateUHtml,
		"appengine/template/u.html",
	)
}

func appengineTemplateUHtml() (*asset, error) {
	bytes, err := appengineTemplateUHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "appengine/template/u.html", size: 466, mode: os.FileMode(0664), modTime: time.Unix(1588580445, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x36, 0x16, 0x98, 0xd9, 0xdf, 0x67, 0x1e, 0xaa, 0x80, 0xd4, 0x8d, 0xbb, 0x62, 0x4, 0xb7, 0x2e, 0xdd, 0xe7, 0x17, 0x5b, 0x78, 0x9d, 0x37, 0x8c, 0xe9, 0xe6, 0x49, 0xfa, 0xb1, 0x19, 0xd0, 0xf2}}
	return a, nil
}

var _appengineHelpguideMd = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x56\x5d\x6f\x1c\xb7\x0e\x7d\x9f\x5f\xc1\xac\x1f\xe2\x0d\xf6\xce\x26\xce\x05\x2e\xe0\xdc\xd4\x48\xdc\x8f\x04\x68\x1b\x03\x36\x1a\x04\x86\x91\x68\x25\xee\x8c\x60\x8d\xa8\xea\x63\x07\xd3\x5f\x5f\x50\x9a\x99\xec\xae\x9d\x36\x4f\xf6\x4a\x22\x79\x78\x78\x48\xce\xc9\x09\x7c\x6c\x45\x04\x1d\x20\xb6\x3a\x5c\x54\xd5\xe7\x6b\xe9\x69\xb3\x31\xf8\x99\x0f\x05\x74\x29\x68\x09\xa1\x1c\x6a\xdb\x40\x40\xbf\xd3\x12\x21\xb2\x5d\x4f\xfe\x3e\x40\xaf\x63\x0b\x6f\x9c\x33\x08\xbf\xf1\xf3\xba\xaa\x6e\x5a\x84\x88\xbe\xdb\xb7\xf4\xb8\x45\x1f\x20\x12\xc4\x16\x41\xc8\x08\xb4\x05\x97\x36\x46\x87\x96\xef\xf9\x34\x90\x6d\x02\x0c\x94\xc0\xe8\x10\xd1\xe6\xd7\x04\xc2\x02\x59\xa3\x2d\xc2\x16\x51\x8d\xfe\x27\x24\x92\x6c\xd0\x21\x06\x76\x17\x7b\x02\x27\x7c\x0c\xe7\x55\xf5\xa2\x66\xfc\x42\x7e\xb8\x06\xe1\x9c\xd1\x52\x44\x4d\x16\xfa\x56\xcb\x16\x92\x33\x24\x54\x0e\xe5\xf7\xb1\x8f\x71\x19\x4f\xab\x43\x24\x3f\xbc\x02\x61\x55\x75\x56\x67\x86\xa0\xc7\xcd\x23\xde\x94\x0e\xce\x88\x61\x74\xc7\x2e\x18\xcc\x98\x3b\xaa\x92\x56\x5d\x55\x27\x27\xf0\x8e\x7a\x50\x94\x53\x6c\x30\x42\x88\xc2\x47\x54\x17\x55\x75\x43\xb0\xc1\x46\xdb\x03\xb2\xd9\xae\xa4\xb2\xd3\x41\xc7\xcc\xd1\x08\x01\x44\x84\x36\x46\x17\xce\xd7\xeb\x29\x52\x2d\x8c\x41\x83\xb5\x94\xaf\xd8\x26\xe8\xc6\x82\xb6\xa5\x40\xbf\x10\x35\x06\x39\x19\x70\x5a\xde\x17\xa8\x29\xa0\xb7\xa2\xc3\xfc\x5c\x92\x1b\x72\x84\x37\x57\xef\xe1\x1e\x87\x7c\xa8\xa8\xb7\x4c\x55\xbe\xb8\x9d\xe9\xbc\xe3\x7f\x29\xdc\xe5\x37\x68\x23\xfa\xfc\x40\x92\xd3\xa8\x26\x07\x1c\x9b\x4f\x67\xab\xba\xaa\x3e\x71\xd8\x23\x66\x40\x0a\x0b\x1b\xae\x28\x5a\x4e\xeb\xcb\xb7\xf3\x5a\xa7\xf5\x84\xf9\x4b\x21\x94\xa5\xe0\x3c\x6d\xb5\x41\x70\xa2\xc1\x22\x8e\xfd\x93\xaf\xe5\x79\x58\x91\x4f\x94\x72\xf4\xad\x36\x9c\xc3\x66\x80\xc5\x1b\x63\x16\x70\x7a\x54\x58\x61\xcc\xb1\xf5\x72\x05\x8b\x5f\x69\x87\xea\xe1\xeb\x92\x95\xd1\xf7\xa8\x80\xec\xbe\xbe\x96\x2b\x20\x0f\x8b\xb7\x03\x48\x32\xe4\x1f\x31\x3d\xa2\x66\x33\x80\xf0\x91\x1b\xad\x58\x2c\xeb\xaa\x7a\x72\x1b\xb0\xe9\x98\x75\xc5\xea\x8f\x9e\xcc\xdd\x29\x73\x76\xbe\x5e\xeb\xae\x49\xbe\x96\xd4\xad\x2f\xff\x3a\xdb\xbd\x3c\xfb\x9f\xa9\x9d\x6d\x96\x33\x55\xa5\x16\x1d\xda\x04\x1b\xe1\xf7\xb5\x5c\x88\x9b\x6b\xb5\x1a\xb5\x3d\x01\xfa\x9e\x5e\x59\x4d\x95\x9c\x54\x83\x0a\x6e\x5b\xf4\x38\xc9\xa5\xae\xaa\x4b\xc3\xea\xa3\xa2\x0c\xd6\xb1\x96\x64\x67\xa5\x8c\xb8\x56\x59\xa7\xb2\x25\x0a\x08\x8b\x6b\x6e\x92\xbd\xbe\xa8\xeb\x7a\x51\xc3\x4f\x59\x75\x05\xd4\xa1\xde\x94\x16\x86\x1a\x1e\x1c\xc7\x1d\x55\x57\xd5\xff\x75\xd7\x40\xf0\xf2\xf5\x62\x52\xd9\x57\xca\xfe\xb8\xf9\xf0\x73\xfb\xfc\x8c\x19\x5b\x40\x8b\xba\x69\xe3\xeb\x97\xcf\x9f\xc3\xfa\x87\xcc\xdf\xa5\xb0\xb9\x71\x27\x4a\x72\x63\x51\x8a\xf3\x9c\x91\xd4\xb9\x14\xd1\x5f\x64\x2e\x3d\xf2\x0c\xb5\x04\x28\xc2\x00\xbd\x18\x40\x26\xef\xd1\x46\x33\x30\xb6\xd9\x4b\x4e\xa1\x54\xfb\x5b\x0e\x6b\xb8\x32\x28\x02\xc2\xad\xf4\x28\x22\xb7\x31\xe8\x10\x12\xde\xdd\xe6\x3f\xe1\x2e\x8f\x49\xa5\x0e\xa6\x47\x72\x8e\x7c\xe4\x1b\x4b\xf6\x3f\x87\x1e\x27\xf1\x77\x62\x58\x41\x4b\x3d\xee\xd0\xaf\x78\x1e\x64\x06\xdf\xdd\xdc\x5c\x65\x56\x1f\xc7\x59\x67\x29\xf1\x7d\x4e\x30\x82\x22\x99\x46\x41\x0e\x18\xe1\xf4\xc9\x72\x05\x9b\x54\xe6\xd5\x6d\xa0\xe4\xf3\x9c\x56\x78\x37\xfe\xb8\xcb\xdb\x65\x27\xb4\x11\xdc\xde\x33\x14\x08\x51\x1b\x33\xc3\x38\x9e\xb6\x5c\x51\x4f\x7d\x98\x16\x45\x5e\x42\x19\xd6\xd6\x6b\xb4\x2a\x40\x2b\x76\xf8\xb5\x87\xca\x84\xf8\xc8\x2e\xcb\x22\x3b\xdc\x2c\x64\x81\x62\x8b\x1e\x14\xf2\x26\x09\x70\x8a\x4d\x0d\xfa\xaa\x25\x8b\xcb\x3c\x90\x26\x47\x17\x55\xf5\x09\xc3\x0a\xf4\x36\xbb\x98\xf0\x05\xd1\xe1\x41\x47\x08\x29\x29\xd9\xc8\x9e\x33\xac\xe2\x2b\x8b\x79\x3a\x3a\x2e\xec\x35\x8e\xbe\x50\xe6\x1c\x37\x68\xa8\x87\x2d\x31\xaa\x28\xb4\x09\x6c\xd9\x52\x3f\xee\x1f\xde\xb8\x25\xad\x1f\xd1\x88\x2c\xf9\x71\x76\x70\xc9\xcb\xfe\x11\x11\x43\x2c\x14\xb1\x10\x75\x28\xbb\xfc\xde\x52\x3f\xca\x06\x54\xc2\xac\x18\x30\xba\xd3\xb1\xb0\xcb\x9b\xa2\x45\x7b\xe0\xe4\x41\x8b\x1f\x14\x6e\x6a\xb8\x7d\x0a\x14\x86\xfb\x48\xae\xac\x28\x5b\xf6\x86\xbe\x49\x16\x43\x3e\x3a\x75\x1e\x47\x29\x5e\x8a\x28\x8c\xb6\x62\x39\x2e\x74\x95\x13\x52\xe8\x72\x25\xc9\x82\x80\x2d\xf6\xb0\x15\x32\x92\xe7\x35\xf8\xec\x41\xb0\x79\x5a\xe5\x50\x7b\x61\xa2\xb8\x47\xf0\x94\x9a\xd6\x0c\xf0\xe2\xac\xe5\x64\x93\x53\xdc\x3a\xcc\x15\xe4\x32\x71\x41\xad\x34\x49\x71\x7e\x86\xa4\x30\xfb\x97\xaf\xaa\x67\x33\x1f\x4f\x77\x08\x46\xf0\x5a\x77\x68\x51\xfd\x03\x8e\xd3\x69\xac\x75\x29\x44\xd6\x10\x5b\x3c\x1e\x7d\x59\x3e\x2d\x9e\xed\xed\xc8\xb9\xd7\x84\x73\x4f\xc3\xf8\x95\x02\x5b\x8f\x7f\x26\xb4\x72\xe0\xaf\x8a\xb3\xff\xb6\xa5\xfe\xef\x43\x86\x54\x6f\x3b\x5e\x29\x41\x77\xda\x08\x3f\x37\xfd\x38\x60\x3d\x09\xd5\x09\x77\x51\x55\xbf\x53\x0d\x6f\x53\x04\xf7\x6f\x83\xa4\x68\xfc\x29\x4f\x2f\x9e\xb0\x18\x62\x6e\xa4\xe9\xfe\x7c\xfe\xea\x68\x74\x6c\xd3\x26\x0f\x4e\xab\x43\x2b\x6c\x6c\xc3\xbc\xb2\xd7\xe5\xf9\xda\x62\x5f\x4d\x4d\xff\x9d\xa6\xd5\xb8\x2c\xbe\x37\x92\xc7\x9c\x52\x58\x17\xcd\x56\x7f\x07\x00\x00\xff\xff\x60\x1a\xd1\x4a\xd6\x0a\x00\x00")

func appengineHelpguideMdBytes() ([]byte, error) {
	return bindataRead(
		_appengineHelpguideMd,
		"appengine/helpguide.md",
	)
}

func appengineHelpguideMd() (*asset, error) {
	bytes, err := appengineHelpguideMdBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "appengine/helpguide.md", size: 2774, mode: os.FileMode(0664), modTime: time.Unix(1588595568, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xea, 0x34, 0x82, 0xc4, 0xbc, 0x27, 0x9d, 0xc0, 0xe9, 0xc3, 0x2d, 0x24, 0x34, 0xa7, 0x86, 0x2d, 0x71, 0x12, 0x29, 0x58, 0xf1, 0xcd, 0x59, 0xbd, 0x51, 0x6d, 0xd5, 0x7f, 0x19, 0xe4, 0xc, 0xfb}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"appengine/template/content.html": appengineTemplateContentHtml,
	"appengine/template/root.html":    appengineTemplateRootHtml,
	"appengine/template/u.html":       appengineTemplateUHtml,
	"appengine/helpguide.md":          appengineHelpguideMd,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"appengine": &bintree{nil, map[string]*bintree{
		"helpguide.md": &bintree{appengineHelpguideMd, map[string]*bintree{}},
		"template": &bintree{nil, map[string]*bintree{
			"content.html": &bintree{appengineTemplateContentHtml, map[string]*bintree{}},
			"root.html":    &bintree{appengineTemplateRootHtml, map[string]*bintree{}},
			"u.html":       &bintree{appengineTemplateUHtml, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
