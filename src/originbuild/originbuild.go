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

// R.class
func RunResultRClass(path, apiLib string) {
    command := "aapt package -f -m" \
        + " -M " + "\HelloAndroid3\AndroidManifest.xml" \
        + " -J " + "\HelloAndroid3\gen" \
        + " -S " + "\HelloAndroid3\res" \
        + " -I " + "platforms\android-10\android.jar"

    rCmd := exec.Command(command)
    rResult, rErr := rCmd.Output()
    if rErr != nil {
        fmt.Fprintf(os.Stderr, "The command failed to perform: %s\n", rErr)
        return false
    }
    if DEBUG {
        fmt.Fprintf(os.Stdout, "----update Result: %s \n", rResult)
    }
}

// AIDI
func RunResultAIDL(path string) {
    // TODO nothing
}

// classes
func RunResultClass(path string) {
    command := "javac -encoding GB18030 -target 1.6" \
        + " -bootclasspath " + "\platforms\android-10\android.jar" \
        + " -d " + "\HelloAndroid3\bin" + " " \
        + "\HelloAndroid3\src\com\******\HelloAndroid3\HelloAndroid3.java" + " " \
        + "\HelloAndroid3\gen\com\******\HelloAndroid3\R.java"

    rCmd := exec.Command(command)
    rResult, rErr := rCmd.Output()
    if rErr != nil {
        fmt.Fprintf(os.Stderr, "The command failed to perform: %s\n", rErr)
        return false
    }
    if DEBUG {
        fmt.Fprintf(os.Stdout, "----update Result: %s \n", rResult)
    }
}

// DEX
func RunResultDEX(path string) {
    command := "dx --dex" \
        + " --output=" + "\HelloAndroid3\bin\classes.dex" + " " \
        + "c:\docume~1\******\workspace\HelloAndroid3\bin\"

    rCmd := exec.Command(command)
    rResult, rErr := rCmd.Output()
    if rErr != nil {
        fmt.Fprintf(os.Stderr, "The command failed to perform: %s\n", rErr)
        return false
    }
    if DEBUG {
        fmt.Fprintf(os.Stdout, "----update Result: %s \n", rResult)
    }
}

// Resource
func RunResultR(path string) {
    command := "aapt package -f" \
        + " -M " + "\HelloAndroid3\AndroidManifest.xml" \
        + " -S " + "\HelloAndroid3\res" \
        + " -A " + "\HelloAndroid3\assets" \
        + " -I " + "platforms\android-10\android.jar" \
        + " -F " + "\HelloAndroid3\bin\resources.ap_"

    rCmd := exec.Command(command)
    rResult, rErr := rCmd.Output()
    if rErr != nil {
        fmt.Fprintf(os.Stderr, "The command failed to perform: %s\n", rErr)
        return false
    }
    if DEBUG {
        fmt.Fprintf(os.Stdout, "----update Result: %s \n", rResult)
    }
}

// unsign apk
func RunResultUnsignAPK(path string) {
    command := "apkbuilder " + "\HelloAndroid3\bin\HelloAndroid3.apk" \
        + " -v -u -z " + "\HelloAndroid3\bin\resources.ap_" \
        + " -f " + "\HelloAndroid3\bin\classes.dex" \
        + " -rf " + "\HelloAndroid3\src"

    rCmd := exec.Command(command)
    rResult, rErr := rCmd.Output()
    if rErr != nil {
        fmt.Fprintf(os.Stderr, "The command failed to perform: %s\n", rErr)
        return false
    }
    if DEBUG {
        fmt.Fprintf(os.Stdout, "----update Result: %s \n", rResult)
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
func RunResultSignAPK(path string) {
    command := "jarsigner" \
        + " -keystore " + "\.android\debug.keystore" \
        + " -storepass " + "android" \
        + " -keypass " + "android" \
        + " -signedjar " + "\HelloAndroid3\bin\HelloSign.apk" \
        + " " \
        + "\HelloAndroid3\bin\HelloAndroid3.apk" \
        + " " + androiddebugkey

    rCmd := exec.Command(command)
    rResult, rErr := rCmd.Output()
    if rErr != nil {
        fmt.Fprintf(os.Stderr, "The command failed to perform: %s\n", rErr)
        return false
    }
    if DEBUG {
        fmt.Fprintf(os.Stdout, "----update Result: %s \n", rResult)
    }
}

//对签名后的.apk文件进行对齐处理（不进行对齐处理是不能发布到Google Market的）
// use zipalign tool to align apk
func RunResultAlignAPK(path string) {
    // TODO zipalign工具
}
