# Gob Utility 使用示例

这个示例演示了增强版 gob 工具包的所有功能。

## 功能特性

1. **基本读写功能**:
    - `WriteToGobFile`: 将数据写入gob文件
    - `ReadFromGobFile`: 从gob文件读取数据

2. **字节切片编解码**:
    - `EncodeToBytes`: 将数据编码为gob格式的字节切片
    - `DecodeFromBytes`: 将gob格式的字节切片解码为数据

3. **增强功能**:
    - `AppendToGobFile`: 追加数据到gob文件
    - `WriteToGobFileWithBackup`: 带备份的写入操作
    - `ReadLastFromGobFile`: 读取文件中的最后一条记录
    - `GetGobFileSize`: 获取gob文件大小

## 运行示例

在项目根目录下执行以下命令:

```bash
go run examples/gob_example.go
```

## 输出说明

示例程序会展示以下操作:

1. 基本的gob文件读写
2. 编码为字节切片和解码
3. 追加数据到现有文件
4. 带备份的写入操作
5. 获取文件大小信息

所有操作都会成功执行，并在控制台显示结果。

## 清理

示例运行完成后，会自动清理创建的临时文件。

# 加密解密功能示例

这个示例演示了加密工具包的所有功能。

## 功能特性

1. **对称加密 (AES-GCM)**:
    - `AESGCMEncrypt`: 使用AES-GCM算法加密数据
    - `AESGCMDecrypt`: 使用AES-GCM算法解密数据

2. **非对称加密 (RSA)**:
    - `GenerateRSAKeyPair`: 生成RSA密钥对
    - `RSAEncrypt`: 使用RSA公钥加密数据
    - `RSADecrypt`: 使用RSA私钥解密数据

3. **密钥编码/解码**:
    - `EncodePrivateKeyToPEM`: 将RSA私钥编码为PEM格式
    - `EncodePublicKeyToPEM`: 将RSA公钥编码为PEM格式
    - `DecodePrivateKeyFromPEM`: 从PEM格式解码RSA私钥
    - `DecodePublicKeyFromPEM`: 从PEM格式解码RSA公钥

## 运行示例

在项目根目录下执行以下命令:

```bash
go run examples/crypto_example.go
```

## 输出说明

示例程序会展示以下操作:

1. AES-GCM加密和解密
2. RSA密钥生成和加解密
3. RSA密钥的PEM编码和解码

所有操作都会成功执行，并在控制台显示结果。
