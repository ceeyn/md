在Go代码中，通过 `os.Getenv("SERVER_NAME")` 来获取环境变量 `SERVER_NAME` 的值。这意味着 `SERVER_NAME` 必须在操作系统的环境变量中定义。环境变量通常通过以下几种方式设置和获取：

### 1. 在操作系统的命令行中设置

#### Linux/MacOS

你可以在运行程序之前通过命令行设置环境变量。例如：

```
sh
复制代码
export SERVER_NAME=my_server
```

然后运行你的Go程序：

```
sh
复制代码
go run main.go
```

#### Windows

你可以在命令行窗口中设置环境变量：

```
sh
复制代码
set SERVER_NAME=my_server
```

然后运行你的Go程序：

```
sh
复制代码
go run main.go
```

### 2. 在Shell脚本中设置

你可以将环境变量设置放在一个Shell脚本中，然后运行该脚本。例如：

```
sh
复制代码
#!/bin/bash
export SERVER_NAME=my_server
go run main.go
```

### 3. 在Docker容器中设置

如果你的应用运行在Docker容器中，你可以在Dockerfile或docker-compose文件中设置环境变量。

#### Dockerfile

```
Dockerfile
复制代码
ENV SERVER_NAME=my_server
```

#### docker-compose.yml

```
yaml
复制代码
version: '3'
services:
  my_service:
    image: my_image
    environment:
      - SERVER_NAME=my_server
```

### 4. 在Kubernetes中设置

如果你的应用运行在Kubernetes集群中，你可以在Pod或Deployment的配置文件中设置环境变量。

```
yaml
复制代码
apiVersion: v1
kind: Pod
metadata:
  name: mypod
spec:
  containers:
  - name: mycontainer
    image: myimage
    env:
    - name: SERVER_NAME
      value: "my_server"
```

### 5. 在程序启动时通过代码设置

虽然这种方式不常见，但你也可以在程序启动时通过代码设置环境变量：

```
go
复制代码
os.Setenv("SERVER_NAME", "my_server")
```

然后再获取它：

```
go
复制代码
serverName := os.Getenv("SERVER_NAME")
```

### 示例代码

以下是一个示例，展示了如何获取环境变量 `SERVER_NAME`：

```
go
复制代码
package main

import (
	"fmt"
	"os"
)

func main() {
	serverName := os.Getenv("SERVER_NAME")
	if serverName == "" {
		fmt.Println("SERVER_NAME is not set")
	} else {
		fmt.Println("SERVER_NAME:", serverName)
	}
}
```

### 