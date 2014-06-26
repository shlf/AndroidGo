package main

import (
	"cfg"
	"flag"
	"fmt"
	"linkedlist"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
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
var publishFlag *string = flag.String("publish", RELEASE_PUBLISH, "Output release or debug edition.")
var tooleFlag *string = flag.String("tool", "", "Update project use Git(default) or other code store manager tools.")
var projectPathFlag *string = flag.String("ppath", "", "Update project path.")

var paths = new(linkedlist.LinkedList)

func Help() {
	fmt.Println("------------------------")
	fmt.Println("Usage :")
	fmt.Println("       -h : Show command help tip.")
	fmt.Println("       -c : Check your environment.")
	fmt.Println("       -v : Show command version number.")
	fmt.Println("       -p : Input android project path.")
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

func RunShell(command string) bool {
	cmd := exec.Command("/bin/sh", "-c", command)
	result, _err := cmd.Output()
	if _err != nil {
		fmt.Fprintf(os.Stderr, "The command failed to perform: %s (Command: %s) \n", _err, command)
		return false
	}

	if DEBUG {
		fmt.Fprintf(os.Stdout, "----run Result: %s (OK)\n", result)
	}

	return true
}

// update and build project
func Run(path string) bool {
	// run clean project
	cleanCmd := "cd " + path + " && ant clean"
	fmt.Println("----cleanCmd :", cleanCmd)
	if !RunShell(cleanCmd) {
		return false
	}

	// update project.
	rs := []rune(path)
	projectName := string(rs[strings.LastIndex(path, "/")+1:])
	fmt.Printf("----((((projectName:%s))))", projectName)
	updateCmd := exec.Command("android", "update", "project", "-p", path, "-n", projectName)
	updateResult, updateErr := updateCmd.Output()
	if updateErr != nil {
		fmt.Fprintf(os.Stderr, "The command failed to perform: %s (Command: android update project -p) \n", updateErr)
		return false
	}

	if DEBUG {
		fmt.Fprintf(os.Stdout, "----update Result: %s \n", updateResult)
	}

	// run ant to build apk or lib
	var publishCmd string
	if *publishFlag == RELEASE_PUBLISH {
		publishCmd = "cd " + path + " && ant release"
	} else if *publishFlag == DEBUG_PUBLISH {
		publishCmd = "cd " + path + " && ant debug"
	}
	fmt.Println("----publishCmd :", publishCmd)

	return RunShell(publishCmd)
}

// parse project.properties file in project
func ParseCfg(path string) {
	if DEBUG {
		fmt.Println("<<" + path + ">>")
	}
	// reserve last item, delete old pne
	_, ok := paths.Find(path)
	if ok {
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
		fmt.Println("----%v", cfgmap)
	}

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
	if RunShell(updateCmd) {
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
		Help()
		return
	}

	if *buildPath != "" {
		if DEBUG {
			fmt.Println("----Android project path :", *buildPath)
		}

		// delete end last '/' char
		rs := []rune(*buildPath)
		rl := len(rs)
		lastIndex := strings.LastIndex(*buildPath, "/")

		var path string
		if rl-1 == lastIndex {
			path = string(rs[:strings.LastIndex(*buildPath, "/")])
			if DEBUG {
				fmt.Printf("----((((path:%s))))----\n", path)
			}
		} else {
			path = *buildPath
		}

		// parse project.properties and to check depends.
		ParseCfg(path)

		// list paths to show
		var counter int32 = 0
		to_s := func(node *linkedlist.Node) {
			// convert item into string
			node.Value = reflect.ValueOf(node.Value).String()
			if str_v, ok := node.Value.(string); ok {
				if DEBUG {
					fmt.Println("<<< node.Value : " + str_v + " >>>")
				}
				counter++
			} else {
				fmt.Errorf("can't convert to strings %s, result %s", str_v, ok)
			}
		}
		paths.Map(to_s)
		fmt.Println("list count :", counter)

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
						if !Run(str_v) {
							fmt.Println("------------------------")
							fmt.Println("<<< Over, Failed : you have some problem >>>")
							return
						}
					} else {
						fmt.Errorf("can't convert to strings %s, result %s", str_v, ok)
					}
				}
			} else {
				fmt.Println("------------------------")
				fmt.Println("<<< Over, Success : Luck Dog! >>>")
				return
			}
		}
	} else {
		Help()
	}
}
