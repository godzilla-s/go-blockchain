## 加密算法 

使用go自带的加密包实现区块链中的相关加密函数。使用的是ECC解密与加密。  
Weierstrass椭圆曲线的格式：y² = x³ + ax + b，其中a、b为系数。  
ECDSA： 椭圆曲线签名算法。  

关于椭圆曲线加密，网上描述大多数太难理解，可以看下知乎上有一篇浅显一定的描述： https://www.zhihu.com/question/22399196 

源码部分解析: 
```
type CurveParams struct {
    P       *big.Int // 决定有限域的p的值（必须是素数）
    N       *big.Int // 基点的阶（必须是素数）
    B       *big.Int // 曲线公式的常量（B!=2）
    Gx, Gy  *big.Int // 基点的坐标(即椭圆曲线上的一点，也称为基点)
    BitSize int      // 决定有限域的p的字位数
}
```


```
// 公钥
type PublicKey struct {
	elliptic.Curve
	X, Y *big.Int
}
// 私钥
type PrivateKey struct {
	PublicKey
	D *big.Int
}
``` 
对于计算公式： G * k = K (G,K是公钥(即X，Y)，k是私钥)  

### 第三方加密包 
```
// 基于比特币的机密算法
github.com/btcsuite/btcd/btcec
```

### 测试运行  
测试使用椭圆签名算法： 
```
go run test/main.go
```