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

const APP_VERSION = "0.1"
const CFG_NAME = "project.properties"
const REFERENCE = "android.library.reference."

const DEBUG = false

// The flag package provides a default help printer via -h switch
var helpFlag *bool = flag.Bool("h", false, "Print the command help.")
var checkFlag *bool = flag.Bool("c", false, "Check your environment.")
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
var buildPath *string = flag.String("p", "", "Input android project path.")

var paths = new(linkedlist.LinkedList)

func Help() {
	fmt.Println("------------------------")
	fmt.Println("Usage :")
	fmt.Println("       -h : show command help tip.")
	fmt.Println("       -c : check your environment.")
	fmt.Println("       -v : show command version number.")
	fmt.Println("       -p : input android project path.")
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

func RunAnt(command string) {
	cmd := exec.Command("/bin/sh", "-c", command)
	result, _err := cmd.Output()
	if _err != nil {
		if DEBUG {
			fmt.Fprintf(os.Stderr, "The command failed to perform: %s (Command: %s) \n", _err, command)
		}
		return
	}

	if DEBUG {
		fmt.Fprintf(os.Stdout, "----run Result: %s \n", result)
	}
}

// update and build project
func Run(path string) {
	// update project.
	updateCmd := exec.Command("android", "update", "project", "-p", path)
	updateResult, updateErr := updateCmd.Output()
	if updateErr != nil {
		if DEBUG {
			fmt.Fprintf(os.Stderr, "The command failed to perform: %s (Command: android update project -p) \n", updateErr)
		}
		return
	}
	fmt.Fprintf(os.Stdout, "----update Result: %s \n", updateResult)

	// run clean project
	cleanPath := "cd " + path + " && ant clean"
	fmt.Println("----cleanPath:", cleanPath)
	RunAnt(cleanPath)

	// run ant release to build apk or lib
	releasePath := "cd " + path + " && ant release"
	fmt.Println("----releasePath:", releasePath)
	RunAnt(releasePath)
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

func main() {
	showCopyright()
	flag.Parse() // Scan the arguments list

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

		ParseCfg(*buildPath)

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

		// build project and depend project.
		for {
			if counter >= 0 {
				counter--
				node, _err := paths.Get(counter)
				if node != nil && _err == nil {
					// find item
					// convert item into string
					node.Value = reflect.ValueOf(node.Value).String()
					if str_v, ok := node.Value.(string); ok {
						Run(str_v)
					} else {
						fmt.Errorf("can't convert to strings %s, result %s", str_v, ok)
					}
				}
			} else {
				fmt.Println("OOOOOOO Over :", counter)
				return
			}
		}
	} else {
		Help()
	}
}
