// Code generated by go-bindata. DO NOT EDIT.
// sources:
// appengine/template/root.html (310B)
// appengine/template/u.html (464B)

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

var _appengineTemplateRootHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\x50\x3d\x4f\xc4\x30\x0c\xdd\xef\x57\x18\xef\xd7\xac\x48\x24\x1d\xf8\x5a\x61\x28\x03\x63\x2e\xb6\x54\xa3\xb6\x89\x12\xab\x27\x54\xf5\xbf\xa3\x92\x56\x70\x5e\xec\xe7\xe7\xf7\x6c\xd9\xde\x3d\xbf\x3d\x75\x9f\xef\x2f\xd0\xeb\x38\xb4\x27\x7b\x24\xf6\xd4\x9e\x00\x00\xec\xc8\xea\x21\xf4\x3e\x17\x56\x87\x1f\xdd\xeb\xf9\x1e\xc1\xfc\x27\x27\x3f\xb2\xc3\x59\xf8\x9a\x62\x56\x84\x10\x27\xe5\x49\x1d\x5e\x85\xb4\x77\xc4\xb3\x04\x3e\xff\x02\xdc\x75\x2a\x3a\x70\xbb\x2c\x4d\xb7\x15\xeb\x6a\x4d\xed\x54\xb6\x84\x2c\x49\x2b\xd8\x62\x60\x85\x4b\x8c\x5a\x34\xfb\x04\x0e\x96\xa5\x79\x3c\xe0\xba\x3e\x54\x91\x39\x54\xd6\xd4\xeb\xed\x25\xd2\xf7\xee\x48\x32\x83\x90\x43\x9f\x12\xb6\xd6\x90\xcc\x37\xab\xa0\xe4\xe0\xd0\x90\x14\x35\x39\x46\x6d\xbe\xca\x36\xf6\xe7\x58\xad\xac\xa9\xef\xf9\x09\x00\x00\xff\xff\xa5\x14\xbc\xb3\x36\x01\x00\x00")

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

	info := bindataFileInfo{name: "appengine/template/root.html", size: 310, mode: os.FileMode(0664), modTime: time.Unix(1588530134, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x8a, 0x2f, 0x3e, 0xe7, 0x57, 0x60, 0x7c, 0x4, 0x65, 0x22, 0x90, 0xa6, 0xde, 0x6, 0x13, 0xc5, 0xcb, 0x51, 0xd4, 0xa7, 0xfc, 0x72, 0xa1, 0xf7, 0xba, 0x23, 0x5c, 0xcd, 0xc8, 0xde, 0x86, 0x5a}}
	return a, nil
}

var _appengineTemplateUHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x91\x3f\x6f\xf3\x20\x10\xc6\xf7\x7c\x8a\x7b\xd9\x63\xd6\x57\x2a\x78\xe9\x9f\xb5\x1d\xd2\xa1\x23\x85\x4b\x7c\x29\xc6\x16\x77\x71\x14\x59\xfe\xee\x95\x83\xa3\x5a\x51\x87\x32\x00\xc7\x3d\xcf\x8f\x47\x60\xfe\x3d\xbd\x3e\xee\x3e\xde\x9e\xa1\x91\x36\xd6\x1b\x73\x5b\xd0\x85\x7a\x03\x00\x60\x5a\x14\x07\xbe\x71\x99\x51\xac\x7a\xdf\xbd\x6c\xff\x2b\xd0\xeb\x66\x72\x2d\x5a\x35\x10\x9e\xfb\x2e\x8b\x02\xdf\x25\xc1\x24\x56\x9d\x29\x48\x63\x03\x0e\xe4\x71\x7b\x2d\xd4\xe2\x13\x92\x88\xf5\x38\x56\xbb\x79\x33\x4d\x46\x97\x93\xd2\x8d\x94\xbe\x20\x63\xb4\x8a\xe5\x12\x91\x1b\x44\x51\xd0\x64\xdc\x5b\xa5\x59\x9c\x90\xd7\xfb\x2e\xc9\x75\xaa\x3c\xf3\x0d\xcb\x3e\x53\x2f\xa5\x98\x47\x44\x81\x93\xcb\x07\x06\x0b\xe3\x58\x4d\xd3\x43\xd1\xe9\xb5\x70\x71\x01\x67\xff\x83\x3f\xb2\x4e\x7d\xee\x0e\x19\x99\xab\x23\xab\xfa\xce\xf3\x97\x88\x9e\xd7\x90\x12\xd3\xe8\xf2\xb4\xe6\xb3\x0b\x97\x85\x15\x68\x00\x0a\x56\xb9\xbe\x9f\xef\x09\x34\xfc\x16\x2c\x10\x8b\x3e\xdd\x65\x31\xba\x70\x8c\x2e\x1f\xf7\x1d\x00\x00\xff\xff\x0c\xb4\x8d\x9d\xd0\x01\x00\x00")

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

	info := bindataFileInfo{name: "appengine/template/u.html", size: 464, mode: os.FileMode(0664), modTime: time.Unix(1588530142, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x83, 0x8c, 0xaa, 0x1, 0xfc, 0x6c, 0xe6, 0x5b, 0xf5, 0x1, 0x81, 0x27, 0xd0, 0x13, 0x98, 0x6, 0xe7, 0xff, 0xf5, 0xad, 0x61, 0xe5, 0x73, 0xf0, 0x76, 0xd3, 0x46, 0x52, 0x58, 0xa0, 0xa6, 0x97}}
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
	"appengine/template/root.html": appengineTemplateRootHtml,
	"appengine/template/u.html":    appengineTemplateUHtml,
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
			"root.html": &bintree{appengineTemplateRootHtml, map[string]*bintree{}},
			"u.html":    &bintree{appengineTemplateUHtml, map[string]*bintree{}},
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
