# mercury

## 微服务端口
bff 8080  
user 8091  
article 8092  
interactive 8093  
comment 8094  
captcha 8095  
sms 8096  
oauth2 8097  
ranking 8098  
crontask 8099

## TODO
1. [article] userrpc improve

## 启动顺序
1. user, sms, interactive, comment, ~~oauth2~~, ~~crontask~~  
2. article, captcha
3. ranking
4. bff