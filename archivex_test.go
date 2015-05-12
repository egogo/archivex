//////////////////////////////////////////
// archivex_test.go
// Jhonathan Paulo Banczek - 2014
// jpbanczek@gmail.com - jhoonb.com
//////////////////////////////////////////

package archivex

import (
	"testing"
    "github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"bytes"
	"strings"
	"regexp"
)

func Test_archivex(t *testing.T) {

	// interface
	arcvx := []Archivex{&ZipFile{}, &TarFile{}}
    extensions := []string{"zip", "tar.gz"}
    for _, ext := range extensions {
        os.Remove("filetest_absolute_without_root."+ext)
    	os.Remove("filetest_absolute_with_root."+ext)
    	os.Remove("filetest_relative_without_root."+ext)
    	os.Remove("filetest_relative_with_root."+ext)
    }

    dir, _ := os.Getwd()
    
	for idx, arc := range arcvx {
	    ext := extensions[idx]
		// Absolute, no root
        err := arc.Create("filetest_absolute_without_root." + ext)
        checkError(t, err)
        err = arc.AddAll(dir+"/testfolder/", false)
        arc.Close()
        entries, err := getEntries("filetest_absolute_without_root." + ext)
        checkError(t, err)
        expected := []string{"1/", "1/bla.txt", "test1.txt", "test2.txt"}
        assert.Equal(t, expected, entries, ext + ": Must include files and containing directories.")
		
		// Absolute with root
        err = arc.Create("filetest_absolute_with_root." + ext)
        checkError(t, err)
        err = arc.AddAll(dir+"/testfolder/", true)
        arc.Close()
        entries, err = getEntries("filetest_absolute_with_root." + ext)
        checkError(t, err)
        expected = []string{"testfolder/", "testfolder/1/", "testfolder/1/bla.txt", "testfolder/test1.txt", "testfolder/test2.txt"}
        assert.Equal(t, expected, entries, ext + ": Must include files and containing directories.")
        
        // Relative, no root
        err = arc.Create("filetest_relative_without_root." + ext)
        checkError(t, err)
        err = arc.AddAll("testfolder/", false)
        arc.Close()
        entries, err = getEntries("filetest_relative_without_root." + ext)
        checkError(t, err)
        expected = []string{"1/", "1/bla.txt", "test1.txt", "test2.txt"}
        assert.Equal(t, expected, entries,  ext +": Must include files and containing directories.")
        
        // Relative with root
        err = arc.Create("filetest_relative_with_root." + ext)
        checkError(t, err)
        err = arc.AddAll("testfolder/", true)
        arc.Close()
        entries, err = getEntries("filetest_relative_with_root." + ext)
        checkError(t, err)
        expected = []string{"testfolder/", "testfolder/1/", "testfolder/1/bla.txt", "testfolder/test1.txt", "testfolder/test2.txt"}
        assert.Equal(t, expected, entries, ext + ": Must include files and containing directories.")
	}
	
	for _, ext := range extensions {
        os.Remove("filetest_absolute_without_root."+ext)
    	os.Remove("filetest_absolute_with_root."+ext)
    	os.Remove("filetest_relative_without_root."+ext)
    	os.Remove("filetest_relative_with_root."+ext)
    }
    
}

// func for check errors
func checkError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func getEntries(filename string) ([]string, error) {
    if(strings.HasSuffix(filename, ".zip")){ return getZipEntries(filename) }
    return getTarEntries(filename)
}

func getTarEntries(filename string) ([]string, error) {
    out, err := runCmd("tar", "-tf", filename)
	lines := strings.Split(out, "\n")
	entries := lines[0:len(lines)-1]
	return entries, err
}

func getZipEntries(filename string) ([]string, error) {
    r, _ := regexp.Compile("\\s+")
    
    out, err := runCmd("zipinfo", filename)
	lines := strings.Split(out, "\n")
	var entries []string
	for _, line := range lines[1:len(lines)-2] {
	    parts := r.Split(line, -1)
	    entries = append(entries, parts[len(parts)-1])
	}
	return entries, err
}

func runCmd(command string, args ...string) (string, error) {
    cmd := exec.Command(command, args...)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    return out.String(), err
}
