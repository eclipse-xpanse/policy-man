<p align='center'>
<a href="https://github.com/eclipse-xpanse/policy-man/actions/workflows/ci.yml" target="_blank">
    <img src="https://github.com/eclipse-xpanse/policy-man/actions/workflows/ci.yml/badge.svg" alt="build">
</a>

<a href="https://opensource.org/licenses/Apache-2.0" target="_blank">
    <img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="coverage">
  </a>
</p>

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
  -c, --config string      Specify the config file
  -h, --help               help for policy-man
  -a, --host string        The host of the HTTP server (default "localhost")
      --log.level string   The level of the log (default "info")
      --log.path string    The path of the log (default "stdout")
  -m, --mode string        The mode of the HTTP server.[release/debug/test] (default "release")
  -p, --port string        The port of the HTTP server (default "8090")
  -v, --version            Show the version number
```

### Evaluate the input by a policy list

Only `allow` and `deny` will be evaluated. If the variable `allow` be evaluated as false, or the variable `deny` be 
evaluated as true, The policy will be evaluated as false.

```shell
$ curl -X POST http://localhost:8090/evaluate/policies -H 'Content-Type: application/json' -d '
{
    "policy_list": [
        "import future.keywords.if\nimport future.keywords.in\n\ndefault allow := false\n\nallow if {\n    input.method == \"GET\"\n    input.path == [\"salary\", input.subject.user]\n}\n\nallow if is_admin\n\nis_admin if \"admin\" in input.subject.groups",
        "import future.keywords.if\nimport future.keywords.in\n\ndefault deny := false\n\nallow if {\n    input.method == \"GET\"\n    input.path == [\"salary\", input.subject.user]\n}\n\nallow if is_admin\n\nis_admin if \"admin\" in input.subject.groups"
    ],
    "input": "{\"method\":\"GET\",\"path\":[\"salary\",\"bob\"],\"subject\":{\"user\":\"bob\",\"groups\":[\"sales\",\"marketing\"]}}"
}'
 
{"isSuccessful":true}
```

### Use Swagger UI

Open internet browser and navigate to the url [http://localhost:8090/swagger/index.html](http://localhost:8090/swagger/index.html).
View and Call APIs on the page of swagger UI.

### Update OpenAPI documentation

All files of the RESTful API documentation are in the directory [./openapi/docs](./openapi/docs), when the service API
or API annotations are updated, these files should be updated by the following commands:

```shell
make api_doc
```

All the above commands are written to the file [Makefile](./Makefile), You can also use commands in the chapter
[Build from source](#Build from source) directly to update these files.

## Dependencies File

All third-party related content is listed in the [DEPENDENCIES](DEPENDENCIES) file.