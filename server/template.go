// Code generated by go-bindata.
// sources:
// web/template/index.html
// DO NOT EDIT!

package server

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

var _webTemplateIndexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x94\xdf\x6f\xe4\x26\x10\xc7\x9f\x6f\xff\x8a\x29\x7d\x66\xcc\x6f\x9b\x0a\x5b\x4a\xb6\x77\xbd\x97\xaa\x7d\xa8\x2a\xf5\x69\xe5\xd8\x9c\x8d\xca\xda\xd1\xe2\x78\x2f\xfd\xeb\x2b\xec\x5b\x25\xa9\x72\x95\x76\x61\x98\x61\xe0\xf3\x1d\x03\xae\xf8\xe1\xe7\xdf\x8e\x7f\xfc\xf5\xfb\x47\x18\x97\x73\x6c\x0e\xee\xd6\xf9\xb6\x6f\x0e\x1f\xdc\x12\x96\xe8\x9b\xbb\xf8\xf0\x74\x86\xfb\x70\x59\xc6\xbe\x7d\x4e\xae\xd8\xdd\x87\x0f\x2e\x86\xe9\x6f\xb8\xf8\x58\x93\xb4\x3c\x47\x9f\x46\xef\x17\x02\xe3\xc5\x7f\xa9\x49\xd1\xa5\x54\x9c\xdb\x30\x61\x97\x12\x69\x0e\xae\xd8\x57\x75\x0f\x73\xff\x9c\x93\xfb\xb0\x42\x17\xdb\x94\x6a\x12\xe7\x21\x4c\xa7\x53\x9a\xbb\xd0\x46\xd2\x7c\x37\x44\x1f\x9e\x96\x65\x9e\xd2\xdb\x29\x6f\x63\xa7\xd3\x6e\xd0\x6e\x9e\x96\x36\x4c\xfe\x42\x1a\xb7\xbb\xfe\x3f\x01\xf6\x8e\x34\x2e\xad\xc3\x7f\xa6\x86\x6e\x9e\x12\xbc\x1e\x9c\x4e\xc3\x3c\x0f\xd1\xbf\x71\x52\xea\xa7\xf6\x21\xfa\x9e\xc0\x35\xf4\xcb\x58\x13\xc1\x08\x8c\x3e\x0c\xe3\xb2\xdb\x6b\xf0\xd7\xfb\xf9\x6b\x4d\x18\x30\x10\xf9\x47\xe0\xeb\x39\x4e\xa9\x26\xe3\xb2\x3c\xfe\x54\x14\xd7\xeb\x15\xaf\x12\xe7\xcb\x50\x08\xc6\x58\x91\xd6\x81\x34\x6e\x80\x2f\x21\xc6\x9a\x4c\xf3\xe4\xc9\x66\xd3\xcb\x53\xf4\x35\xf1\xab\x9f\xe6\xbe\x27\x8d\x7b\x6c\x97\x11\xfa\x9a\xfc\xca\x2d\x1a\xe0\x0c\x85\x28\x3b\x46\xb1\x64\x96\x22\x33\x8a\x72\x94\x96\x22\xaf\x04\x15\xc8\x94\xfe\xcc\xd9\x2a\xb1\x32\xd5\xa8\x51\x56\xa2\x55\x68\x20\xff\x33\x19\xa7\x1c\xad\x35\x20\x91\xf1\x6a\x15\xa8\xf9\x28\x51\x48\xd1\x71\xac\x6c\x0e\x96\x4a\x80\x40\x5b\x09\xaa\x50\x32\xfd\xcd\x2e\x51\xea\x7f\xc8\x37\xd4\x1f\x95\xa8\xf4\x27\x45\x1a\x57\x64\xb4\xd7\x80\x59\x77\x27\xb0\x04\x06\x0a\xad\x51\x14\x2b\xab\xc1\xa0\xe1\x15\x15\xa8\x84\x8c\x74\xdb\x8e\x0a\xd4\xcc\x76\x5b\x18\xcd\x86\x8d\x56\x6b\x2a\x51\x56\x66\xb3\x04\x1a\xa6\x81\x51\x85\xd5\x86\x65\xa8\x46\x6d\x35\x55\xc8\x85\xfc\xcc\xb3\xee\x8c\x6f\xef\xec\xa6\x67\x6f\xd9\xae\x31\x53\xbc\xd0\x4a\x75\x57\x69\xf9\x0e\xad\x42\xc5\x34\x70\x8e\x19\x44\x50\x34\x14\x25\xcf\xd5\x14\xea\x66\x59\x60\x14\x8d\x41\xbe\x8d\xe4\xcd\xfb\xa7\x46\xcd\x77\x88\x77\xf6\xcf\x04\x1d\x03\x8e\x86\xab\xac\x07\x24\x72\x05\xdb\x6c\x50\xa8\x6c\x94\x28\x55\xae\x80\x7d\x81\xfc\x74\x7f\x7f\x64\xfa\xfd\x92\x4a\xb4\x65\xd9\x71\x54\xa6\xca\x67\x0b\xcb\xca\xa0\x66\x1a\x24\x56\x42\x02\x47\x65\x4d\x14\xf9\x7b\xd3\xad\x3d\xf2\x5c\xcb\x0c\x04\x5c\xa0\xb1\x7a\x2f\x09\x03\x83\xcc\xee\x0b\x70\x10\x28\x6e\x48\x59\xcb\x86\x04\x19\xe9\xa8\x91\x5b\x01\x1a\x4b\x69\xa0\x44\xb9\x9d\x15\x5b\x96\x70\x03\x79\x61\xfe\x78\xa7\xa4\x7c\xc5\x5c\x0c\x8d\xcb\xa7\xba\x71\xe9\xb1\xfd\xee\x9d\x4c\xfe\xb2\x86\xce\xd3\xa9\x3d\x7b\xd2\x1c\xe7\x69\x09\xd3\x93\x87\x6b\x58\x46\xf8\x65\xbb\x79\xae\xc8\xf9\x8d\x2b\xf6\x9c\xc6\x15\x7d\x58\xdf\xb4\x07\x57\xec\x8f\x8d\x2b\xf6\x87\xed\xdf\x00\x00\x00\xff\xff\x9d\x5d\x41\x19\xf1\x04\x00\x00")

func webTemplateIndexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_webTemplateIndexHtml,
		"web/template/index.html",
	)
}

func webTemplateIndexHtml() (*asset, error) {
	bytes, err := webTemplateIndexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "web/template/index.html", size: 1265, mode: os.FileMode(420), modTime: time.Unix(1526005090, 0)}
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
	"web/template/index.html": webTemplateIndexHtml,
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
	"web": &bintree{nil, map[string]*bintree{
		"template": &bintree{nil, map[string]*bintree{
			"index.html": &bintree{webTemplateIndexHtml, map[string]*bintree{}},
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

