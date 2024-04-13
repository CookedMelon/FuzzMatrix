# FuzzMatrix

<!-- PROJECT SHIELDS -->

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![License][license-shield]][license-url]

<!-- PROJECT LOGO -->

## 目录

- [FuzzMatrix](#fuzzmatrix)
  - [目录](#目录)
    - [上手指南](#上手指南)
          - [开发前的配置要求](#开发前的配置要求)
          - [**安装步骤**](#安装步骤)
    - [文件目录说明](#文件目录说明)
    - [部署](#部署)
    - [使用到的框架](#使用到的框架)
    - [贡献者](#贡献者)
      - [如何参与开源项目](#如何参与开源项目)
    - [鸣谢](#鸣谢)

### 上手指南

###### 开发前的配置要求

1. clang-10
2. npm >= 10.2.0
3. python3

###### **安装步骤**

1. Get a free API Key at [https://example.com](https://example.com)
2. Clone the repo

```sh
git clone https://github.com/CookedMelon/FuzzMatrix.git
```

### 文件目录说明

eg:

```bash
filetree
├─README.md
├─demo    #测试程序和数据
│  ├─json
│  └─re2
└─visfuzz
    ├─fuzz    #fuzz框架
    │  ├─afl  #AFL框架
    │  └─llvm #LibFuzzer框架，利用llvm框架构建了函数调用图
    └─visualization #web网页，有几个第三方库（Bootstrap、Chart.js、FontAwesome、jQuery）
        ├─css
        ├─d3
        ├─fig
        ├─js
        ├─scss
        └─vendor
```

### 部署

后端服务

```base
cd server_side
./build.sh
./run.sh
```

前端服务

```bash
cd client_side
npm install
npm run serve
```

### 使用到的框架

|名称|描述|版本要求|
|----|----|------|
|Vue3|前端框架|>= 3.2.13|
|element-plus|前端组件|>=2.6.2|
|D3|控制流图组件|>= 7.9.0|
|AFL|后端Fuzzer|2.52b|
|LLVM|编译器|10|

### 贡献者

请阅读**CONTRIBUTING.md** 查阅为该项目做出贡献的开发者。

#### 如何参与开源项目

贡献使开源社区成为一个学习、激励和创造的绝佳场所。你所作的任何贡献都是**非常感谢**的。

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### 鸣谢

- https://github.com/ChijinZ/VisFuzz/tree/master/
- https://github.com/antonio-morales/Fuzzing101

<!-- links -->

[your-project-path]: CookedMelon/FuzzMatrix
[contributors-shield]: https://img.shields.io/github/contributors/CookedMelon/FuzzMatrix.svg?style=flat-square
[contributors-url]: https://github.com/CookedMelon/FuzzMatrix/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/CookedMelon/FuzzMatrix.svg?style=flat-square
[forks-url]: https://github.com/CookedMelon/FuzzMatrix/network/members
[stars-shield]: https://img.shields.io/github/stars/CookedMelon/FuzzMatrix.svg?style=flat-square
[stars-url]: https://github.com/CookedMelon/FuzzMatrix/stargazers
[issues-shield]: https://img.shields.io/github/issues/CookedMelon/FuzzMatrix.svg?style=flat-square
[issues-url]: https://img.shields.io/github/issues/CookedMelon/FuzzMatrix.svg
[license-shield]: https://img.shields.io/github/license/CookedMelon/FuzzMatrix.svg?style=flat-square
[license-url]: https://github.com/CookedMelon/FuzzMatrix/blob/master/LICENSE
