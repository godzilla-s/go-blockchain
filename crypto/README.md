## 加密算法 

使用go自带的加密包实现区块链中的相关加密函数。使用的是ECC解密与加密。  
Weierstrass椭圆曲线的格式：y² = x³ + ax + b，其中a、b为系数。  
ECDSA： 椭圆曲线签名算法。  


### 测试运行  
测试使用椭圆签名算法： 
```
go run main.go -f crypto
```