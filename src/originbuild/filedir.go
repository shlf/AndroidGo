package originbuild

import (
    "container/list"
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "strings"
)

const (
    REGULAR_JAVA = "\\w+.java"
    REGULAR_AIDL = "\\w+.aidl"
)

func GetFilesPath(fullPath, regular string) []string {
    listStr := list.New()
//    var javaPath string

    filepath.Walk(fullPath, func(path string, fi os.FileInfo, err error) error {

        if nil == fi {
            return err
        }
        if fi.IsDir() {
            return nil
        }

        name := fi.Name()

        reg := regexp.MustCompile(regular)

        if reg.MatchString(name) {
            fmt.Println("<<< path : " + path)
//            javaPath += path + " "
            listStr.PushBack(path)
        }

        return nil
    })

//    return javaPath
    sli := []string{}
    for el := listStr.Front(); nil != el; el = el.Next() {
        sli = append(sli, el.Value.(string))
    }
    fmt.Printf("<<< sli : %s\n", sli)
    return sli
}

func DeleteEndChar(name, char string) string {
    // delete end last '/' char
    rs := []rune(name)
    rl := len(rs)
    lastIndex := strings.LastIndex(name, char)

    var path string
    if rl-1 == lastIndex {
        path = string(rs[:strings.LastIndex(name, char)])
        if DEBUG {
            fmt.Printf("----((((path:%s))))----\n", path)
        }
    } else {
        path = name
    }

    return path
}

func CheckPath(path string) bool {
    dir, err := os.Stat(path)
    if err == nil {
        if dir.IsDir() {
            return true
        }
    }

    return false
}

func CheckFile(path string) bool {
    f, err := os.Open(path)
    if err != nil && os.IsNotExist(err) {
        fmt.Printf("file not exist!\n")
        return false
    }
    fmt.Printf("file exist!\n")
    defer f.Close()

    return true
}
