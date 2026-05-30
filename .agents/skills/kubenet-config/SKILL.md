---
name: kubenet-config
description: 编写与编辑 kubenet 天空地海一体化网络仿真配置；引用 OpenAPI 配置定义进行字段对齐；调用 kubenet-config-generator 生成 TLE 与 kubenet YAML 并完成核验。关键词: kubenet, walker, tle, shell-config, yaml, openapi, 配置校验, 生成, merge
---

## Skill Purpose
该 skill 用于在本仓库内完成 kubenet 配置的编写、修订、生成与核验，覆盖以下场景：
- 根据需求创建或修改 shell 配置 YAML（多壳层）。
- 参考 OpenAPI 定义对字段语义做对齐。
  - 路径位于 `go list -m -f '{{.Dir}}' dslab.sjtu/kubenet/api` 下的 `openapi/kubenet-config.yaml`，完整路径例如`/home/node1/go/pkg/mod/git.dslab.cc/dslab/kubenet-api.git@v1.0.40/openapi/kubenet-config.yaml`。
- 调用 `kubenet-config-generator` CLI 生成 TLE 和 kubenet config。
- 对生成结果执行可复现的核验流程（命令成功、输出存在、基础结构正确）。

## When To Use
当用户需求包含任一关键词时应触发：
- “kubenet 配置编写/编辑”
- “walker 壳层参数调整”
- “根据 OpenAPI 对齐配置字段”

## Repository Facts (Do Not Invent)
- CLI 入口：`go run ./main.go`
- 生成命令组：`generate walker tle`、`generate walker kubenet`
- 合并命令：`merge`
- Shell YAML 关键字段定义源：`pkg/shell_config/shell_config.go`
- 示例配置：`config/multi_shell.yaml.example`
- 模板示例：`template/kubenet_template.yaml.example`
- 批量生成参考：`.mk/gen.mk`

## OpenAPI Reference Rules
1. 若用户给出 OpenAPI 文档路径或 URL：
- 优先按 OpenAPI 的字段定义、类型、约束进行配置生成。
- 对照本仓库可落地字段，仅写入 CLI/代码真实支持的字段。

2. 若用户未给出 OpenAPI 文档：
- 明确说明“当前仓库未发现本地 openapi.yaml/json”。
- 采用本地结构体作为临时事实源：`pkg/shell_config/shell_config.go`。
- 若用户提到可视化/仿真器 API，可提示其提供规范来源（例如 kubenet visual API 的 OpenAPI 文档）后再做严格对齐。

3. 对齐策略：
- 字段名冲突时，以“OpenAPI 语义 + 本地 CLI 可执行性”双重约束决策。
- 对无法落地到 CLI 的字段，必须在答复中标注为“暂无法生效”。

## Tools to Use

### kubenet-config-generator

kubenet-config-generator 是一个 CLI 工具，用于生成 kubenet 配置文件。
它支持从命令行输入参数进行执行，也支持配合详细的配置文件进行执行，并且可以生成 TLE 文件和 kubenet 配置文件。
如果用户没有指定，该工具通过在项目根目录下运行`go run main.go`调用。
如果不确定用法，可以通过`-h`查看所有命令及其用法，可以递归查看子命令用法。

#### 示例

##### 合并多个 kubenet 配置文件

```bash
go run ./main.go merge -i tasks/20200202/inputs/s1.yaml -i tasks/20200202/inputs/s2.yaml -o tasks/20200202/intermediate/kubenet/merges1s2.yaml
```

##### 生成单壳层 TLE 和 kubenet 配置

编写 shell-config.yaml

```yaml
shells:
  - shell-index: 1 # 和 generate walker kubenet 命令里的 -s 参数一致
    satellite-num: 100 # -N
    plane-num: 10 # -P
    needed-plane-num: 10 # -E
    orbit-inclination: 53 # -I
    height: 550 # -H
    phase: 0 # -F
    initial-delta: 0 # -D
    offset: -1 # --offset
    fake: false # --fake
    cross-sphere: true # --cross-sphere
    ground-access: true # --ground-access
    routing-area: 1 # -r
    schedule: EasyOrbit
    gen-flow: walker-neighbour
    flow-bandwidth: 1Gbps
    name-rule: wxwx-1 # --name-rule, 目前只有 wxwx-1 可选
```

从`template/kubenet_template.yaml`拷贝模板到`tasks/20200202/intermediate/kubenet/s1_template.yaml`，并根据需要修改。

```bash
# 1) 生成 TLE
go run ./main.go generate walker tle \
	-c shell-config.yaml \
	-s 1 \
	-o tasks/20200202/intermediate/tle/s1.tle

# 2) 生成 kubenet 配置
go run ./main.go generate walker kubenet \
	-c shell-config.yaml \
	-s 1 \
	-t tasks/20200202/intermediate/kubenet/s1_template.yaml \
	-i tasks/20200202/intermediate/tle/s1.tle \
	-o output/kubenet/s1.yaml
```

### OpenAPI-based Configuration Editing

参考 OpenAPI 配置定义，按照用户需求编写或修改 YAML 配置文件，确保字段语义与 OpenAPI 定义对齐。
生成文件放在`工作目录/input`下。

## Editing Constraints
- 不得凭空新增仓库不支持的 CLI 参数。
- 不得修改与本次任务无关的配置文件。
- 变更配置时优先最小修改，保留用户已有字段顺序和注释风格。
- 用户若无指定，需要把最终的配置文件生成在`output/kubenet`目录下，并保持与输入文件相似的命名规范。
- 所有中间文件，包括任何 TLE 文件和 YAML 文件，均保存在`工作目录/intermediate`目录下，并保持与输入文件相似的命名规范。
- 不得以任何方式编辑中间文件和结果文件，除非用户明确要求或确认。

## Authoring Workflow
按以下顺序执行，不跳步：

1. 创建工作目录
在 `output/tasks` 下创建以当前日期时间命名的子目录，例如 `output/tasks/2024-06-01_15-30-00`，用于存放本次任务的所有输出文件。

2. 需求分析
对用户需求进行分析，确定为单步执行还是多步执行模式。
- 单步执行旨在通过一次 kubenet-config-generator 最小闭环序列生成配置文件或者借助 OpenAPI 定义直接编写目标配置文件。
- 多步执行旨在通过多次迭代生成和修改配置文件（即对单步执行的多次组合），最终完成复杂配置的生成。

在工作目录下创建 `planning.md` 文件，记录需求分析结果、执行模式选择理由以及初步的执行步骤规划。
等待用户审阅、修改该文件后，再继续后续步骤。

3. 规划执行步骤
- 单步执行
  - 选项A：借助 kubenet-config-generator 生成，参考最小闭环序列示例。
  - 选项B：仅依照 OpenAPI 定义编写/修改配置。
- 多步执行
  - 根据需求确定执行步骤构成，向用户展示，需要获得用户确认。
  - 需要分解生成粒度，必要时先生成多个 YAML，再通过 kubenet-config-generator 执行 merge。

对 `planning.md` 文件进行更新，详细列出每一步的目标、输入输出文件、预期命令等信息。
等待用户审阅、修改该文件后，再继续后续步骤。

4. 开始执行最小步骤
每一最小步骤需要向用户展示简要执行结果，并等待用户确认。

5. 报告结果
- 返回修改了哪些配置文件。
- 返回执行的关键命令。
- 返回核验结论（成功/失败、失败原因、下一步建议）。

## 多步执行思路

将复杂的场景需求进行分解，例如：对于一个空天地海复杂场景，要分解成空、天、地、海四个大部分，每个部分再单独拆解，
例如空间段，根据需要，拆解为多个星座、多个壳层，每个星座就可以编写单独的`shell-config.yaml`，生成单独的 TLE 和 kubenet 配置文件，最后通过 `merge` 命令合并成一个完整的 kubenet 配置文件。

你需要把你的思路和用户进行讨论，得到确认后再继续。

## Verification Checklist
至少完成以下核验并在答复中写明：
- 命令退出码为 0。
- 目标输出文件存在且非空。
- 生成日志无 `failed`/`error` 关键失败信息。
- 若用户提供 OpenAPI：列出已对齐字段与未对齐字段。

可选增强核验：
- 对输出 YAML 做结构检查（包含预期节点数量、链路字段、命名前缀等）。
- 若环境可用，调用仿真器管理接口进行 apply 前静态校验。

## Response Contract
执行此 skill 时，回复建议按以下结构输出：
1. 本次目标
2. 配置改动
3. 执行命令
4. 核验结果
5. 待用户确认项（如有）

## Notes
- `cmd/generate/ground` 当前为占位命令，不作为主要生成路径。
- 若用户要求“严格 OpenAPI 校验”，但未提供可访问规范文件，应先请求规范来源，再继续实施。
