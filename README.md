AndroidGo
=========

Auto build android project use ant command.

Android工程的自动化编译脚本工具------androidGo
  此工具采用最新的计算机GO语言(golang)完成，能在linux系统上顺利工作，支持32bits和64bits。

场景：<br>
  如project A 依赖 project B（而B也依赖D），project B 依赖 project C，project C 依赖 project D。
  那么你只需要告知此工具project A的绝对路径，工具即可自动查找依赖，并依次编译。

帮助：<br>
 ------------------------
 Usage :<br>
       -h : Show command help tip.<br>
       -c : Check your environment.<br>
       -v : Show command version number.<br>
       -p : Input android project path.<br>
       -publish : Choice release(default) or debug mode to publish.<br>
       -tool    : Update project use Git(default) or other code store manager tools.<br>
       -ppath   : Update project path.
 ------------------------
 -tool和-ppath一起使用，用于从仓库中更新代码;<br>
 -tool可以指定Git(暂时只支持Git)、svn、hg等<br>
 -ppath用于指定仓库目录
