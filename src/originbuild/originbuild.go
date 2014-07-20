package originbuild

import (
    "fmt"
    "os"
    "os/exec"
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
    IsRealLibrary bool

    Keystore   string
    Keyalias   string
    Kstorepass string
    Keypass    string
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

// R.class
func RunRClass(apk *Apk, sdkPath string, ch chan int) {
    genPath := apk.Path + "/gen"
    if !CheckPath(genPath) {
        os.Mkdir(genPath, os.FileMode(0777))
    }

    rCmd := exec.Command("aapt", "package", "-f", "-m", "-M", apk.Path+"/AndroidManifest.xml", "-J", genPath, "-S", apk.Path+"/res", "-I", sdkPath+"/platforms/"+apk.Api+"/android.jar")
    rResult, rErr := rCmd.Output()

    if DEBUG {
        fmt.Fprintf(os.Stdout, "----Command: %s\n", rCmd.Args)
        fmt.Fprintf(os.Stdout, "----run Result: %s \n", rResult)
    }

    if rErr == nil {
        ch <- 1
    } else {
        fmt.Fprintf(os.Stderr, "The command failed to perform: %s\n", rErr)
        ch <- 0
    }
}

// AIDI
func RunAIDL(path string) {
    // TODO nothing
}

// classes
func RunClass(apk *Apk, sdkPath string, ch chan int) {
    //    command := "javac -encoding GB18030 -target 1.6" +
    //        " -bootclasspath " + "/platforms/android-10/android.jar" +
    //        " -d " + path + "/bin" + " " +
    //        "/HelloAndroid3/src/com/******/HelloAndroid3/HelloAndroid3.java" + " " +
    //        "/HelloAndroid3/gen/com/******/HelloAndroid3/R.java"

    rCmd := exec.Command("javac", "-encoding", "GB18030", "-target", "1.6", "-bootclasspath", sdkPath+"/platforms/"+apk.Api+"/android.jar", "-d", apk.Path+"/bin")
    rResult, rErr := rCmd.Output()

    if DEBUG {
        fmt.Fprintf(os.Stdout, "----Command: %s\n", rCmd.Args)
        fmt.Fprintf(os.Stdout, "----run Result: %s \n", rResult)
    }

    if rErr == nil {
        ch <- 1
    } else {
        fmt.Fprintf(os.Stderr, "The command failed to perform: %s\n", rErr)
        ch <- 0
    }
}

// DEX
func RunDEX(path string, ch chan int) {
    //    command := "dx --dex" +
    //        " --output=" + path + "/bin/classes.dex" + " " +
    //        path + "/bin/"

    rCmd := exec.Command("dx", "--dex", "--output="+path+"/bin/classes.dex", path+"/bin/")
    rResult, rErr := rCmd.Output()

    if DEBUG {
        fmt.Fprintf(os.Stdout, "----Command: %s\n", rCmd.Args)
        fmt.Fprintf(os.Stdout, "----run Result: %s \n", rResult)
    }

    if rErr == nil {
        ch <- 1
    } else {
        fmt.Fprintf(os.Stderr, "The command failed to perform: %s\n", rErr)
        ch <- 0
    }
}

// Resource
func RunResource(path, api, sdkPath string, ch chan int) {
    //    command := "aapt package -f" +
    //        " -M " + path + "/AndroidManifest.xml" +
    //        " -S " + path + "/res" +
    //        " -A " + path + "/assets" +
    //        " -I " + "platforms/android-10/android.jar" +
    //        " -F " + path + "/bin/resources.ap_"

    rCmd := exec.Command("aapt", "package", "-f", "-M", path+"/AndroidManifest.xml", "-S", path+"/res", "-A", path+"/assets", "-I", sdkPath+"/platforms/"+api+"/android.jar", "-F", path+"/bin/resources.ap_")
    rResult, rErr := rCmd.Output()

    if DEBUG {
        fmt.Fprintf(os.Stdout, "----Command: %s\n", rCmd.Args)
        fmt.Fprintf(os.Stdout, "----run Result: %s \n", rResult)
    }

    if rErr == nil {
        ch <- 1
    } else {
        fmt.Fprintf(os.Stderr, "The command failed to perform: %s\n", rErr)
        ch <- 0
    }
}

// unsign apk
func RunUnsignAPK(path, projectName string, ch chan int) {
    //    command := "apkbuilder " + path + "/bin/HelloAndroid3.apk" +
    //        " -v -u -z " + path + "/bin/resources.ap_" +
    //        " -f " + path + "/bin/classes.dex" +
    //        " -rf " + path + "/src"

    rCmd := exec.Command("apkbuilder", path+"/bin/"+projectName+".apk", "-v", "-u", "-z", path+"/bin/resources.ap_", "-f", path+"/bin/classes.dex", "-rf", path+"/src")
    rResult, rErr := rCmd.Output()

    if DEBUG {
        fmt.Fprintf(os.Stdout, "----Command: %s\n", rCmd.Args)
        fmt.Fprintf(os.Stdout, "----run Result: %s \n", rResult)
    }

    if rErr == nil {
        ch <- 1
    } else {
        fmt.Fprintf(os.Stderr, "The command failed to perform: %s\n", rErr)
        ch <- 0
    }
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
func RunSignAPK(apk *Apk, ch chan int) {
    //    command := "jarsigner -keystore " + "/.android/debug.keystore" +
    //        " -storepass " + "android" +
    //        " -keypass " + "android" +
    //        " -signedjar " + path + "/bin/HelloSign.apk" +
    //        " " +
    //        path + "/bin/HelloAndroid3.apk" +
    //        " " + "androiddebugkey"

    rCmd := exec.Command("jarsigner", "-keystore", apk.Keystore, "-storepass", apk.Kstorepass, "-keypass", apk.Keypass, "-signedjar", apk.Path+"/bin/"+apk.ProjectName+"-release.apk", apk.Path+"/bin/"+apk.ProjectName+".apk", apk.Keyalias)
    rResult, rErr := rCmd.Output()

    if DEBUG {
        fmt.Fprintf(os.Stdout, "----Command: %s\n", rCmd.Args)
        fmt.Fprintf(os.Stdout, "----run Result: %s \n", rResult)
    }

    if rErr == nil {
        ch <- 1
    } else {
        fmt.Fprintf(os.Stderr, "The command failed to perform: %s\n", rErr)
        ch <- 0
    }
}

//对签名后的.apk文件进行对齐处理（不进行对齐处理是不能发布到Google Market的）
// use zipalign tool to align apk
func RunAlignAPK(path string) {
    // TODO zipalign工具
}
