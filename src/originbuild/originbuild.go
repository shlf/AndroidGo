package originbuild

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

const DEBUG = true

/**
步骤中提到的工具如下表：
名称         功能介绍                               在操作系统中的路径
aapt        Android资源打包工具                     ${ANDROID_SDK_HOME}/platform-tools/appt
aidl        Android接口描述语言转化为.java文件的工具   ${ANDROID_SDK_HOME}/platform-tools/aidl
javac       Java Compiler                         ${JDK_HOME}/javac或/usr/bin/javac
dex         转化.class文件为Davik VM能识别的.dex文件  ${ANDROID_SDK_HOME}/platform-tools/dx
apkbuilder  生成apk包                              ${ANDROID_SDK_HOME}/tools/apkbuilder
jarsigner   .jar文件的签名工具                       ${JDK_HOME}/jarsigner或/usr/bin/jarsigner
zipalign    字节码对齐工具                           ${ANDROID_SDK_HOME}/tools/zipalign
*/

type Apk struct {
    Path          string
    ProjectName   string
    Api           string
    IsLibrary     bool
    IsRealLibrary bool

    Keystore   string
    Keyalias   string
    Kstorepass string
    Keypass    string
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

func RunCommand(args []string) (bool, error) {
    rCmd := &exec.Cmd{
        Path: args[0],
        Args: args,
    }
    if filepath.Base(args[0]) == args[0] {
        if lp, err := exec.LookPath(args[0]); err != nil {
            //rCmd.lookPathErr = err
        } else {
            rCmd.Path = lp
        }
    }

    rResult, rErr := rCmd.Output()

    if DEBUG {
        fmt.Fprintf(os.Stdout, "----Command: %s\n", rCmd.Args)
        fmt.Fprintf(os.Stdout, "----run Result: %s \n", rResult)
    }

    if rErr == nil {
        return true, nil
    } else {
        fmt.Fprintf(os.Stderr, "The command failed to perform: %s\n", rErr)
        return false, rErr
    }
}

// R.java
func RunRClass(apk *Apk, sdkPath string) bool {
    //check gen dir
    genPath := apk.Path + "/gen"
    if !CheckPath(genPath) {
        os.Mkdir(genPath, os.FileMode(0777))
    }

    args := []string{}
    args = append(args, "aapt")
    args = append(args, "package")
    args = append(args, "-f")
    args = append(args, "-m")
    args = append(args, "-M")
    args = append(args, apk.Path+"/AndroidManifest.xml")
    args = append(args, "-J")
    args = append(args, genPath)
    args = append(args, "-S")
    args = append(args, apk.Path+"/res")
    args = append(args, "-I")
    args = append(args, sdkPath+"/platforms/"+apk.Api+"/android.jar")

    if result, _ := RunCommand(args); result {
        fmt.Println("Built R.java success.")
        return true
    } else {
        fmt.Println("Built R.java failed.")
        return false
    }
}

// AIDI
func RunAIDL(apk *Apk, sdkPath string) bool {
    //check gen dir
    genPath := apk.Path + "/gen"
    if !CheckPath(genPath) {
        os.Mkdir(genPath, os.FileMode(0777))
    }

    aidlFiles := GetFilesPath(apk.Path, REGULAR_AIDL)
    files := len(aidlFiles)
    if len(aidlFiles) > 0 {
        name := "aidl"

        counter := 0
        for _, value := range aidlFiles {
            args := []string{}
            args = append(args, name)
            args = append(args, "-I"+apk.Path+"/src")
            args = append(args, "-p"+sdkPath+"/platforms/"+apk.Api+"/framework.aidl")
            args = append(args, "-o"+genPath)
            args = append(args, value)

            if result, _ := RunCommand(args); result {
                fmt.Fprintf(os.Stderr, "The %s file built java success\n", value)
                counter++
            } else {
                fmt.Fprintf(os.Stderr, "The %s file built java failed\n", value)
            }
        }
        if counter == files {
            fmt.Println("Built AIDL success.")
            return true
        } else {
            fmt.Println("Built AIDL failed.")
            return false
        }
    } else {
        fmt.Fprintf(os.Stdout, "----project: %s not contain aidl\n", apk.ProjectName)
        return true
    }
}

// classes
func RunClass(apk *Apk, sdkPath string) bool {
    binPath := apk.Path + "/bin/classes"
    if !CheckPath(binPath) {
        os.MkdirAll(binPath, os.FileMode(0777))
    }

    javaFiles := GetFilesPath(apk.Path, REGULAR_JAVA)

    args := []string{}
    args = append(args, "javac")
    args = append(args, "-encoding")
    args = append(args, "GB18030")
    args = append(args, "-target")
    args = append(args, "1.6")
    args = append(args, "-bootclasspath")
    args = append(args, sdkPath+"/platforms/"+apk.Api+"/android.jar")
    args = append(args, "-d")
    args = append(args, binPath)

    for _, value := range javaFiles {
        args = append(args, value)
    }

    if result, _ := RunCommand(args); result {
        fmt.Println("Built all classes success.")
        return true
    } else {
        fmt.Println("Built all classes failed.")
        return false
    }
}

// DEX
func RunDEX(path string) bool {
    args := []string{}
    args = append(args, "dx")
    args = append(args, "--dex")
    args = append(args, "--output="+path+"/bin/classes.dex")
    args = append(args, path+"/bin/classes")

    if result, _ := RunCommand(args); result {
        fmt.Println("Built DEX success.")
        return true
    } else {
        fmt.Println("Built DEX failed.")
        return false
    }
}

// Resource
func RunResource(path, api, sdkPath string) bool {
    assetsPath := path + "/assets"
    if !CheckPath(assetsPath) {
        os.Mkdir(assetsPath, os.FileMode(0777))
    }

    args := []string{}
    args = append(args, "aapt")
    args = append(args, "package")
    args = append(args, "-f")
    args = append(args, "-M")
    args = append(args, path+"/AndroidManifest.xml")
    args = append(args, "-S")
    args = append(args, path+"/res")
    args = append(args, "-A")
    args = append(args, assetsPath)
    args = append(args, "-I")
    args = append(args, sdkPath+"/platforms/"+api+"/android.jar")
    args = append(args, "-F")
    args = append(args, path+"/bin/resources.ap_")

    if result, _ := RunCommand(args); result {
        fmt.Println("Built Resource success.")
        return true
    } else {
        fmt.Println("Built Resource failed.")
        return false
    }
}

// unsigner apk
func RunUnsignAPK(path, projectName, sdkPath string) bool {
    args := []string{}
    args = append(args, "apkbuilder")
    args = append(args, path+"/bin/"+projectName+"-unsigner.apk")
    args = append(args, "-v")
    args = append(args, "-u")
    args = append(args, "-z")
    args = append(args, path+"/bin/resources.ap_")
    args = append(args, "-f")
    args = append(args, path+"/bin/classes.dex")
    args = append(args, "-rf")
    args = append(args, path+"/src")

    if result, rErr := RunCommand(args); result {
        fmt.Println("Built unsigner apk success.")
        return true
    } else {
        if strings.Contains(rErr.Error(), "executable file not found in") || strings.Contains(rErr.Error(), "no such file or directory") {
            fmt.Println("apkbuilder not found")

            createApkBuilderCmd := "cd " + sdkPath + "/tools && cat android | sed -e 's/com.android.sdkmanager.Main/com.android.sdklib.build.ApkBuilderMain/g' > apkbuilder"
            if !RunShell(createApkBuilderCmd) {
                fmt.Println("apkbuilder command create faild")
                return false
            } else {
                chmodCmd := "cd " + sdkPath + "/tools && chmod a+x apkbuilder"
                if !RunShell(chmodCmd) {
                    fmt.Println("apkbuilder command chmod faild")
                    return false
                } else {
                    RunUnsignAPK(path, projectName, sdkPath)
                }
            }
        } else {
            fmt.Println("found apkbuilder, but built unsigner apk failed.")
            return false
        }
    }
    return false
}

/**
-keystore及后面的路径指明密钥文件的位置，
-storepass是用于密钥库完整性的口令，-keypass是专用密钥的口令
-signedjar及后面的路径指明签完名的apk文件的路径，紧接着的是需要签名的apk的路径，最后面是密钥的别名

debug.keystore的name和passwords信息是在SDK文档中找到的，具体内容如下：
The SDK tools create the debug keystore/key with predetermined names/passwords:
Keystore name: "debug.keystore" Keystore password: "android" Key alias: "androiddebugkey" Key password: "android" CN: "CN=Android Debug,O=Android,C=US"
*/
// sign apk
func RunSignAPK(apk *Apk, signType string, ch chan<- int) {
    //    command := "jarsigner -keystore " + "/.android/debug.keystore" +
    //        " -storepass " + "android" +
    //        " -keypass " + "android" +
    //        " -signedjar " + path + "/bin/HelloSign.apk" +
    //        " " +
    //        path + "/bin/HelloAndroid3.apk" +
    //        " " + "androiddebugkey"
    //jarsigner -verbose -keystore liufeng.keystore -signedjar notepad_signed.apk notepad.apk liufeng.keystore

    if signType == "DEBUG" || signType == "debug" {

    } else if signType == "RELEASE" || signType == "release" {

    }

    args := []string{}
    args = append(args, "jarsigner")
    args = append(args, "-verbose")
    args = append(args, "-sigalg")
    args = append(args, "MD5withRSA")
    args = append(args, "-digestalg")
    args = append(args, "SHA1")
    args = append(args, "-keystore")
    args = append(args, apk.Keystore)
    args = append(args, "-storepass")
    args = append(args, apk.Kstorepass)
    args = append(args, "-keypass")
    args = append(args, apk.Keypass)
    args = append(args, "-signedjar")
    args = append(args, apk.Path+"/bin/"+apk.ProjectName+"-release.apk")
    args = append(args, apk.Path+"/bin/"+apk.ProjectName+"-unsigner.apk")
    args = append(args, apk.Keyalias)

    if result, _ := RunCommand(args); result {
        ch <- 1
    } else {
        ch <- 0
    }
}

// sign debug apk
func RunDebugSignAPK(apk *Apk) bool {
    //Keystore名字：“debug.keysotre”
    //Keystore密码：“android”
    //Key别名：“androiddebugkey”
    //Key密码：“android”
    //CN：“CN=Android Debug,O=Android,C=US”

    keystoreFile := "~/.android/debug.keystore"
    keystorePwd := "android"
    keyPwd := "android"
    keyAlias := "androiddebugkey"

    debugSignApkPath := apk.Path+"/bin/"+apk.ProjectName+"-debugsign.apk"

    command := "jarsigner -verbose -sigalg MD5withRSA -digestalg SHA1 -keystore "+keystoreFile+" -storepass "+keystorePwd+" -keypass "+keyPwd+" -signedjar "+debugSignApkPath +" "+apk.Path+"/bin/"+apk.ProjectName+"-unsigner.apk" +" "+keyAlias

    if RunShell(command) {
        fmt.Println("Built sign debug apk success.")
        result := RunAlignAPK(debugSignApkPath, apk.Path+"/bin/"+apk.ProjectName+"-debugsign-aligned.apk")
        return result
    } else {
        fmt.Println("Built sign debug apk failed.")
        return false
    }
}

//对签名后的.apk文件进行对齐处理（不进行对齐处理是不能发布到Google Market的）
// use zipalign tool to align apk
func RunAlignAPK(signApkFullPath, outName string) bool {
    // TODO zipalign工具
    // zipalign -v 4 notepad_signed.apk notepad_signed_aligned.apk
    args := []string{}
    args = append(args, "zipalign")
    args = append(args, "-v")
    args = append(args, "4")
    args = append(args, signApkFullPath)
    args = append(args, outName)

    if result, _ := RunCommand(args); result {
        fmt.Println("align apk success.")
        return true
    } else {
        fmt.Println("align apk failed.")
        return false
    }
}

func RunBuild(apk *Apk, sdkpath string) bool {
    // clear bin & gen
    genPath := apk.Path + "/gen"
    if CheckPath(genPath) {
        os.RemoveAll(genPath)
    }
    binPath := apk.Path + "/bin"
    if CheckPath(binPath) {
        os.RemoveAll(binPath)
    }

    if !RunRClass(apk, sdkpath) {
        return false
    }

    if !RunAIDL(apk, sdkpath) {
        return false
    }

    if !RunClass(apk, sdkpath) {
        return false
    }

    if !RunDEX(apk.Path) {
        return false
    }

    if !RunResource(apk.Path, apk.Api, sdkpath) {
        return false
    }

    if !RunUnsignAPK(apk.Path, apk.ProjectName, sdkpath) {
        return false
    }

    if !RunDebugSignAPK(apk) {
        return false
    }

    return true
}

func RuntineBuild(apk *Apk, sdkpath string, ch chan int) {
    if RunBuild(apk, sdkpath) {
        ch <- 1
        return
    }

    ch <- 0
}
