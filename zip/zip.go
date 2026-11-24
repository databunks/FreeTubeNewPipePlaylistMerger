package zip
	
import (
    "archive/zip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
    "strings"
)

func DecompressFolder(path string) {
    reader, err := zip.OpenReader(path)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer reader.Close()

    for _, file := range reader.File {
        fmt.Printf("Extracting file: %s\n", file.Name)
        err := ExtractFile(file)
        if err != nil {
            fmt.Println(err)
        }
    }
}

func CompressFiles(path string){
    timestamp := time.Now().Format("20060102150405")
    newZipFile, err := os.Create("mergedPlaylists_" + timestamp + ".zip")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer newZipFile.Close()

    zipWriter := zip.NewWriter(newZipFile)
    defer zipWriter.Close()


    AddFile(zipWriter, "./" + path + "/newpipe.db")
    AddFile(zipWriter, "./" + path + "/newpipe.settings")
    AddFile(zipWriter, "./" + path + "/preferences.json")
}

func ExtractFile(f *zip.File) error {
    rc, err := f.Open()
    if err != nil {
        return err
    }
    defer rc.Close()

    path := filepath.Join("output", f.Name)
    if f.FileInfo().IsDir() {
        os.MkdirAll(path, os.ModePerm)
    } else {
        os.MkdirAll(filepath.Dir(path), os.ModePerm)
        outFile, err := os.Create(path)
        if err != nil {
            return err
        }
        defer outFile.Close()

        _, err = io.Copy(outFile, rc)
    }
    return err
}

func AddFile(zipWriter *zip.Writer, filename string) {
    fileToZip, err := os.Open(filename)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer fileToZip.Close()

    parsedPath := strings.Split(filename, "/")
    w, err := zipWriter.Create(parsedPath[len(parsedPath) - 1])
    if err != nil {
        fmt.Println(err)
        return
    }

    if _, err := io.Copy(w, fileToZip); err != nil {
        fmt.Println(err)
    }
}

func RemoveFolder(path string){
	err := os.RemoveAll(path)

	if (err != nil){
		panic(err)
	}
}

