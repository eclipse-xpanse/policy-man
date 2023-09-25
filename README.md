# policy-man
OPA based policy engine.

## How to use

### Build from source

```shell
make build
```

### Start the policy-man


```shell
./policy-man -h

  ___  ___  _    ___ _____   __        __  __   _   _  _
 | _ \/ _ \| |  |_ _/ __\ \ / /  ___  |  \/  | /_\ | \| |
 |  _/ (_) | |__ | | (__ \ V /  |___| | |\/| |/ _ \| .' |
 |_|  \___/|____|___\___| |_|         |_|  |_/_/ \_\_|\_|

Usage:
  policy-man [flags]

Flags:
  -c, --config string      config file (default is ./config.yaml)
  -h, --help               help for policy-man
  -a, --host string        The host of the HTTP server
      --log.level string   The level of the log (default "warn")
      --log.path string    The path of the log (default "stdout")
  -p, --port string        The port of the HTTP server
```

## Evaluate the input by a policy list

Only `allow` and `deny` will be evaluated. If the variable `allow` be evaluated as false, or the variable `deny` be 
evaluated as true, The policy will be evaluated as false.

```shell
$ curl -X POST http://localhost:8080/evaluate/policies -H 'Content-Type: application/json' -d '
{
    "policy_list": [
        "import future.keywords.if\nimport future.keywords.in\n\ndefault allow := false\n\nallow if {\n    input.method == \"GET\"\n    input.path == [\"salary\", input.subject.user]\n}\n\nallow if is_admin\n\nis_admin if \"admin\" in input.subject.groups",
        "import future.keywords.if\nimport future.keywords.in\n\ndefault deny := false\n\nallow if {\n    input.method == \"GET\"\n    input.path == [\"salary\", input.subject.user]\n}\n\nallow if is_admin\n\nis_admin if \"admin\" in input.subject.groups"
    ],
    "input": "{\"method\":\"GET\",\"path\":[\"salary\",\"bob\"],\"subject\":{\"user\":\"bob\",\"groups\":[\"sales\",\"marketing\"]}}"
}'
 
{"isSuccessful":true}
```