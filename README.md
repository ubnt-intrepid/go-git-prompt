# git-prompt-go
Informative and fast Git prompt for any shell (Bash, Zsh, and PowerShell).

This project is inspired from oliviervedier's [zsh-git-prompt](https://github.com/olivierverdier/zsh-git-prompt).

# Usage (development)

for Zsh user:
```zsh
PROMPT='%~ $(git-prompt-go) >'
```

for PowerShell user:
```ps1
function prompt {
  write-host "$(pwd) " -nonewline
  write-host "$(git-prompt-go)" -nonewline
  return "`n> "
}
```

## License
This software is released under the MIT license.
See [LICENSE](LICENSE) for details.
