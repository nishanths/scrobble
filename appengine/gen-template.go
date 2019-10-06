// Code generated by go-bindata.
// sources:
// appengine/template/fs-snippet.tmpl
// appengine/template/root.tmpl
// appengine/template/u.tmpl
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

var _appengineTemplateFsSnippetTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x93\x51\x6f\x9b\x3e\x14\xc5\xdf\xf9\x14\x6e\x1e\x6a\xdc\x5a\x4e\xfb\xf4\x97\x62\xa1\xbf\xd6\x6d\x91\xf6\x52\x4d\xca\xb6\x87\x45\xa8\x72\xcc\xc5\xf1\x0a\x36\xb5\x0d\x15\x22\xfe\xee\x53\x48\x48\x96\x75\x1b\x6f\xdc\x7b\xee\xf1\xef\x1e\xf0\x30\x14\x50\x6a\x03\x68\xb6\x5c\xad\x8c\x6e\x1a\x08\xb3\x18\x93\xf9\x4d\xf2\xaa\x4d\x61\x5f\xd7\xf8\xa9\xf4\x4f\x0e\x44\xd1\xe3\x1c\x65\xa8\x6c\x8d\x0c\xda\x9a\x94\xa0\x21\x41\x08\x21\x5d\xa2\x74\x92\x2e\x57\x38\x5f\xe3\xb0\x6d\xeb\x4d\x09\x50\x6c\x84\x7c\xc6\xf9\x24\xdc\x3f\xff\x14\xb2\xad\x2e\x20\x25\xa3\x38\x26\x31\xb9\x20\x28\x60\xd3\xaa\x03\x81\xa8\x3c\xf0\x8b\xe6\xd6\xfa\x30\xf6\xb0\x0f\x42\x69\xa3\x58\xd9\x56\x95\x0f\xd6\xf5\x4c\xda\x1a\x5f\xaa\xad\x3b\x18\xe1\xfb\xef\x1f\xfe\xfb\xad\x67\x44\x0d\xbe\x11\x12\x0e\x8a\xe5\x0a\xf3\x24\x3d\xed\x5c\x53\x43\x81\x06\x5a\x51\x4b\x15\xed\xc9\x39\x01\x40\xda\xa0\x9a\xa0\x41\x97\x69\xcd\xa4\x35\xde\x56\x80\xae\xaf\xd1\xe9\x85\x55\x56\x11\x34\x5c\x16\x52\xbc\x6c\xab\x6a\xb5\x07\x45\xa7\xa3\x91\xb4\xa6\xac\xb4\x0c\x0c\x7d\xae\x40\x78\x40\x1e\xc2\x14\xdd\xec\x82\x72\x96\x33\x4c\x78\x44\x0e\x42\xeb\x0c\x8f\x23\x8f\xca\xea\x35\xe4\xd9\x89\x5a\xd0\x0d\x19\x14\x7b\xf9\x5f\xb1\x17\xd6\xb4\x7e\x9b\xae\x05\xdd\xe4\x64\xa1\xd8\x93\x68\xf4\xd8\xe7\x91\x2b\xf6\x92\xad\x73\x3e\x5a\xd8\xcc\x30\xe9\x40\x04\xf8\x58\x41\x0d\x26\xa4\x81\x70\xcb\x84\xef\x8d\xcc\xee\xb9\x65\xde\xc9\x0c\x6f\x43\x68\xfc\x62\x3e\xc7\xb7\xd3\x47\xb8\xc5\x73\x3f\x2f\x3d\xfb\xe1\xf1\xc1\xa8\xcf\x0c\x53\x10\x8e\x2e\xfe\xa1\xff\x22\xd4\xa3\xa8\x21\x0d\x64\x7d\x97\xf3\x9e\x35\xc2\x81\x09\x8f\xb6\x00\xa6\x8d\x07\x17\x1e\xa0\xb4\x0e\x52\x4b\x7b\x72\xb0\x50\x4c\x17\x60\x82\x2e\xfb\xf3\x4a\x9a\x76\x64\x50\x69\x45\x87\x56\x17\x0b\x1d\x09\xd7\x65\xda\x91\x7d\xa5\x23\xfb\x55\x3c\x84\xaf\x1e\xdc\x37\xe1\xfc\x79\xea\x38\x73\x50\x40\x07\x26\xbc\x71\xc4\x63\x19\xd3\xc1\x2c\x34\x6d\x16\x5d\x24\x71\xa2\xf0\xdb\x36\x14\xf6\xd5\x9c\x67\xf6\x03\x33\x07\x72\x46\xaf\xee\x47\x4f\x07\x3e\x08\x17\xfe\xa8\xb8\x3b\x3b\xed\x7f\x80\x8b\xc3\xc5\xa8\x3b\x96\x67\xf4\x4a\x38\xd5\x8e\x79\xb1\x0a\x8c\x0a\xdb\xdd\x4e\x9c\xa7\xa7\x34\xde\x49\x69\xdb\x37\x2b\xd8\x0c\x8b\x43\x03\xf3\x2e\xeb\x76\xbb\x21\xf2\x8e\x09\x29\xc3\xa7\x22\xd3\x5c\xa5\x76\xdc\x7f\x22\xa9\x40\xb8\x7d\x50\xef\xad\x7d\xd6\xf0\x2b\x78\xe4\x49\x24\xc7\x9b\x4d\x0b\x2b\x47\x20\xfa\x97\xbb\x42\xb1\x97\x4e\x37\x01\x53\xdc\x7a\x70\x98\xf0\xe4\x66\x9e\x0c\x03\x98\x22\xc6\xe4\x67\x00\x00\x00\xff\xff\x8e\xfd\x90\xeb\x64\x04\x00\x00")

func appengineTemplateFsSnippetTmplBytes() ([]byte, error) {
	return bindataRead(
		_appengineTemplateFsSnippetTmpl,
		"appengine/template/fs-snippet.tmpl",
	)
}

func appengineTemplateFsSnippetTmpl() (*asset, error) {
	bytes, err := appengineTemplateFsSnippetTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "appengine/template/fs-snippet.tmpl", size: 1124, mode: os.FileMode(420), modTime: time.Unix(1570400762, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _appengineTemplateRootTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\x90\xcd\x6e\xeb\x20\x10\x85\xf7\x79\x8a\xb9\xb3\x8f\xd9\x5e\xa9\xe0\x45\x7f\xb2\x6d\xa5\xb8\x8b\x2e\x09\x8c\xe4\xa9\x6c\x83\x60\xe4\xa8\x42\xbc\x7b\xe5\x92\xa8\x69\xcb\x06\x0e\x87\xef\x43\xa0\xff\x3d\x3e\x3f\x0c\x6f\x2f\x4f\x30\xca\x3c\xf5\x3b\x7d\x9d\xc8\xfa\x7e\x07\x00\xa0\x67\x12\x0b\x6e\xb4\x29\x93\x18\x7c\x1d\x0e\xfb\xff\x08\xea\xb6\x5c\xec\x4c\x06\x57\xa6\x73\x0c\x49\x10\x5c\x58\x84\x16\x31\x78\x66\x2f\xa3\xf1\xb4\xb2\xa3\xfd\x57\xc0\x0b\x27\x2c\x13\xf5\xa5\x74\xc3\xb6\xa8\x55\xab\xb6\xd3\xda\xec\x12\x47\x69\x61\x1b\x13\x09\x9c\x42\x90\x2c\xc9\x46\x30\x50\x4a\x77\x7f\x8d\xb5\xde\x35\x48\xdd\x52\x7f\x14\xa5\x08\xcd\x71\xb2\x42\x80\x87\xe3\x71\xe1\x18\x49\xb0\xd6\x5f\xac\x56\xed\xe5\xfa\x14\xfc\xc7\x45\xe5\x79\x05\xf6\x06\x6d\x8c\xd8\x6b\xe5\x79\xfd\x71\x07\xe4\xe4\x0c\x2a\xcf\x59\x54\x0a\x41\xba\xf7\xbc\x1d\xfb\x36\x36\x95\x56\xed\x6b\x3f\x03\x00\x00\xff\xff\xa3\x24\x6d\xc8\x72\x01\x00\x00")

func appengineTemplateRootTmplBytes() ([]byte, error) {
	return bindataRead(
		_appengineTemplateRootTmpl,
		"appengine/template/root.tmpl",
	)
}

func appengineTemplateRootTmpl() (*asset, error) {
	bytes, err := appengineTemplateRootTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "appengine/template/root.tmpl", size: 370, mode: os.FileMode(420), modTime: time.Unix(1570400762, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _appengineTemplateUTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x52\xbd\x6e\xe3\x30\x0c\xde\xf3\x14\x8c\xa6\xbb\x21\xd6\x7a\xc0\x49\x5e\xee\x9a\xb5\x05\x92\x0e\x1d\x55\x89\x89\x99\xfa\x47\x10\x19\x07\x81\xe1\x77\x2f\x1c\x25\x69\x90\xb6\x40\x39\xd8\xa2\xbe\x1f\x52\x94\xcc\xfc\xff\xe3\xbf\xf5\xcb\xd3\x03\x54\xd2\xd4\xe5\xcc\x5c\x7e\xe8\x42\x39\x03\x00\x30\x0d\x8a\x03\x5f\xb9\xc4\x28\x56\x3d\xaf\x97\x8b\x3f\x0a\xf4\x2d\xd8\xba\x06\xad\xea\x09\x0f\xb1\x4b\xa2\xc0\x77\xad\x60\x2b\x56\x1d\x28\x48\x65\x03\xf6\xe4\x71\x71\x4a\xd4\x59\x27\x24\x35\x96\xc3\x50\xac\xa7\xc5\x38\x1a\x9d\x77\x32\x5a\x53\xfb\x06\x09\x6b\xab\x58\x8e\x35\x72\x85\x28\x0a\xaa\x84\x1b\xab\x34\x8b\x13\xf2\x7a\xd3\xb5\x72\xfa\x14\x9e\xf9\x62\xcb\x3e\x51\x94\x9c\x4c\x51\xa3\xc0\xde\xa5\x2d\x83\x85\x61\x28\xc6\xf1\x6f\xe6\xe9\x5b\xe2\x27\x15\x6d\xe0\xd7\xfc\x24\x2b\x62\xa2\xde\x09\xfe\x86\xe1\x8a\x4e\x31\x0c\x82\x4d\xac\x9d\x20\xa8\xe5\x6a\xd5\x52\x8c\x28\x6a\x1c\xaf\xa4\xf1\xfb\x3a\xc0\xc9\x7f\x1c\x63\xc7\xba\x8d\xa9\xdb\x26\x64\x2e\x76\xac\xca\x3b\xcd\x4f\x46\xe1\xf9\xd6\x24\x8f\xc3\xe8\x7c\x85\xe6\xb5\x0b\xc7\xb3\x57\xa0\x1e\x28\x58\xe5\x62\x9c\xea\x04\xea\xbf\x6a\x2c\x10\x8b\xde\xdf\xf5\x62\x74\xf6\x31\x3a\x3f\x90\xf7\x00\x00\x00\xff\xff\xef\xcd\x61\xac\x38\x02\x00\x00")

func appengineTemplateUTmplBytes() ([]byte, error) {
	return bindataRead(
		_appengineTemplateUTmpl,
		"appengine/template/u.tmpl",
	)
}

func appengineTemplateUTmpl() (*asset, error) {
	bytes, err := appengineTemplateUTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "appengine/template/u.tmpl", size: 568, mode: os.FileMode(420), modTime: time.Unix(1539186795, 0)}
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
	"appengine/template/fs-snippet.tmpl": appengineTemplateFsSnippetTmpl,
	"appengine/template/root.tmpl": appengineTemplateRootTmpl,
	"appengine/template/u.tmpl": appengineTemplateUTmpl,
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
			"fs-snippet.tmpl": &bintree{appengineTemplateFsSnippetTmpl, map[string]*bintree{}},
			"root.tmpl": &bintree{appengineTemplateRootTmpl, map[string]*bintree{}},
			"u.tmpl": &bintree{appengineTemplateUTmpl, map[string]*bintree{}},
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

