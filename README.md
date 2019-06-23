# Kubeless YAML
The easiest way to deploy kubeless functions using YAML files. *You need to [install and setup kubeless](https://kubeless.io/docs/quick-start/) first.*

```
$ go get -u github.com/herlon214/kubeless-yaml
```

```
$ kubeless-yaml
Deploy your kubeless functions using yaml files

Usage:
  kubeless-yaml [command]

Available Commands:
  deploy      Deploy folders recursively that contains kubeless.yaml configs
  help        Help about any command

Flags:
  -h, --help   help for kubeless-yaml

Use "kubeless-yaml [command] --help" for more information about a command.
```

## Getting Started

You can deploy recursively the `example/` folder using one command, it will seek for `kubeless.yaml` files within the folders:
```
$ kubeless-yaml deploy -f example

INFO[0000] Getting deployed functions list...
INFO[0000] Found 0 deployed functions...
INFO[0000] -> [example/python]
INFO[0000] Deploying function 'python-hello'
time="2019-06-23T10:56:17+02:00" level=info msg="Deploying function..."
time="2019-06-23T10:56:17+02:00" level=info msg="Function python-hello submitted for deployment"
time="2019-06-23T10:56:17+02:00" level=info msg="Check the deployment status executing 'kubeless function ls python-hello'"

INFO[0000] Deploying function 'python-world'
time="2019-06-23T10:56:17+02:00" level=info msg="Deploying function..."
time="2019-06-23T10:56:17+02:00" level=info msg="Function python-world submitted for deployment"
time="2019-06-23T10:56:17+02:00" level=info msg="Check the deployment status executing 'kubeless function ls python-world'"

INFO[0000] -> [example/python/other]
INFO[0000] Deploying function 'python-hello2'
time="2019-06-23T10:56:18+02:00" level=info msg="Deploying function..."
time="2019-06-23T10:56:18+02:00" level=info msg="Function python-hello2 submitted for deployment"
time="2019-06-23T10:56:18+02:00" level=info msg="Check the deployment status executing 'kubeless function ls python-hello2'"

INFO[0000] -> [example/python/other/other]
INFO[0001] Deploying function 'python-hello3'
time="2019-06-23T10:56:18+02:00" level=info msg="Deploying function..."
time="2019-06-23T10:56:18+02:00" level=info msg="Function python-hello3 submitted for deployment"
time="2019-06-23T10:56:18+02:00" level=info msg="Check the deployment status executing 'kubeless function ls python-hello3'"

INFO[0001] All done!
```

Your `kubeless.yaml` file can contain many function configs:
```yaml
- name: "python-hello"
  version: "1.0.0"
  runtime: python2.7 # Run "kubeless get-server-config" to see your supported runtimes
  file: test.py
  handler: test.hello

- name: "python-world"
  version: "1.0.0"
  runtime: python2.7
  file: test.py
  handler: test.world
```


For each function it will run a preflight and compare the name and checksum with the deployed functions. If the function's name is found on the deployed functions but the checksum doesn't match, it will update the deployed function. If the function's name is not found it will be deployed normally. It will also skip a function if the deployed's function checksum is the same of the function being deployed.

```
$ kubeless-yaml deploy -f example
INFO[0000] Getting deployed functions list...
INFO[0000] Found 4 deployed functions...
INFO[0000] -> [example]
INFO[0000] -> [example/python]
INFO[0000] Skipping function 'python-hello' due to no changes
INFO[0000] Skipping function 'python-world' due to no changes
INFO[0000] -> [example/python/other]
INFO[0000] Skipping function 'python-hello2' due to no changes
INFO[0000] -> [example/python/other/other]
INFO[0000] Skipping function 'python-hello3' due to no changes
INFO[0000] All done!
```