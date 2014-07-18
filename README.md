#AndroidGo#

Auto build android project use ant command.

Android工程的自动化编译脚本工具------androidGo<br>
此工具采用最新的计算机GO语言(golang)完成，能在linux系统上顺利工作，支持32bits和64bits。

##场景：##
如project A 依赖 project B（而B也依赖D），project B 依赖 project C，project C 依赖 project D。
那么你只需要告知此工具project A的绝对路径，工具即可自动查找依赖，并依次编译。

##帮助：##
>     Usage :
       -h : Show command help tip.
       -c : Check your environment.
       -v : Show command version number.
       -p : Input android project path.
       -publish : Choice release(default) or debug mode to publish.
       -tool    : Update project use Git(default) or other code store manager tools.
       -ppath   : Update project path.
 -tool和-ppath一起使用，用于从仓库中更新代码;<br>
 -tool可以指定Git(暂时只支持Git)、svn、hg等<br>
 -ppath用于指定仓库目录
 
 测试：<br>
------------
 1.Master分支:顺序执行编译 
 >Tool name : androidGo-sip
 
 2.Develop-goruntime分支:具有多核并行处理能力 
 >Tool name : androidGo-mult
 
 现对比如下（编译相同的Android工程）:
 <table>
    <tr>
        <td>Tool</td><td>CPU Number</td><td>Dispend Time(Unit:ns)</td>
    </tr>
    <tr>
        <td>androidGo-sip</td><td>1</td><td>175881123215</td>
    </tr>
    <tr>
        <td>androidGo-sip</td><td>4</td><td>98466618185</td>
    </tr>
    <tr>
        <td>androidGo-mult</td><td>1</td><td>113663692767</td>
    </tr>
    <tr>
        <td>androidGo-mult</td><td>4</td><td>65844843454</td>
    </tr>
</table>
从表格可以看出，在机器配置相同、工程一样的情况下，具有多核并行能力的编译工具更加节省时间<br>
尤其在多核环境下的表现能力，尤为突出。<br>
我的工程较大，但是在这种情况下，优化后在4核上消耗近1分钟，没优化单核上消耗接近3分钟，时间差不多压缩近1/3

