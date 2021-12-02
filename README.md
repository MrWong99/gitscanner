# Git Repository checker

This is a fast Go implementation to check Git repositories (local or remote) for some common security issues.
It relies heavily on [go-git](https://github.com/go-git/go-git).

## Usage

```
$ ./gitscanner --help
Usage of ./gitscanner:
  -repositories string
        A comma-separated list of repositories to perform checks against. Can be in these formats:
        - http(s)://<remote URL>
        - git@<remote URL>
        - file://<path>
  -branch-pattern string
        Optional pattern to match refs against. Only matches will be processed in checks that rely on refs.
  -username string
        An optional username for http basic auth.
  -password string
        An optional password for http basic auth.
  -ssh-private-key-file string
        An optional path to a SSH private key file in PEM format.
  -ssh-private-key-password string
        An optional password if the given private key file is encrypted.
```

## Performed checks

* **[binaryfile.SearchBinaries](/checks/binaryfile):** Searches for any binary files on each branch (local or remote) that matches the `-branch-pattern`.
* **[unicode.SearchUnicode](/checks/unicode):** Searches for specific unicode characters in each file on each branch (local or remote) that matches the `-branch-pattern`. See [trojan-source.pdf](https://trojansource.codes/trojan-source.pdf).
