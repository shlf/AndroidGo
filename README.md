AndroidGo
=========

Auto build android project use ant command.

Android工程的自动化编译脚本工具------androidGo
  此工具采用最新的计算机GO语言(golang)完成，能在linux系统上顺利工作，支持32bits和64bits。

场景：
  如project A 依赖 project B（而B也依赖D），project B 依赖 project C，project C 依赖 project D。
  那么你只需要告知此工具project A的绝对路径，工具即可自动查找依赖，并依次编译。
