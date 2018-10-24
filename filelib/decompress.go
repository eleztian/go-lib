package filelib

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	PackTypeTAR   = "tar"
	PackTypeZIP   = "zip"
	PackTypeTARGZ = "tar.gz"
)

func DeCompress(srcFile string, dstDir string, srcType string) error {
	// 查看目标文件夹是否存在，如果存在是否为空文件夹
	_, err := os.Stat(srcFile)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(srcFile, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		os.RemoveAll(dstDir)
		err = os.MkdirAll(dstDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	switch srcType {
	case PackTypeTAR:
		err = deCompressTar(srcFile, dstDir)
	case PackTypeZIP:
		err = deCompressZip(srcFile, dstDir)
	case PackTypeTARGZ:
		err = decompressTarGz(srcFile, dstDir)
	default:
		err = fmt.Errorf("not support this file type: %s", srcType)
	}
	if err != nil {
		os.RemoveAll(dstDir) // 移除,回退
	}
	return err
}

func deCompressTar(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	tr := tar.NewReader(srcFile)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := filepath.Join(dest, hdr.Name)
		d, f := filepath.Split(filename)
		if hdr.FileInfo().IsDir() {
			os.MkdirAll(filename, hdr.FileInfo().Mode())
		} else {
			os.MkdirAll(d, hdr.FileInfo().Mode())
			if f != "" {
				file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, hdr.FileInfo().Mode())
				if err != nil {
					return err
				}
				_, err = io.Copy(file, tr)
				if err != nil {
					file.Close()
					return err
				}
				file.Close()
			}
		}
	}
	return nil
}

func deCompressZip(zipFile, dest string) error {
	zr, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zr.Close()
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		name := f.Name
		if f.Flags == 0 { // 使用GBK编码
			name, err = GbkToUtf8([]byte(name))
		}

		filename := filepath.Join(dest, name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filename, f.Mode()); err != nil {
				return err
			}
			continue
		}
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(file, rc)
		if err != nil {
			file.Close()
			rc.Close()
			return err
		}
		file.Close()
		rc.Close()
	}
	return nil
}

func decompressTarGz(gzipFile, dest string) error {
	sf, err := os.Open(gzipFile)
	if err != nil {
		return err
	}
	defer sf.Close()
	gr, err := gzip.NewReader(sf)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := filepath.Join(dest, hdr.Name)
		d, f := filepath.Split(filename)
		if hdr.FileInfo().IsDir() {
			os.MkdirAll(filename, hdr.FileInfo().Mode())
		} else {
			os.MkdirAll(d, hdr.FileInfo().Mode())
			if f != "" {
				file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, hdr.FileInfo().Mode())
				if err != nil {
					return err
				}
				_, err = io.Copy(file, tr)
				if err != nil {
					file.Close()
					return err
				}
				file.Close()
			}
		}
	}
	return nil
}

func GbkToUtf8(s []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return "", e
	}
	return string(d), nil
}

func Utf8ToGbk(s []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return "", e
	}
	return string(d), nil
}

// GetSupportPackType 返回fiellib包支持的文件格式
func GetSupportPackType() []string {
	return []string{"tar", "tar.gz", "zip"}
}
