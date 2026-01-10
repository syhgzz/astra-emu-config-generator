# kubenet-config-generator

kubenet-config-generator 是用于生成**大规模多壳层**卫星场景的kubenet配置的工具。输入支持命令行和配置文件（shell-config，YAML）两种方式，输出kubenet配置文件（kubenet-config，YAML）。 

**⚠️ 配置文件并不会上传到远程仓库，请自行在本地妥善保管**

## Directory Structure

```
kubenet-config-generator/
├── .mk/                                # Makefile 片段
│   └── gen.mk                          # 生成相关 Makefile 片段
│   └── test.mk                         # 测试相关 Makefile 片段
├── cmd/                                # Cobra 命令目录
├── config/                             # shell-config 配置文件目录
│   └── multi_shell.yaml.example        # 示例配置文件
├── output/                             # 输出文件目录
│   ├── tle/                            # 生成 TLE 文件目录
│   └── kubenet/                        # 生成 kubenet-config 文件目录
├── template/                           # 模板文件目录
│   └── kubenet_template.yaml.example   # kubenet 配置模板示例文件
├── .gitignore                          # Git 忽略文件
├── Makefile                            # Makefile 主文件
├── README.md                           # 项目说明文件
└── make.env.example                    # Makefile 环境变量示例文件
```

## Usage

1. 复制示例文件
    ```bash
    cp make.env.example make.env
    cp config/multi_shell.yaml.example config/multi_shell.yaml
    cp template/kubenet_template.yaml.example template/kubenet_template.yaml
    ```
2. 编辑 `template/kubenet_template.yaml` 文件，配置 kubenet 模版内容

3. 编辑 `config/multi_shell.yaml` 文件，配置多壳层场景，参数见下文表格

4. 编辑 `make.env` 文件，配置文件名，指定壳层等

5. 运行 Makefile 生成 kubenet 配置文件。本工具支持单壳层和多壳层两种生成模式，多壳层场景下会自动做配置文件合并
    ```bash
    # clean output directory
    make clean-output
    # generate multiple shells
    make gen-all
    # generate one shell
    make gen-one
    ```

## Configuration

本工具输入支持命令行和配置文件（shell-config，YAML）两种方式。
命令行使用参考`.mk/gen.mk`中的命令。如果`-c`参数为空，则从命令行读取壳层参数，否则从`-c`指定配置文件读取。这里介绍两种模式的配置参数。


|     |     |     |     |     |     |     |     |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 参数名 | 必须 | 类型  | 单位  | 命令行参数 | yaml文件配置 | 默认值 | 约束  |
| 配置文件路径 |     | string |     | \-c<br>\--config |     | ""  |     |
| shell index | ✅   | int |     | \-s<br>\--shell-index | shell-index | 1   |     |
| 卫星总数 |     | int |     | \-N<br>\--num | satellite-num | 300 |     |
| 平分球体总轨道面数量 |     | int |     | \-P<br>\--plane | plane-num | 10  | 整除卫星总数 |
| 需要生成的轨道面数量 |     | int |     | \-E<br>\--neede-plane | needed-plane-num | 10  | ≤ plane-num |
| 相位  |     | float64 |     | \-F<br>\--phase | phase | 0   | < plane-num |
|     |     |     |     |     |     |     |     |
| tle文件输出路径 | ✅   | string |     | \-o<br>\--output |     | -  |     |
| 轨道高度 |     | float64 | km  | \-H<br>\--height | height | 500 |     |
| 第一个轨道面的起始赤经角（初始RAAN）|     | float64 | 度   | \-D<br>\--delta-init | initial delta | 0   | \[0, 360) |
| 轨道倾角 |     | float64 | 度   | \-I<br>\--inclination | orbit-inclination | 50  | \[0, 180\] |
| 命名规则 |     | string |     | \--name-rule | name-rule |     | 目前只有wxwx-1 |
|     |     |     |     |     |     |     |     |
| 模版路径 | ✅   | string |     | \-t<br>\--template |     | -  |     |
| tle文件路径 | ✅   | string |     | \-i<br>\--input |     | -  |     |
| Kubenet config输出路径 | ✅   | string |     | \-o<br>\--output |     | -  |     |
| offset |     | int32 |     | \--offset | offset | 0   | 正负不超过轨道面卫星数量 |
| fake |     | bool |     | \--fake | fake | false |     |
| cross sphere |     | bool |     | \--cross-sphere | cross-sphere | false |     |
| ground access |     | bool |     | \--ground-access | ground-access | false |     |
| access intf num |     | int |     | \--access-intf-num |     | 16  |     |
| routing area |     | int32 |     | \-r<br>\--routing-area | routing-area | -   |     |
| gen-flow |     | string |     | \--gen-flow |  gen-flow   | ""  | 指定flow task生成模式，可选值：walker-neighbour、random   |
| image |     | string |     | \--image |     | ""  |     |
| node name prefix |     | string |     | \--node-name-prefix |     | ""  |     |
| 调度策略 | | string | | | schedule | "" | 枚举值：Random, EasyOrbit