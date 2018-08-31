// Code generated by go-bindata.
// sources:
// appengine/template/index.html
// appengine/template/u.html
// DO NOT EDIT!

package main

import (
	"bytes"
	"compress/gzip"
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
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
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

var _appengineTemplateIndexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\x50\x3d\x4f\xc4\x30\x0c\xdd\xef\x57\x18\xef\xd7\xac\x48\x24\x1d\xf8\x5a\x61\x28\x03\x63\x2e\xb6\x54\xa3\xb6\x89\x12\xab\x07\xaa\xfa\xdf\x51\x49\x2b\xc0\x8b\xfd\xfc\xfc\x9e\x2d\xdb\x9b\xc7\x97\x87\xee\xfd\xf5\x09\x7a\x1d\x87\xf6\x64\x8f\xc4\x9e\xda\x13\x00\x80\x1d\x59\x3d\x84\xde\xe7\xc2\xea\xf0\xad\x7b\x3e\xdf\x22\x98\xbf\xe4\xe4\x47\x76\x38\x0b\x5f\x53\xcc\x8a\x10\xe2\xa4\x3c\xa9\xc3\xab\x90\xf6\x8e\x78\x96\xc0\xe7\x1f\x80\xbb\x4e\x45\x07\x6e\x97\xa5\xe9\xb6\x62\x5d\xad\xa9\x9d\xca\x96\x90\x25\x69\x05\x5b\x0c\xac\x70\x89\x51\x8b\x66\x9f\xc0\xc1\xb2\x34\xf7\x07\x5c\xd7\xbb\x2a\x32\x87\xca\x9a\x7a\xbd\xbd\x44\xfa\xda\x1d\x49\x66\x10\x72\xe8\x53\xc2\xd6\x1a\x92\xf9\xdf\x2a\x28\x39\x38\x34\x24\x45\x8d\x4c\xc4\x9f\xcd\x47\xd9\xe6\x7e\x2d\xab\x97\x35\xf5\x3f\xdf\x01\x00\x00\xff\xff\x37\x0d\x40\xdb\x37\x01\x00\x00")

func appengineTemplateIndexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_appengineTemplateIndexHtml,
		"appengine/template/index.html",
	)
}

func appengineTemplateIndexHtml() (*asset, error) {
	bytes, err := appengineTemplateIndexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "appengine/template/index.html", size: 311, mode: os.FileMode(420), modTime: time.Unix(1535676087, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _appengineTemplateUHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x54\x4d\x8f\xdb\x36\x10\xbd\xfb\x57\x70\x75\x08\x25\x84\xa1\x92\x53\x0b\xcb\xdc\xa2\xf9\x58\xa0\x97\x34\xc0\x6e\x0a\xb4\x86\xb1\xa0\xa9\x91\xc4\x2c\x45\x6a\xc9\x91\x5c\x41\xab\xff\x5e\x48\xb2\xd7\xeb\xa4\xad\x0e\xb6\x38\xf3\xde\xd3\x9b\x21\x39\x9b\xab\x8f\xbf\x7f\xb8\xfb\xf3\xcb\x27\x52\x61\x6d\xae\x57\x9b\xd3\x1f\xc8\xfc\x7a\x45\x08\x21\x9b\x1a\x50\x12\x55\x49\x1f\x00\x45\xf4\xf5\xee\xe6\xcd\xcf\x11\x49\x5f\x26\xad\xac\x41\x44\x9d\x86\x43\xe3\x3c\x46\x44\x39\x8b\x60\x51\x44\x07\x9d\x63\x25\x72\xe8\xb4\x82\x37\xf3\x22\x3a\xf2\x50\xa3\x81\xeb\x61\xe0\x77\xd3\xcb\x38\x6e\xd2\x25\xb2\x64\x8d\xb6\x0f\xc4\x83\x11\x51\xc0\xde\x40\xa8\x00\x30\x22\xd8\x37\x20\x22\x84\xbf\x31\x55\x21\x44\xa4\xf2\x50\x88\x28\x0d\x28\x51\xab\xb4\x70\x16\xe7\x1f\x3e\x25\x8f\x42\x41\x79\xdd\xe0\xb2\x98\x1e\x03\x48\x5a\xe9\xcb\x40\x04\x19\x06\x3e\x8e\xd9\x82\x4b\x5f\x02\x7f\x60\xe9\x82\xc4\x57\x33\x8d\x37\x5e\x77\x12\x21\x21\xc3\xea\xa0\x6d\xee\x0e\x5b\x7a\x5f\x84\x7b\x0f\x32\xef\xe9\x8e\x08\x52\xb4\x56\xa1\x76\x36\x9e\x20\x27\xf2\x09\x7a\x73\x4b\x77\x5b\x8a\x55\x5b\xef\x0b\x80\x7c\x2f\xd5\x03\xdd\x9d\x80\xd3\xf3\xbf\x40\x5e\xe9\x1c\xe2\x64\x06\x8f\xab\x71\x75\xe1\x20\x87\x7d\x5b\xce\x0e\xd0\xb7\x90\x5d\xe4\x2a\x17\x70\x4e\xd1\x80\xb2\xd4\xb6\xe4\x45\x6b\x4c\x40\xe7\x7b\xae\x5c\x4d\x2f\xd1\xce\x2f\x3a\xf4\xdd\x5f\x1f\x7f\xfa\x2e\x37\xed\x74\x68\xa4\x82\x05\x71\x73\x4b\xb3\x55\xfc\x5c\x72\xcd\x2c\x03\x86\xcc\x30\xc7\x4a\xd6\x27\xe7\x06\x00\xd1\x96\xd4\x09\x19\x74\x11\xd7\x5c\x39\x1b\x9c\x01\xf2\xea\x15\x79\x5e\x70\xe3\xca\x84\x0c\x97\x81\x98\xde\xb4\xc6\xdc\x4e\x46\xc9\xf3\xa7\xa7\xe3\x55\x18\xad\x90\x93\x2f\x06\x64\x00\x12\x00\x4f\x9d\x8b\x2e\x5c\x46\x3b\x4e\x93\x6c\x24\x1e\xb0\xf5\x36\x1b\x67\x3f\xa5\xa8\xb7\xb0\x13\xcf\xae\x25\xdb\x27\x43\xc9\x1f\x7f\x29\xf9\x23\x6f\xda\x50\xc5\x5b\xc9\xf6\xbb\x64\x5d\xf2\x7b\xd9\xe8\x39\x9f\x8d\x59\xc9\x1f\xc5\x76\xb7\x9c\x17\x27\x2c\x57\x1e\x24\xc2\x27\x03\x35\x58\x8c\x31\xc9\x1c\x97\xa1\xb7\x4a\xbc\xcb\x1c\x0f\x5e\x09\x5a\x21\x36\x61\x9d\xa6\xf4\xf5\x69\x13\x5e\xd3\x34\xa4\x45\xe0\xdf\x02\x5d\x84\x7a\x61\x79\x09\x78\x54\x09\xef\xfb\x3b\x59\x7e\x96\x35\xc4\x98\x6c\xdf\xee\xb2\x9e\x37\xd2\x83\xc5\xcf\x2e\x07\xae\x6d\x00\x8f\xef\xa1\x70\x1e\x62\xc7\xfa\x64\x91\x28\xb9\xce\xc1\xa2\x2e\xfa\x73\x49\x9a\x75\xc9\x50\xc6\x86\x0d\xad\xce\xd7\x7a\x4c\x32\x5d\xc4\x5d\x32\x45\xba\x64\x2a\x25\x00\x7e\x0d\xe0\xff\x90\x3e\x9c\x59\x47\xce\x82\x80\x6e\xba\xc1\xdf\x2b\xd2\x39\x4c\xd9\x60\xd7\x9a\x35\xeb\x6e\x4c\xc6\x93\x8b\x50\xb5\x98\xbb\x83\x3d\x73\x26\x42\xe4\x41\x45\xec\xea\xdd\xac\xe9\x21\xa0\xf4\xf8\xaf\x88\xb7\x67\xa5\xe9\x00\x5c\x7c\x5c\xce\xb8\x63\x38\x62\x57\xd2\x97\xed\xdc\x2f\x6e\xc0\x96\x58\x3d\x3d\xc9\x33\xfb\xd4\x8d\x5f\x95\x72\xed\x0f\x25\x38\x41\xe5\x92\xa0\x59\x27\xba\xa7\xa7\x61\xcc\x3a\x2e\x95\xc2\xdf\x72\xa1\xb3\x32\x76\x73\xfd\x27\x27\x06\xa4\x9f\x1a\xf5\xc1\xb9\x07\x0d\x2f\x8d\x8f\xd9\x6a\x4c\x8e\x17\x9b\xe5\x4e\xcd\x86\xd8\x7f\xdc\x15\x46\x97\x91\x42\x19\x6d\x03\x78\x7a\xdc\xba\xe5\x22\x5f\x0e\xa0\x4d\xba\x8c\xde\xcd\xde\xe5\xfd\x71\x1e\xe5\xba\x23\x3a\x17\x91\x6c\x9a\xe8\x7a\x93\xe6\xba\xbb\x18\x54\x64\x3a\x6d\x51\x9a\xeb\x80\x69\xcb\xbf\x85\x09\x73\x96\x5b\x74\x36\xe9\x32\xd8\xff\x09\x00\x00\xff\xff\x4d\x7b\x66\x6a\xf0\x05\x00\x00")

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

	info := bindataFileInfo{name: "appengine/template/u.html", size: 1520, mode: os.FileMode(420), modTime: time.Unix(1535724061, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
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
	"appengine/template/index.html": appengineTemplateIndexHtml,
	"appengine/template/u.html": appengineTemplateUHtml,
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
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
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
			"index.html": &bintree{appengineTemplateIndexHtml, map[string]*bintree{}},
			"u.html": &bintree{appengineTemplateUHtml, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
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
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
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
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

