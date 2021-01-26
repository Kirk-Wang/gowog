# 试玩 GOWOG ，初探 OpenAI（使用 NeuroEvolution 神经进化）与 Golang 多人在线游戏开发

![screenshot](document/images/screenshot.gif)

GOWOG：
* 原项目：https://github.com/giongto35/gowog
* 我调整过的：https://github.com/Kirk-Wang/gowog

GOWOG 是一款迷你的，使用 Golang 编写的多人 Web 游戏。

## 试玩游戏

Demo：http://game.giongto35.com

## 在 Agent 上的 AI 实验

由于服务器，客户端和消息是分离的，因此很容易与后端进行通信。
此项目是用 `Python` 编写的 `AI agent`，可以学习与环境的交互。
这个实验是利用 `neuroevolution` (神经进化）在迷宫中寻找一条路径。

油管 Demo：https://www.youtube.com/watch?v=pWbb1m91mhU


## 本地 Docker 运行

`run_local.sh`

```sh
#!/bin/bash
docker build . --build-arg HOSTNAME=localhost:8080 -t gowog_local|| exit
docker stop gowog
docker rm gowog
docker run --privileged -d --name gowog -p 8080:8080 gowog_local server -prod client/build/
```

执行，我本地是 `Mac`：

```sh
./run_local.sh
```

打开 [http://localhost:8080](http://localhost:8080)

## 本地开发

游戏包含两部分：服务器和客户端。服务器使用 `Golang`，客户端使用 `Node.JS` 和 `Phaser` 游戏引擎。

### 服务器

为少调整过的项目：[Kirk-Wang/gowog](https://github.com/Kirk-Wang/gowog)

我的本地环境是 `go 1.14`。

本地启动：

```sh
go run cmd/server/*
```

服务器将监听 `8080`。

### 客户端

```sh
npm install
npm run dev -- --env.HOST_IP=localhost:8080 # HOST_IP -> 服务器地址
```

进入 [http://localhost:3000](http://localhost:3000)

### 注意

在开发过程中，客户端在端口 `3000` 上运行，服务器端在端口 `8080` 上运行。

在 `production` 和 `docker`环境中，已构建 `Client`，golang 服务器将在同一端口 `8080` 上返回客户端页面。因此，如果我们运行 `docker` 环境，则游戏将在浏览器中的 [http://localhost:8080](http://localhost:8080) 运行。

### 通讯约定

服务器和客户端之间的通信包基于 `protobuf`。 安装 `protoc` 以生成 `protobuf`。
* http://google.github.io/proto-lens/installing-protoc.html

每次你在 `server/message.proto` 中有更改（package singature）时。请运行：

```sh
cd server
./generate.sh
```

## 游戏前端设计

![](document/images/phaser-es6-webpack.png)

这个前端项目是基于：
* https://github.com/RenaudROHLINGER/phaser-es6-webpack

```js
├── client
│   ├── index.html
│   ├── src
│   │   ├── config.js: javascript config
│   │   ├── index.html
│   │   ├── main.js
│   │   ├── sprites
│   │   │   ├── Leaderboard.js: Leaderboard object
│   │   │   ├── Map.js: Map object
│   │   │   ├── Player.js: Player object
│   │   │   └── Shoot.js: Shoot object
│   │   ├── states
│   │   │   ├── Boot.js Boot screen
│   │   │   ├── const.js
│   │   │   ├── Game.js: Game master
│   │   │   ├── message_pb.js: Protobuf Message
│   │   │   ├── Splash.js
│   │   │   └── utils.js
│   │   └── utils.js
```

每个对象都是从 `Sprite` 继承的类。
玩家包含 `shootManager`，每次射击时，shoot manager 都会生成新的 bullet。

## 游戏后端设计方案

### Components（组件）

游戏中主要有 `5` 个实体。他们的状态是私有的

| 实体 | 私有状态 |  |
| ------ | ----- | ----------- |
| Client | websocket.Conn | client hold 住 websocket 连接 |
| Hub | Client | Hub 处理所有通讯, 包含所有 client 列表 |
| ObjManager | Player, Shoot, ... | ObjManager 包含所有 Player 和 Shoot，处理游戏逻辑 |
| Game Master | ObjManager, Hub | Master object 由 ObjManager 和 Hub 组成 |

### Architecture（架构图）

![Architecture](document/images/architecture.png)

不同的实体通过包装在函数中的 `channel` 彼此调用。

### Client 与 Server 交互设计方案

**Player connect（玩家连接）**

![PlayerConnect](document/images/playerconnect.png)

**Player Disconnect（玩家断开连接）**

![PlayerDisconnect](document/images/playerdisconnect.png)

**Client input（客户端输入）**

![ClientInput](document/images/playeraction.png)

### Profile

`Profile` 是研究 `Golang` 性能并找出 `slow components` 的方法。运行服务器时，可以使用标志 `--cpuprofile` 和`--memprofile` 来配置 server。

```sh
cd server
go run cmd/server/* --cpuprofile --memprofile
```

### 代码结构

```
├── server
│   ├── cmd
│   │   └── server
│   │       └── server.go: Entrypoint running server
│   ├── game
│   │   ├── common
│   │   ├── config
│   │   │   └── 1.map: Map represented 0 and 1
│   │   ├── gameconst
│   │   ├── game.go: Game master objects, containing logic and communication
│   │   ├── mappkg
│   │   ├── objmanager
│   │   ├── playerpkg
│   │   ├── shape
│   │   ├── shootpkg
│   │   ├── types.go
│   │   └── ws
│   │       ├── wsclient.go
│   │       └── wshub.go
│   ├── generate.sh: Generate protobuf for server + client + AI environment
│   ├── message.proto
│   └── Message_proto
│       └── message.pb.go
├── Dockerfile
└── run_local.sh
```

## AI 训练设计方案

此仓库包含遵循 `openAI Gym` 格式和训练脚本的 `CS2D` 环境。
训练脚本使用 `NeuroEvolution(神经进化)`在迷宫中找到到达目的地的最短路径。

https://www.youtube.com/watch?v=pWbb1m91mhU

### 运行

按照的说明运行 `gowog` 环境。即本地 Docker 运行：

```sh
./run_local.sh
```

使用 `virtualenv` 设置 `python3` 虚拟环境（直接用 Docker 吧~）。
* 安装 `requirements.txt` 所包含的库。

运行训练脚本

```sh
python train_ga.py -n save_file_name
```

`save_fie_name` 是保存权重（`weights`）的地方。
在下一次，如果我们指定了一个现有的文件，它将继续从该文件的最后一次运行中的权重（`weights`）进行训练。

### Genetic Algorithm（遗传算法）

_cs2denv_ga.py 的实现_

基于机器学习的目的，`CS2D Agent` 是在 `CS2D` 上构建的。
它遵循 `openAI gym`，支持 `agent` 的基本方法，包括：`reset()`、`step()`、`observation_space` 和`action_space`。

`ObservationSpace` 是一个一维数组，它由来自服务器的 `update_player` 消息构造而成

1. Player position（玩家位置）, player size（玩家大小）, number of columns（列数）, number of rows（行数）, block width（块宽度）, block height（块高度）
2. 到左，右，上，下到最近 block（块）的距离。此输入是为了避免碰撞
3. 玩家在二进制块地图（binary block map）中的位置。地图是0和1的2维数组（0为空，1为块）
4. binary block map

The Reward is the 1 / distance to the goal. If the agent is close to the goal by 100, the reward is 1 and the episode finishes.
* 奖励是 1 / distance(距离的目标)。如果 agent 接近目标 100 点，那么奖励就是 1，情节结束。

### NeuroEvolution（神经进化）

_train_ga.py 的实现_

神经网络（Neural Network）通过使输入（观察空间）通过神经网络来获得最佳动作。

NeuroEvolution（神经进化）是使用进化算法不断改进人工神经网络的AI。对于每次迭代（生成），程序将基于前一次迭代中的最佳设置生成一组新的神经网络权重。 由先前的 `NN(神经网络)` 生成一个 `NN` 的过程叫做 `Mutate`，它给神经网络中的每个参数添加随机噪声。

一个特别的改进是，我们只存储应用于神经网络的噪声种子列表，而不是存储所有的代权值。因为在同一个种子下，所有的随机化都是相同的，所以一个种子可以代表一个网络的突变算子。我们不需要保留每一代的所有权值，我们只需要存储一组从开始到当前一代的种子，然后从这组种子中重新构造权值来得到所有神经网络的权值。

代码是基于 Maxim Lapan 的 "Deep Reinforcement Learning Hands-On"

* https://github.com/PacktPublishing/Deep-Reinforcement-Learning-Hands-On/blob/master/Chapter16/04_cheetah_ga.py













