# git-prompt-go
Informative and fast Git prompt for any shell (Bash, Zsh, and PowerShell).

This project is inspired from oliviervedier's [zsh-git-prompt](https://github.com/olivierverdier/zsh-git-prompt).

# Usage (development)

for Zsh user:
```zsh
PROMPT='%~ $(git-prompt-go) >"
```

for PowerShell user:
```ps1
function prompt {
  write-host "$(pwd) " -nonewline
  write-host "$(git-prompt-go)" -nonewline
  return "`n> "
}
```

## TODO
- [ ] ステータスの取得の一部を libgit2 で置き換えて高速化する
  * 特に `git stash` が Git for Windows だと遅い（シェルスクリプトで実装されているため）

- [ ] （本来の目的である） Rust 版への移植

## License
This software is released under the MIT license.
See [LICENSE](LICENSE) for details.
