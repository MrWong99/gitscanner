# Git Repository security checker

This is a fast Go implementation to check Git repositories (local or remote) for some common security issues.
It relies heavily on [go-git](https://github.com/go-git/go-git).

## Usage

To simply start the scanner as webservice configure a port to use and optionally (highly recommended) a TLS certificate (chain) and private key file:

`./gitscanner -port 16092 -ssl-certificate-chain-file ~/.tls/my-cert-chain.crt.pem -ssl-private-key-file ~/.tls/tls.key.pem`

![UI Example](/ui_example.png)

If you instead just want to run in single execution mode with CLI don't provide any port and then refer to the possible configuration parameters:

```
$ ./gitscanner --help
Usage of ./gitscanner:
  -branch-pattern string
        Optional pattern to match refs against. Only matches will be processed in checks that rely on refs.
  -email-pattern string
        Pattern to match all commiter and author emails against. This will be used for the commitmeta.CheckCommits check.
  -name-pattern string
        Pattern to match all commiter and author names against. This will be used for the commitmeta.CheckCommits check.
  -password string
        An optional password for http basic auth.
  -port int
        When provided this will startup a webserver including ui that can be used to perform the checks via browser. (default -1)
  -repositories string
        A comma-separated list of repositories to perform checks against. Can be in these formats:
        - http(s)://<remote URL>
        - git@<remote URL>
        - file://<path>
  -ssh-private-key-file string
        An optional path to a SSH private key file in PEM format.
  -ssh-private-key-password string
        An optional password if the given private key file is encrypted.
  -ssl-certificate-chain-file string
        An optional path to a TLS certificate (chain) in PEM format to enable HTTPS. Only used when port is set.
  -ssl-private-key-file string
        An optional path to a TLS private key file in PEM format to enable HTTPS. Only used when port is set.
  -username string
        An optional username for http basic auth.
```

## Performed checks

* **[SearchBinaries](/checks/binaryfile/binaryfile.go):** Searches for any binary files on each branch (local or remote) that matches the `branchPattern`.
* **[SearchIllegalUnicodeCharacters](/checks/unicode/unicode.go):** Searches for specific unicode characters in each file on each branch (local or remote) that matches the `branchPattern`. See [trojan-source.pdf](https://trojansource.codes/trojan-source.pdf).
* **[CheckCommitMetaInformation](/checks/commitmeta/commitmeta.go):** Checks every commits author and committer name and email for expected match against `emailPattern` and `namePattern`.
* **[SearchBigFiles](/checks/filesize/filesize.go):** Searches for files on each branch (local or remote) that matches the `branchPattern` and is bigger than `filesizeThresholdByte`.
  
## Build locally

1. Install Go.
2. `go build .`

## Add new tests

Adding tests is very simple:

1. Write a type implementing the [Checker interface](checks/checks.go#L39)
```go
package myawesometest

import (
    mygit "github.com/MrWong99/gitscanner/git"
    "github.com/MrWong99/gitscanner/utils"
)

type MyTest struct{
}

func (*MyTest) Check(wrapRepo *mygit.ClonedRepo, output chan<- *utils.SingleCheck) error {
    defer close(output)
    // perform checks here and write any found issues into the output channel
}
```
2. Add the function to the list of possible checks in [main.go](main.go#L45-47)

## REST API

When started in server mode *gitscanner* will provide the following endpoints:

### POST /api/v1/checkRepos - Perform checks for given paths

**Status Codes:**

* `200`: checks were performed. Singular could still have failed though.
* `400`: the request body was malformed.

**Request Body:**

* `path`: a comma separated list of urls to clone. They can be in these formats:
  * `http(s)://<remote URL>`
  * `git@<remote URL>`
  * `file://<path>` -> will only search on the local filesystem of the server
* `checks`: a list of check identifiers to determine with checks are to be performed. See [GET checks](#get-apiv1checks---retrieve-the-list-of-possible-checks).

*Example:*

```json
{
  "path": "git@github.com:go-git/go-git.git,https://gitlab.com/gitlab-org/gitlab.git",
  "checkNames": [
    "SearchBinaries",
    "CheckCommitMetaInformation",
    "SearchIllegalUnicodeCharacters",
    "SearchBigFiles"
  ]
}
```

**Response Body:**

* Array of objects with:
  * `date`: ISO encoded timestamp of when the check was started
  * `repository`: the repository that was checked
  * `error`: if any error occured while opening the repo it will be contained here, else empty string
  * `checks`: list of checks that contained suspicious results. Each check consists of:
    * `origin`: where this issue was found. This can be multiple things, e.g. a path to a file or a commit hash.
    * `branch`: the branch (if any) on which the issue was found.
    * `checkName`: the name of the check that found this issue.
    * `acknowledged`: currently always false, will be used in future updates.
    * `additionalInfo`: list of non-specified key-value pairs that differ from check to check.

*Example:*

```json
[
    {
        "date": "2021-12-03T00:22:00.6155686+01:00",
        "repository": "git@github.com:MrWong99/micasuca.git",
        "error": "",
        "checks": [
            {
                "origin": "Commit 65508c0d5f0ea52ce3d93f77f471359f4ec1d1bc",
                "branch": "",
                "checkName": "CheckCommitMetaInformation",
                "acknowledged": false,
                "additionalInfo": {
                    "authorEmail": "shady.dude@inter.net",
                    "authorName": "Jeff",
                    "commitMessage": "Table view (in progress...)\n",
                    "commiterEmail": "shady.dude@inter.net",
                    "commiterName": "jeffHacker",
                    "numberOfParents": 1
                }
            },
            {
                "origin": "gradle/wrapper/gradle-wrapper.jar",
                "branch": "refs/remotes/origin/master",
                "checkName": "SearchBigFiles",
                "acknowledged": false,
                "additionalInfo": {
                    "filemode": "0100644",
                    "filesize": "54.3 kB"
                }
            },
            {
                "origin": "gradle/wrapper/gradle-wrapper.jar",
                "branch": "refs/remotes/origin/master",
                "checkName": "SearchBinaries",
                "acknowledged": false,
                "additionalInfo": {
                    "filemode": "0100644",
                    "filesize": "54.3 kB"
                }
            },
            {
                "origin": "gradlew",
                "branch": "refs/remotes/origin/master",
                "checkName": "SearchIllegalUnicodeCharacters",
                "acknowledged": false,
                "additionalInfo": {
                    "character":"'\\u202a'",
                    "filemode": "0100644",
                    "filesize": "5.3 kB"
                }
            }
        ]
    },
    {
        "date": "2021-12-03T00:22:01.0774318+01:00",
        "repository": "https://github.com/Mnaaz/JavaChat",
        "error": "",
        "checks": []
    }
]
```

### GET /api/v1/checkDefinitions - Retrieve the list of possible checks

**Status Codes:**

* `200`: Checks retrieved.

**Response Body:**

* List of strings for all registered checks

*Example:*

```json
[
    "SearchBinaries",
    "SearchIllegalUnicodeCharacters",
    "CheckCommitMetaInformation",
    "SearchBigFiles"
]
```

### GET /api/v1/config - Retrieve the current application configuration

**Status Codes:**

* `200`: configuration was returned.

**Response Body:**

* `branchPattern`: pattern to match branches against. The *SearchBinaries* and *SearchUnicode* checks use this.
* `namePattern`: pattern to match the commiter and author names against. The *CheckCommitAuthor* check uses this.
* `emailPattern`: pattern to match the commiter and author emails against. The *CheckCommitAuthor* check uses this.

*Example:*

```json
{
    "branchPattern": "refs/remotes/origin/main|.*v\\d+",
    "namePattern": ".*",
    "emailPattern": ".*@shady.com"
}
```

### PUT /api/v1/config - Set the application configuration

**Status Codes:**

* `200`: config was updated.
* `400`: the request body was malformed.

**Request Body:**

* `branchPattern`: regex pattern to match branches against. The *SearchBinaries* and *SearchUnicode* checks use this.
* `namePattern`: regex pattern to match the commiter and author names against. The *CheckCommitAuthor* check uses this.
* `emailPattern`: regex pattern to match the commiter and author emails against. The *CheckCommitAuthor* check uses this.

*Example:*

```json
{
    "branchPattern": "refs/remotes/origin/main|.*v\\d+",
    "namePattern": ".*",
    "emailPattern": ".*@shady.com"
}
```

### PUT /api/v1/config/sshkey - Set the ssh private key to use when using ssh during clone

**Status Codes:**

* `200`: config was updated.
* `400`: the request body was malformed.

**Request Body:**

* `key`: the ssh private key. Can be additionally encrypted with
* `password`: the password that this key was encrypted with if any

*Example:*

```json
{
    "key": "-----BEGIN RSA PRIVATE KEY-----\nyOut41nK1mdUMB?\n-----END RSA PRIVATE KEY-----",
    "password": ""
}
```

### PUT /api/v1/config/basicauth - Set the username and password when usic basic authentication during clone

**Status Codes:**

* `200`: config was updated.
* `400`: the request body was malformed.

**Request Body:**

* `username`: the username to use.
* `password`: the password to use.

*Example:*

```json
{
    "username": "SecureMan",
    "password": "1n5EcuR3"
}
```

### GET /api/v1/checks?from=<unix millis>&to=<unix millis>&checkNames=<stringArr> - Retrieve previously performed checks that are stored in DB

**Status Codes:**

* `200`: checks retrieved successfully
* `500`: checks could not be read from DB

**Query Parameters:**

* `from`: milliseconds since 1970-01-01 as start date from which checks should be included.
* `to`: milliseconds since 1970-01-01 as end date until which checks should be included.
* `checkNames`: comma-separated list of check names to include in the results.

**Response Body:**

Same as in [/api/v1/checkRepos](#post-apiv1checkrepos---perform-checks-for-given-paths).

### PUT api/v1/acknowledged/{checkID} - Set the acknowledged flag of given check

**Path Params:**

* `checkID`: the id of the check that should be set to acknowledged.

**Request Body:**

* `acknowledged`: boolean indecating wether this check was acknowledged or not.

*Example:*

```json
{
    "acknowledged": true
}
```
