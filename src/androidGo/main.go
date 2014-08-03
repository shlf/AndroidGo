package main

import (
    "cfg"
    "flag"
    "fmt"
    "linkedlist"
    "log"
    "originbuild"
    "os/exec"
    "reflect"
    "runtime"
    "strings"
    "time"
)

const APP_VERSION = "0.2"
const CFG_NAME = "project.properties"
const REFERENCE = "android.library.reference."
const RELEASE_PUBLISH = "release"
const DEBUG_PUBLISH = "debug"
const TOOL_GIT = "git"

const DEBUG = true

// The flag package provides a default help printer via -h switch
var helpFlag *bool = flag.Bool("h", false, "Print the command help.")
var checkFlag *bool = flag.Bool("c", false, "Check your environment.")
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
var buildPath *string = flag.String("p", "", "Input android project path.")
var sdkPath *string = flag.String("sdk", "", "Android SDK path.")
var publishFlag *string = flag.String("publish", RELEASE_PUBLISH, "Output release or debug edition.")
var tooleFlag *string = flag.String("tool", "", "Update project use Git(default) or other code store manager tools.")
var projectPathFlag *string = flag.String("ppath", "", "Update project path.")

var paths = new(linkedlist.LinkedList)
var apksCfg = make(map[string]*originbuild.Apk)
var runtineApksCfg = make(map[string]*originbuild.Apk)
var startTime time.Time

func Help() {
    fmt.Println("------------------------")
    fmt.Println("Usage :")
    fmt.Println("       -h : Show command help tip.")
    fmt.Println("       -c : Check your environment.")
    fmt.Println("       -v : Show command version number.")
    fmt.Println("       -p : Input android project path.")
    fmt.Println("       -sdk : Input android SDK path.")
    fmt.Println("       -publish : Choice release(default) or debug mode to publish.")
    fmt.Println("       -tool : Update project use Git(default) or other code store manager tools.")
    fmt.Println("       -ppath : Update project path.")
    fmt.Println("------------------------")
}

func showCopyright() {
    fmt.Println("This Command Information:")
    fmt.Println("  Author : SHLF(shenlongfei)")
    fmt.Println("  Date   : 2014-06-24")
    fmt.Println("  Email  : shenlongfei1989@gmail.com")
}

func checkEnv() bool {
    javaResult := false
    javaCmd := exec.Command("java", "-version")
    _, javaErr := javaCmd.Output()
    if javaErr != nil {
        fmt.Println("java environment is bad.")
    } else {
        fmt.Println("java environment is good.")
        javaResult = true
    }

    androidResult := false
    androidCmd := exec.Command("android", "list")
    _, androidErr := androidCmd.Output()
    if androidErr != nil {
        fmt.Println("android SDK environment is bad.")
    } else {
        fmt.Println("android SDK environment is good.")
        androidResult = true
    }

    antResult := false
    antCmd := exec.Command("ant", "-version")
    _, antErr := antCmd.Output()
    if antErr != nil {
        fmt.Println("ant environment is bad.")
    } else {
        fmt.Println("ant environment is good.")
        antResult = true
    }

    return (javaResult && androidResult && antResult)
}

// parse project.properties file in project
func ParseCfg(path string) {
    if DEBUG {
        fmt.Println("<<" + path + ">>")
    }
    // reserve last item, delete old pne
    if _, ok := paths.Find(path); ok {
        paths.Remove(path)
    }
    paths.Append(path)

    cfgPath := path + "/" + CFG_NAME
    cfgmap := make(map[string]string)
    err := cfg.Load(cfgPath, cfgmap)
    if err != nil {
        log.Fatal(err)
    }

    if DEBUG {
        fmt.Println("----map =", cfgmap)
    }

    //================================
    thisApk := new(originbuild.Apk)
    thisApk.Path = path

    rs := []rune(path)
    projectName := string(rs[strings.LastIndex(path, "/")+1:])
    thisApk.ProjectName = projectName

    for key, value := range cfgmap {
        if strings.Contains(key, "android.library") {
            if value == "true" {
                thisApk.IsRealLibrary = true
                thisApk.IsLibrary = true
            } else {
                thisApk.IsRealLibrary = false
                thisApk.IsLibrary = false
            }
        }

        if strings.Contains(key, REFERENCE) {
            thisApk.IsRealLibrary = false
        }

        if strings.Contains(key, "target") {
            thisApk.Api = value
        }

        if strings.Contains(key, "key.store") {
            thisApk.Keystore = value
        }
        if strings.Contains(key, "key.alias") {
            thisApk.Keyalias = value
        }
        if strings.Contains(key, "key.store.password") {
            thisApk.Kstorepass = value
        }
        if strings.Contains(key, "key.alias.password") {
            thisApk.Keypass = value
        }
    }

    // find need goruntine project
    if thisApk.IsRealLibrary {
        runtineApksCfg[path] = thisApk
    } else {
        apksCfg[path] = thisApk
    }
    //================================

    for key, value := range cfgmap {
        if strings.Contains(key, REFERENCE) {
            var dependProjectPath string = ""

            count := strings.Count(value, "../")
            if count > 0 {
                // hand depend project path
                pathPice := strings.Split(path, "/")
                endIndex := len(pathPice) - count
                newPath := pathPice[:endIndex]

                var rootPath string
                for _, vpath := range newPath {
                    rootPath += vpath + "/"
                }
                if DEBUG {
                    fmt.Println("----rootPath =", rootPath)
                }

                endDir := strings.Trim(value, "../")
                if DEBUG {
                    fmt.Println("----endDir =", endDir)
                }

                dependProjectPath = rootPath + endDir
            } else {
                if strings.Index(value, "/") == 0 {
                    dependProjectPath = value
                } else {
                    dependProjectPath = path + "/" + value
                }
            }

            if dependProjectPath != "" {
                fmt.Println("---->> dependProjectPath =", dependProjectPath)
                // check depend project depend
                ParseCfg(dependProjectPath)
            }
        }
    }
}

func gitUpdate() bool {
    updateCmd := "cd " + *projectPathFlag + " && git pull"
    fmt.Println("----updateCmd :", updateCmd)
    if originbuild.RunShell(updateCmd) {
        fmt.Println("Update code OK!")
        return true
    } else {
        fmt.Println("Please check your environment.")
        return false
    }
}

func main() {
    showCopyright()
    flag.Parse() // Scan the arguments list

    // get android sdk path.
    sdkpath := originbuild.DeleteEndChar(*sdkPath, "/")

    // hand code update logic
    if *projectPathFlag != "" {
        if *tooleFlag == "" {
            fmt.Println("Please input your code manager tool, like git or svn or others.")
            return
        } else {
            if *tooleFlag == TOOL_GIT {
                gitUpdate()
            } else {
                fmt.Println("Sorry! Not support", *tooleFlag)
            }
        }
        return
    } else {
        if *tooleFlag != "" {
            fmt.Println("Please input tool name, like git or svn or others")
            return
        }
    }

    if *checkFlag {
        if checkEnv() {
            fmt.Println("------------------------")
            fmt.Println("Everything is OK.")
        } else {
            fmt.Println("------------------------")
            fmt.Println("Please to check you environment (java, android, ant).")
        }
        return
    }

    if *versionFlag {
        fmt.Println("Version:", APP_VERSION)
        return
    }

    if *helpFlag {
        goto HELP
        return
    }

    if *buildPath != "" {
        if DEBUG {
            fmt.Println("----Android project path :", *buildPath)
        }

        path := originbuild.DeleteEndChar(*buildPath, "/")

        // parse project.properties and to check depends.
        ParseCfg(path)
        startTime = time.Now()

        // go apksCfg
        NCPU := runtime.NumCPU()
        fmt.Println("MMMMMMMM CPU Number =", NCPU)
        runtime.GOMAXPROCS(NCPU)
        ch := make(chan int, NCPU)

        // built basic lib
        goCount := 0
        for _, value := range runtineApksCfg {
            goCount++

            if DEBUG {
                fmt.Println("<<< runtime Path : "+value.Path+" >>>count:", goCount)
                fmt.Println("<<< runtime ProjectName : " + value.ProjectName)
                fmt.Println("<<< runtime Api : " + value.Api)
                fmt.Println("<<< runtime IsRealLibrary :", value.IsRealLibrary)

                fmt.Println("<<< runtime Keystore : " + value.Keystore)
                fmt.Println("<<< runtime Keyalias : " + value.Keyalias)
                fmt.Println("<<< runtime Kstorepass : " + value.Kstorepass)
                fmt.Println("<<< runtime Keypass : " + value.Keypass)
            }
            if _, ok := paths.Find(value.Path); ok {
                paths.Remove(value.Path)
            }

            go originbuild.RuntineBuild(value, sdkpath, ch)
        }

        result := 0
        for i := 0; i < goCount; i++ {
            temp := <-ch
            fmt.Println("flag :", temp)
            result += temp
        }
        fmt.Println("flag count :", result, "(goruntime count)")
        if result != goCount {
            goto FAILED
        }
        fmt.Println("Built basic lib success.")

        // list paths to show
        var counter int32 = 0
        to_s := func(node *linkedlist.Node) {
            // convert item into string
            node.Value = reflect.ValueOf(node.Value).String()
            if str_v, ok := node.Value.(string); ok {
                if DEBUG {
                    fmt.Println("<<< one by one to build : " + str_v + " >>>")
                }
                counter++
            } else {
                fmt.Errorf("can't convert to strings %s, result %s", str_v, ok)
            }
        }
        paths.Map(to_s)
        fmt.Println("one by one to build list count :", counter)

        counterTmp := counter
        // build project and depend project.
        var index int32 = 0
        for {
            if counter >= 0 {
                counter--
                index++
                node, _err := paths.Get(counter)
                if index <= counterTmp {
                    fmt.Printf("------------start build [%d]------------\n", index)
                }

                if node != nil && _err == nil {
                    // find item
                    // convert item into string
                    node.Value = reflect.ValueOf(node.Value).String()
                    if str_v, ok := node.Value.(string); ok {
                        fmt.Println("----------path =", str_v)
                        if currentApk, result := apksCfg[str_v]; result {
                            if !originbuild.RunBuild(currentApk, sdkpath) {
                                goto FAILED
                                return
                            }
                        }
                    } else {
                        fmt.Errorf("can't convert to strings %s, result %s", str_v, ok)
                    }
                }
            } else {
                endTime := time.Now()
                fmt.Println("------------------------")
                fmt.Printf("Expend Time(%d)\n", endTime.Sub(startTime))
                fmt.Println("<<< Over, Success : Luck Dog! >>>")
                return
            }
        }
    } else {
        goto FAILED
    }

HELP:
    Help()

FAILED:
    endTime := time.Now()
    fmt.Println("------------------------")
    fmt.Printf("-----Expend Time(%d)------\n", endTime.Sub(startTime))
    fmt.Println("<<< Over, Failed : you have some problem >>>")
}
