// Code generated by go-bindata. DO NOT EDIT.
// sources:
// appengine/template/dashboard.html (739B)
// appengine/template/home.html (777B)
// appengine/template/u.html (674B)

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

var _appengineTemplateDashboardHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x92\xbf\x6e\x1b\x31\x0c\xc6\xf7\x3c\x05\xa3\xa1\x4b\x70\xc7\xd4\x53\x91\x48\x6e\x91\xfe\x01\x3a\x04\xcd\xe0\x0e\x1d\x65\x89\xce\x31\xd5\x49\xc2\x91\xb5\x7b\x30\xfc\xee\x85\x7d\x36\x1a\x77\x32\xaa\x41\x12\x49\xfd\xbe\x4f\x04\x68\xaf\x3f\x7d\xfb\xb8\xf8\xf1\xf4\x19\x3a\xed\xd3\xfc\xca\x9e\x0e\xf2\x71\x7e\x05\x00\x60\x7b\x52\x0f\xa1\xf3\x83\x90\x3a\xf3\x7d\xf1\xa5\x79\x67\x00\x5f\x17\xb3\xef\xc9\x99\x35\xd3\xa6\x96\x41\x0d\x84\x92\x95\xb2\x3a\xb3\xe1\xa8\x9d\x8b\xb4\xe6\x40\xcd\x21\x30\x47\x4e\x59\x13\xcd\xb7\xdb\x76\xb1\xbf\xec\x76\x16\xa7\xcc\x54\xbd\x6e\x1a\x58\x95\xac\x02\x4d\x73\x4c\x25\xce\x3f\x61\xa0\xe4\x8c\xe8\x98\x48\x3a\x22\x35\xd0\x0d\xb4\x72\x06\xb1\xa3\x94\x4a\xdb\x8f\x07\xa8\xcd\xa4\x18\xca\xaf\xac\x38\x5b\xdd\xde\x2e\x67\xe6\x02\x8d\x4e\xb5\xca\x1d\xe2\xa4\xf0\x5c\xca\x73\x22\x5f\x59\xda\x50\x7a\x0c\x22\xb3\xf7\x2b\xdf\x73\x1a\xdd\xd7\x87\xc7\x9b\xa7\x44\xbf\x6f\x1e\x4b\x2e\x77\xac\x3e\x7d\xb8\xbd\x7f\xfb\x26\xb2\xd4\xe4\x47\x27\x1b\x5f\x2f\xf1\x43\x51\xaf\x1c\x0e\x7e\x87\xad\x0d\x22\x27\x50\xc2\xc0\x55\xa7\x60\xbf\x42\xc9\xa2\xb0\x2c\x45\x45\x07\x5f\xc1\xc1\x76\xdb\x3e\x9c\xc2\xdd\xee\x7e\xc2\xf0\x35\xf7\xff\xee\x97\x80\x41\x04\x73\x51\xa6\x23\x68\x71\x9a\x18\xbb\x2c\x71\x3c\xea\x44\x5e\x03\x47\x67\x7c\xad\x66\x6e\x31\xf2\xfa\xac\x3d\x90\x21\xfc\x15\x7c\x39\xe9\xf5\x9c\xdb\x17\xd9\x03\x67\xdd\x9c\x31\x91\x45\x31\x7a\xe9\x96\xc5\x0f\xf1\x9f\xe7\x16\xa7\x3f\x58\x9c\x66\xf9\x4f\x00\x00\x00\xff\xff\x1c\x4f\xa4\xd9\xe3\x02\x00\x00")

func appengineTemplateDashboardHtmlBytes() ([]byte, error) {
	return bindataRead(
		_appengineTemplateDashboardHtml,
		"appengine/template/dashboard.html",
	)
}

func appengineTemplateDashboardHtml() (*asset, error) {
	bytes, err := appengineTemplateDashboardHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "appengine/template/dashboard.html", size: 739, mode: os.FileMode(0644), modTime: time.Unix(1669556276, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x3d, 0x4d, 0x1a, 0x30, 0xd3, 0x32, 0x3c, 0xb8, 0x8, 0x14, 0xa0, 0xba, 0xfb, 0xaa, 0x55, 0xf, 0x65, 0xd2, 0xfb, 0xc4, 0x30, 0x1, 0xb9, 0xe, 0x8a, 0xf5, 0xbc, 0xba, 0x24, 0xa5, 0x3a, 0xd7}}
	return a, nil
}

var _appengineTemplateHomeHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x92\xc1\x8e\xd3\x30\x10\x86\xef\xfb\x14\xb3\x3e\x70\x59\x12\x97\x9e\xd0\x6e\x52\xe8\x52\x90\x38\x54\xbb\x12\xe5\xc0\xd1\xb1\x27\x8d\xc1\xf1\x58\x99\x69\x43\x15\xe5\xdd\x51\x1b\xa2\x5d\xd4\x0b\x39\xc4\x1e\xcf\xfc\xdf\x3f\xb2\xa7\xb8\xdd\x3c\x7d\xda\xfd\x78\xfe\x0c\x8d\xb4\x61\x75\x53\xcc\x0b\x1a\xb7\xba\x01\x00\x28\x5a\x14\x03\xb6\x31\x1d\xa3\x94\xea\xfb\xee\x4b\xf6\x5e\x81\x7e\x9d\x8c\xa6\xc5\x52\x1d\x3d\xf6\x89\x3a\x51\x60\x29\x0a\x46\x29\x55\xef\x9d\x34\xa5\xc3\xa3\xb7\x98\x5d\x02\x75\xad\x73\xc8\xb6\xf3\x49\x3c\xc5\x57\xd2\x35\x54\x68\x0e\xe2\xeb\x43\x00\xb6\x1d\x55\x55\xf0\x71\x0f\x8c\xdd\x19\x06\x35\x75\xb0\x4e\x29\x20\x6c\x0f\xec\x6d\x0e\x1b\xea\x63\x20\xe3\x40\x1a\x84\xd6\xd8\xa7\x6f\x60\x52\x7a\x0b\x26\x3a\xa8\x3a\xea\x19\x67\x0c\x32\x50\xbc\x94\xf5\x58\xb1\x17\xcc\xe7\xa6\x82\x8f\xbf\xa0\xc3\x50\x2a\x6b\x22\x45\x6f\x4d\x50\xd0\x74\x58\x97\xaa\x11\x49\x7c\xaf\xf5\x30\xe4\xeb\x94\x36\xd4\x1a\x1f\xc7\x51\xbf\x5c\x84\x78\x09\xb8\x1a\x86\x7c\x77\xde\x8c\x63\xa1\xa7\x93\x29\x7b\x9b\x65\x50\x53\x14\x86\x2c\xbb\x32\x63\x39\x05\xe4\x06\x51\x66\x37\xad\x1b\x0c\x81\xf2\xf6\x74\x11\xe5\x11\x45\x5b\x3a\x44\xd1\xcb\x7a\xb1\xa8\x96\xd7\x0d\x5f\x33\xe6\x8e\x27\xc2\x9e\x68\x1f\xd0\x24\xcf\xb9\xa5\x56\x5b\xe6\xe5\x87\xda\xb4\x3e\x9c\xca\xaf\x8f\xdb\xbb\xe7\x80\xbf\xef\xb6\x14\xe9\xde\x8b\x09\x1f\x17\x0f\xef\xde\x38\xcf\x29\x98\x53\xc9\xbd\x49\xff\xe3\xa7\x59\x8c\x78\x7b\xf1\xbb\xfc\x72\xcb\x3c\x0b\xa7\x17\x9e\x82\xf3\x67\x29\xb2\x40\x45\x24\x2c\x9d\x49\x50\xc2\x30\xe4\x8f\x73\x38\x8e\x0f\x93\x4c\xcf\xba\x42\x4f\x03\x59\x54\xe4\x4e\x7f\x99\xce\x1f\xc1\xbb\x52\x99\x94\xd4\xaa\xd0\xce\x1f\xff\x31\x03\xee\x6c\xa9\xb4\xf3\x2c\xba\xa1\x16\xf3\x9f\x7c\x2e\x7b\x21\x4e\xa8\x42\x4f\x13\xff\x27\x00\x00\xff\xff\xfa\x51\xe1\xf6\x09\x03\x00\x00")

func appengineTemplateHomeHtmlBytes() ([]byte, error) {
	return bindataRead(
		_appengineTemplateHomeHtml,
		"appengine/template/home.html",
	)
}

func appengineTemplateHomeHtml() (*asset, error) {
	bytes, err := appengineTemplateHomeHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "appengine/template/home.html", size: 777, mode: os.FileMode(0644), modTime: time.Unix(1669556276, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xe4, 0x68, 0xe5, 0x12, 0xe, 0x13, 0xca, 0xd5, 0x70, 0xb7, 0x11, 0x7d, 0x4c, 0xe8, 0x92, 0x9d, 0xcd, 0x68, 0xaa, 0xbb, 0x2f, 0xd0, 0xfe, 0xd8, 0x1a, 0xba, 0xf4, 0xb7, 0xd7, 0xc9, 0xe0, 0xf6}}
	return a, nil
}

var _appengineTemplateUHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x92\x4d\x6f\x1b\x21\x10\x86\xef\xf9\x15\x13\x0e\xbd\x44\xbb\xe3\xfa\x54\x25\xb0\xad\xfa\x25\xf5\x10\x35\x07\xf7\xd0\x23\x81\xb1\x97\x94\x05\xc4\x8c\xed\xae\x2c\xff\xf7\xca\x5e\x47\xb1\x7b\x88\xc2\x01\x98\x8f\xf7\x79\x47\x02\x7d\xfd\xf5\xe7\x97\xc5\xef\x87\x6f\xd0\xcb\x10\xbb\x2b\xfd\x7c\x90\xf5\xdd\x15\x00\x80\x1e\x48\x2c\xb8\xde\x56\x26\x31\xea\xd7\xe2\x7b\xf3\x41\x01\x9e\x17\x93\x1d\xc8\xa8\x4d\xa0\x6d\xc9\x55\x14\xb8\x9c\x84\x92\x18\xb5\x0d\x5e\x7a\xe3\x69\x13\x1c\x35\xc7\x40\x9d\x74\x12\x24\x52\xb7\xdb\xb5\x8b\xc3\x65\xbf\xd7\x38\x65\xa6\xea\x75\xd3\xc0\x32\x27\x61\x68\x9a\xb3\x94\x8e\x21\xfd\x81\x4a\xd1\x28\x96\x31\x12\xf7\x44\xa2\xa0\xaf\xb4\x34\x0a\xb1\xa7\x18\x73\x3b\x8c\x47\x65\x9b\x48\xd0\xe5\x75\x12\x9c\x2f\x67\xb3\xc7\xb9\xea\x5e\x60\xaf\x71\x7a\x91\xc2\xb7\x88\x13\x65\x95\xf3\x2a\x92\x2d\x81\x5b\x97\x07\x74\xcc\xf3\x8f\x4b\x3b\x84\x38\x9a\x1f\x9f\xef\x6f\x1e\x22\xfd\xbd\xb9\xcf\x29\xdf\x06\xb1\xf1\xd3\xec\xee\xfd\x3b\x1f\xb8\x44\x3b\x1a\xde\xda\xa2\xde\xe0\x87\x2c\x56\x82\x3b\xfa\x1d\xb7\xd6\x31\x3f\x0b\xd9\xd5\x50\x64\x0a\x0e\xcb\xe5\xc4\x02\x6b\x5b\x57\x0c\x06\x76\xbb\x76\xbf\xbf\x9b\x3a\xf1\xbc\xf5\x4d\x86\x8e\x19\x53\xa9\x79\x55\x89\xf9\x64\xaa\x71\x7a\x78\xfd\x98\xfd\x78\x62\xf9\xb0\x81\xe0\x8d\xb2\xa5\xa8\x4e\xa3\x0f\x9b\x8b\xe1\x80\xab\x7b\x81\x3e\x9d\x33\x9f\xf8\x20\xb8\x18\xec\x42\xe3\x03\x0b\xae\xff\x6b\xd3\x38\x79\x6b\x9c\xbe\xe2\xbf\x00\x00\x00\xff\xff\x2c\xfa\x52\x78\xa2\x02\x00\x00")

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

	info := bindataFileInfo{name: "appengine/template/u.html", size: 674, mode: os.FileMode(0644), modTime: time.Unix(1669556276, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x6e, 0xab, 0x12, 0xbf, 0x7e, 0x2b, 0xf6, 0xb6, 0x1e, 0x16, 0x37, 0x54, 0x99, 0x28, 0xe5, 0xd, 0x8f, 0xa4, 0x2c, 0xbd, 0x79, 0xc5, 0x94, 0xf7, 0x84, 0x20, 0x46, 0x32, 0x49, 0xf0, 0x70, 0x83}}
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
	"appengine/template/dashboard.html": appengineTemplateDashboardHtml,
	"appengine/template/home.html":      appengineTemplateHomeHtml,
	"appengine/template/u.html":         appengineTemplateUHtml,
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
		"template": &bintree{nil, map[string]*bintree{
			"dashboard.html": &bintree{appengineTemplateDashboardHtml, map[string]*bintree{}},
			"home.html":      &bintree{appengineTemplateHomeHtml, map[string]*bintree{}},
			"u.html":         &bintree{appengineTemplateUHtml, map[string]*bintree{}},
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
