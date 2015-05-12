//////////////////////////////////////////
// archivex.go
// Jhonathan Paulo Banczek - 2014
// jpbanczek@gmail.com - jhoonb.com
//////////////////////////////////////////

package archivex

import (
	"archive/tar"
	"archive/zip"
	"io/ioutil"
	"os"
	"strings"
	"log"
)

// interface
type Archivex interface {
	Create(name string) error
	Add(name string, file []byte) error
	AddFile(name string) error
	AddAll(dir string, includeCurrentFolder bool) error
	Close() error
}

// ZipFile implement *zip.Writer
type ZipFile struct {
	Writer *zip.Writer
	Name   string
}

// TarFile implement *tar.Writer
type TarFile struct {
	Writer *tar.Writer
	Name   string
}

func isDir(path string) bool {
	src, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	return src.IsDir()
}

// Create new file zip
func (z *ZipFile) Create(name string) error {
	// check extension .zip
	if strings.HasSuffix(name, ".zip") != true {
		if strings.HasSuffix(name, ".tar.gz") == true {
			name = strings.Replace(name, ".tar.gz", ".zip", -1)
		} else {
			name = name + ".zip"
		}
	}
	z.Name = name
	file, err := os.Create(z.Name)
	if err != nil {
		return err
	}
	z.Writer = zip.NewWriter(file)
	return nil
}

// Add add byte in archive zip
func (z *ZipFile) Add(name string, file []byte) error {

	iow, err := z.Writer.Create(name)
	if err != nil {
		return err
	}
	_, err = iow.Write(file)
	return err
}

// AddFile add file from dir in archive
func (z *ZipFile) AddFile(name string) error {
	bytearq, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	filep, err := z.Writer.Create(name)
	if err != nil {
		return err
	}
	_, err = filep.Write(bytearq)
	if err != nil {
		return err
	}
	return nil
}

// AddAll add all files from dir in archive
func (z *ZipFile) AddAll(dir string, includeCurrentFolder bool) error {
    if includeCurrentFolder { if err := z.addDir(getSubDir(dir, dir, includeCurrentFolder)); err != nil { return err } }
	return z.addAll(dir, dir, includeCurrentFolder)
}

func (z *ZipFile) addAll(dir string, rootDir string, includeCurrentFolder bool) error {
	// capture all name files in dir
	listFile, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	var names []string
	var bdatas [][]byte
    
	for _, arq := range listFile {
		if isDir(dir + arq.Name()) {
		    if err = z.addDir(getSubDir(dir + arq.Name(), rootDir, includeCurrentFolder)+"/"); err != nil { return err }
			z.addAll(dir + arq.Name() + string(os.PathSeparator), rootDir, includeCurrentFolder)
		} else {
			bytearq, err := ioutil.ReadFile(dir + arq.Name())
			if err != nil {
				return err
			}
			names = append(names, arq.Name())
			bdatas = append(bdatas, bytearq)
		}
	}

	subDir := getSubDir(dir, rootDir, includeCurrentFolder)
	for i, file := range bdatas {
		filep, err := z.Writer.Create(subDir + names[i])
		if err != nil {
			return err
		}
		_, err = filep.Write(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (z *ZipFile) addDir(name string) error {
    _, err := z.Writer.Create(name)
    return err
}

func (z *ZipFile) Close() error {
	err := z.Writer.Close()
	return err
}

// Create new Tar file
func (t *TarFile) Create(name string) error {
	// check extension .zip
	if strings.HasSuffix(name, ".tar.gz") != true {
		if strings.HasSuffix(name, ".zip") == true {
			name = strings.Replace(name, ".zip", ".tar.gz", -1)
		} else {
			name = name + ".tar.gz"
		}
	}
	t.Name = name
	file, err := os.Create(t.Name)
	if err != nil {
		return err
	}
	t.Writer = tar.NewWriter(file)
	return nil
}

// Add add byte in archive tar
func (t *TarFile) Add(name string, file []byte) error {

	hdr := &tar.Header{Name: name, Size: int64(len(file))}
	if err := t.Writer.WriteHeader(hdr); err != nil {
		return err
	}
	_, err := t.Writer.Write(file)
	return err
}

// AddFile add file from dir in archive tar
func (t *TarFile) AddFile(name string) error {
	bytearq, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	hdr := &tar.Header{Name: name, Size: int64(len(bytearq))}
	err = t.Writer.WriteHeader(hdr)
	if err != nil {
		return err
	}
	_, err = t.Writer.Write(bytearq)
	if err != nil {
		return err
	}
	return nil

}

// AddAll add all files from dir in archive
func (t *TarFile) AddAll(dir string, includeCurrentFolder bool) error {
    if includeCurrentFolder {
        if err := t.addDir(getSubDir(dir, dir, includeCurrentFolder)); err != nil { return err }
    }
	return t.addAll(dir, dir, includeCurrentFolder)
}

func (t *TarFile) addAll(dir string, rootDir string, includeCurrentFolder bool) error {

	// capture all name files in dir
	listFile, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
    
	var names []string
	var bdatas [][]byte

	for _, arq := range listFile {
		if isDir(dir + arq.Name()) {
		    if err = t.addDir(getSubDir(dir + arq.Name(), dir, includeCurrentFolder)+"/"); err != nil { return  err }
			t.addAll(dir + arq.Name() + string(os.PathSeparator), rootDir, includeCurrentFolder)
		} else {
			bytearq, err := ioutil.ReadFile(dir + arq.Name())
			if err != nil {
				return err
			}
			names = append(names, arq.Name())
			bdatas = append(bdatas, bytearq)
		}
	}

	subDir := getSubDir(dir, rootDir, includeCurrentFolder)
	for i, file := range bdatas {
		hdr := &tar.Header{Name: subDir + names[i], Size: int64(len(file))}
		if err := t.Writer.WriteHeader(hdr); err != nil {
			return err
		}
		_, err = t.Writer.Write(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TarFile) addDir(name string) error {
    h := &tar.Header{Name: name, Size: 0, Typeflag: '5'}
    err := t.Writer.WriteHeader( h )
    return err
}

// Close the file Tar
func (t *TarFile) Close() error {
	err := t.Writer.Close()
	return err
}

func getSubDir(dir string, rootDir string, includeCurrentFolder bool) string {
    subDir := strings.Replace(dir, rootDir, "", 1)
            
    if includeCurrentFolder {
        rootDirParts := strings.Split(rootDir, string(os.PathSeparator))
        subDir = rootDirParts[len(rootDirParts)-2] + string(os.PathSeparator) + subDir
    }
            
    return subDir
}
