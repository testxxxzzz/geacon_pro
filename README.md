# geacon_pro

## 项目介绍
本项目基于[geacon](https://github.com/darkr4y/geacon)项目对cobaltstrike的beacon进行了重构，并适配了大部分Beacon的功能。
该项目仅用于对CobaltStrike协议的学习测试。请勿使用于任何非法用途，由此产生的后果自行承担。

本项目与好兄弟Z3ratu1共同开发，他实现了一版支持4.0版本的[geacon_plus](https://github.com/Z3ratu1/geacon_plus)，我这边实现了一版支持4.1及以上版本的beacon，大致功能类似，有部分功能不同。

暂时只测试了4.3版本，理论上来说4.1+版本均支持，如果有不支持的版本请及时通知我。

目前实现的功能具备免杀性，可过Defender、360核晶（除powershell）、卡巴斯基（除内存操作外，如注入原生cs的dll）、火绒
上述测试环境均为实体机
若想使用免杀powershell和免杀bypassUAC的话，请参考我的另外两个小工具，暂未进行集成。

开发的过程中参考了鸡哥的数篇文章以及许许多多的项目，同时抓包对服务端返回的内容进行猜测，并对服务端java代码进行了部分的理解。

欢迎师傅们交流，欢迎指出问题。


## 使用方法
本项目支持windows、linux、mac平台的使用。

基础的使用方法可参考原项目，windows编译时添加-ldflags "-H windowsgui -s -w"减小程序体积并取消黑框。linux和mac编译的时候添加-ldflags "-s -w"减小程序体积，然后后台运行。

目前项目有部分控制台输出内容，若想删除可在代码中删除。

## 实现功能
### windows平台支持的功能：
sleep、shell、upload、download、exit、cd、pwd、file_browse、ps、kill、getuid、mkdir、rm、cp、mv、run、execute、drives、powershell-import、powershell、execute-assembly（不落地执行c#）、多种线程注入的方法（可自己更换源码）、shinject、dllinject、管道的传输、多种cs原生反射型dll注入（mimikatz、portscan、screenshot、keylogger等）、令牌的窃取与还原、令牌的制作、代理发包等功能

### linux和mac平台支持的功能：
sleep、shell、upload、download、exit、cd、pwd、file_browse、ps、kill、getuid、mkdir、rm、cp、mv

文件管理部分支持图形化交互

### C2profile：
适配了C2profile流量侧的设置与部分主机侧的设置，支持的算法有base64、base64url、mask、netbios、netbiosu、详情见config.go，这里给出示例C2profile，修改完C2profile后请不要忘记在config.go中对相应位置进行修改：
```
set sleeptime "3000";

https-certificate {
    set C "KZ";
    set CN "foren.zik";
    set O "NN Fern Sub";
    set OU "NN Fern";
    set ST "KZ";
    set validity "365";
}

http-get {

	set uri "/www/handle/doc";

	client {
		metadata {
			base64url;
			prepend "SESSIONID=";
			header "Cookie";
		}
	}

	server {
		header "Server" "nginx/1.10.3 (Ubuntu)";
    		header "Content-Type" "application/octet-stream";
        	header "Connection" "keep-alive";
        	header "Vary" "Accept";
        	header "Pragma" "public";
        	header "Expires" "0";
        	header "Cache-Control" "must-revalidate, post-check=0, pre-check=0";

		output {
			mask;
			netbios;
			prepend "data=";
			append "%%";
			print;
		}
	}
}

http-post {
	set uri "/IMXo";
	client {
		
		id {				
			mask;
			netbiosu;
			parameter "doc";
		}

		output {
			mask;
			base64url;
			prepend "data=";
			append "%%";		
			print;
		}
	}

	server {
		header "Server" "nginx/1.10.3 (Ubuntu)";
    		header "Content-Type" "application/octet-stream";
        	header "Connection" "keep-alive";
       	 	header "Vary" "Accept";
        	header "Pragma" "public";
        	header "Expires" "0";
        	header "Cache-Control" "must-revalidate, post-check=0, pre-check=0";
          
		output {
			mask;
			netbios;
			prepend "data=";
			append "%%";
			print;
		}
	}
}

post-ex {
    set spawnto_x86 "c:\\windows\\syswow64\\rundll32.exe";
    set spawnto_x64 "c:\\windows\\system32\\rundll32.exe";
    
    set thread_hint "ntdll.dll!RtlUserThreadStart+0x1000";
    set pipename "DserNamePipe##, PGMessagePipe##, MsFteWds##";
    set keylogger "SetWindowsHookEx";
}
```

### 目前需要改进的地方：
* 堆内存加密目前不稳定，暂未正式使用
* ps暂不支持图形化交互，支持命令行显示
* 部分功能暂未支持x86系统（最近太忙了，会尽快改出来）

### 主体代码结构
#### config
* 公钥、C2服务器地址、https通信、超时的时间、代理等设置
* C2profile设置
#### crypt
* 通信需要的AES、RSA加密算法
* C2profile中加密算法的实现
#### packet
* commands为各个平台下部分功能的实现
* execute_assembly为windows平台下内存执行不落地c#的代码
* heap为windows平台下堆内存加密代码
* http为发包的代码
* inject为windows平台下进程注入的代码
* jobs为windows平台下注入cs原生反射型dll并管道回传的代码
* packet为通信所需的部分功能
* token为windows平台下令牌相关的功能
#### services
对packet里面的功能进行了跨平台封装，方便main.go调用
#### sysinfo
* meta为元信息的处理
* sysinfo为不同平台下有关进程与系统的判断及处理
#### main.go
主方法对各个命令进行了解析与执行，以及对结果和错误进行了返回

## 部分功能的实现细节
* shell直接调用了golang的os/exec库

* 由于os/exec库是新起进程执行的命令，无法以窃取来的令牌的身份执行命令，因此run和execute的实现在没有窃取令牌的时候调用了CreateProcess，窃取令牌后调用CreateProcessWithTokenW以令牌权限来执行命令，run会通过管道回传执行的结果。
因此要注意，如果想以令牌的权限执行命令，那么需要用run或execute而不是shell。

* powershell-import部分与cs的思路一样，先把输入的部分powershell module保存，之后在执行powershell的时候本地开一个端口并把module放上去，powershell直接请求该端口进行不落地的powershell module加载。

* powershell命令直接调用了powershell，会被360监控，可以尝试用免杀的方式执行。

* execute-assembly的实现与cs原生的实现不太一样，cs的beacon从服务端返回的内容的主体部分是c#的程序以及开.net环境的dll。cs的beacon首先拉起来一个进程（默认是rundll32），之后把用来开环境的dll注入到该进程中，然后将c#的程序注入到该进程并执行。考虑到步骤过于繁琐，并且容易拿不到执行的结果，我这里直接用[该项目](https://github.com/timwhitez/Doge-CLRLoad)实现了execute-assembly的功能，但未对全版本windows进行测试。

* 进程注入shinject和dllinject采用的是APC注入。

* cs原生反射型dll注入的思路是先拉起来一个rundll32进程，之后把dll注进去执行，但是会被360核晶报远程线程注入。尝试使用了unhook等方法均失败，最后发现了将dll注入自己是不会被查杀的，因此对服务端下发的dll进行了简单的修改，然后将其注入到木马本身。dll注入通过管道将结果异步地回传给服务端。因此目前的dll反射注入采用了注入自己的方法，后续用户可通过配置文件进行注入方式的更改。

* 令牌的部分目前实现了令牌的窃取、还原、制作。

* 考虑到渗透中常常存在着内网主机上线的情况，即边缘主机出网，内网主机不出网的情况。目前实现的木马暂不支持代理转发的功能，但是可以通过设置config.go中的proxy参数，通过边缘主机的代理进行木马的上线。即如果在边缘主机的8080端口开了个http代理，那么在config.go中设置ProxyOn为true，Proxy为`http://ip:8080`即可令内网的木马上线我们的C2服务器。

* 堆内存加密的方法实现参考了[该文章](https://cloud.tencent.com/developer/article/1949555)。在sleep之前先将除主线程之外的线程挂起，之后遍历堆对堆内存进行加密。sleep结束后解密并将线程恢复。不过该功能较为不稳定，有时在进行堆遍历的时候会突然卡住或者直接退出，并且考虑到后台可能会有keylogger或portscan这种的持久任务，将线程全部挂起有些不合适，如果有师傅有好的想法欢迎来讨论。同时我不太理解为什么go的time.Sleep函数在其他线程都挂起之后调用会一直沉睡，而调用windows.SleepEx就不会有问题，还望师傅们解答。
