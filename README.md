# go-git-prompt
Informative and fast Git prompt for any shell (Bash, Zsh, and PowerShell).

This project is inspired from oliviervedier's [zsh-git-prompt](https://github.com/olivierverdier/zsh-git-prompt).

# Usage (development)

Bash:
```
PS1='\w $(go-git-prompt) % '
```

Zsh:
```
PROMPT='%~ $(go-git-prompt) %% '
```

Fish:
```fish
function fish_prompt
   echo (go-git-prompt)" % "
end
```

and PowerShell:
```ps1
function prompt {
  write-host "$(pwd) " -nonewline
  write-host (go-git-prompt) -nonewline
  return "`n> "
}
```

## License
This software is released under the MIT license.
See [LICENSE](LICENSE) for details.
