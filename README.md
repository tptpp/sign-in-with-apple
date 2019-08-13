# 概述
本文是对AppleID登录接入的相关总结，希望对其他人能有帮助。

苹果在其WWDC19大会上提出了"Sign In with Apple"的概念，类似于微信一键登录，但也有些区别：
* 微信存在一个UnionID、OpenID的概念，苹果只有一个AppleID；
* 微信直接通过api就能拿到用户信息，而苹果拿到的是一个jwt，需要进行加解密。

# 相关参数
根据官方文档(链接见文末)，AppleID登录遵循OAuth2.0协议，主要分为两步：1. 用户授权后获取code；2. 通过code换取token。以下流程在网页端和App端存在一定差异，主要是client_id和redirect_url不同，这里以网页端为例进行说明。

在发起流程前，你需要准备以下参数：

* Team ID，10个字节的字符串，可以在苹果账户后台中看到，位于右上角

* Key Id，10个字节的字符串，可以在苹果账户后台中看到

* Client ID，这里要注意与code授权的平台保持一致。注意网页端与app端的差异，参见遇到的问题一

* Private Key，一个.p8文件，只能从苹果官网下载一次

* Redirect Url，code授权的回调url，app端可以不填

除Redirect Url外，其他几个参数主要用于生成client_secret。根据官方文档，client_secret是如下jwt采用ES256加密的结果：

```
{
    "alg": "ES256",                         // jwt加密算法，固定值
    "kid": "ABC123DEFG"                     // Key Id
}
{
    "iss": "DEF123GHIJ",                    // Team ID
    "iat": 1437179036,                      // jwt生成时间，精确到秒
    "exp": 1493298100,                      // jwt过期时间，unix时间戳，精确到秒
    "aud": "https://appleid.apple.com",     // 授权给此域名，固定值
    "sub": "com.mytest.app"                 // Client ID
}
```

# golang实现
参见[代码](https://github.com/tptpp/sign-in-with-apple/blob/master/main.go)

# 遇到的问题

## 问题一 网页端/App端Client ID不同
* 现象： {"error":"invalid_grant"}
* 原因：client_id填写错误

网页端流程走通后，在调试App端的过程中，总是报错invalid_grant，后来发现是client_id填错了。网页授权登录填写的是Services Id，App端登录需要的是AppId，参见[链接](https://forums.developer.apple.com/thread/118135)
# 参考

* [官方文档](https://developer.apple.com/sign-in-with-apple/get-started/)
* [OAuth 2.0](https://oauth.net/2/)，主要就是code、token、refresh_token的关系
* [jwt](https://jwt.io/introduction/)
* [jwt-go](https://github.com/dgrijalva/jwt-go)
* [参考一](https://medium.com/identity-beyond-borders/how-to-configure-sign-in-with-apple-77c61e336003)
* [参考二](https://developer.okta.com/blog/2019/06/04/what-the-heck-is-sign-in-with-apple#how-sign-in-with-apple-works-hint-it-uses-oauth-and-oidc)
