# 声明
```text
此文件为模板文件，用于给其他服务生成AGENT.md 或者Claude.md 记忆规则使用
```
# 微服务工作模式
```text
AI 为主要任务完成角色，以这个条件为前提，设计AI工作规范，以达到多个服务，多个AI共同协作的目的
````
## AI 服务分类，分为三类
### 1. agent服务类
```text
此服务只有一个，服务命名规则 服务名-agent(word-hero-agent)
此服务是所有的服务总线这个服务担任产品经理，以及架构师的角色
负责给分析需求，架构设计，分发任务
包含以下必要文件
agent_service.yaml 这个文件记录所有的服务类别以及地址等相关信息,以及每个服务的角色和功能
agent_role.yaml 这个文件记录当前服务是那种类别,并且指定了遵循哪个agent服务
service_issues.md 记录每次提交的issue链接，后续追踪
如何分发任务，需求梳理完毕，任务梳理完毕后，应当得出具体功能由哪个服务完成，给那个服务提交issue(gh create issue)
```
### 2.api服务类
```text
agent_role.yaml 这个文件记录当前服务是那种类别,并且指定了遵循哪个agent服务
此服务只有一个，涉及到接口变更的 首先应该给这个服务提交issue,完成api的变更
此服务的任务从自身repo的issue读取，并完成,通过(gh issue)相关命令完成
```
### 3.service 其他具体服务类
```text
agent_role.yaml 这个文件记录当前服务是那种类别,并且指定了遵循哪个agent服务
多个具体服务，逐个读取issue,一个个完成issue 并且关闭issue就行了
此服务的任务从自身repo的issue读取，并完成,通过(gh issue)相关命令完成
```
# 以下规范性约定，所有服务都要遵循
```text
从 agent 了解所有服务信息，或者架构信息,访问地址参考agent_role.yaml 下的 agent_url
从 api 读取最新api协议
每次任务都从issue读取，只完成自己服务的部分，关闭issue
需要读取其他服务代码的时候，允许将代码下载到临时目录.tmp/repos/ 下，该目录只允许临时阅读代码
从agent服务的repo中读取 agent_service.yaml ，以获取其他服务相关信息
```
# 当前服务的服务类别: agent_service.yaml读取

```text
初始化AI 服务，请在agent_service填写完相关服务信息后，依据此文档初始化，生成后续遵从的规范
```