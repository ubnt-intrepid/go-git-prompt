# go-git-prompt

[![Build Status](https://travis-ci.org/ubnt-intrepid/go-git-prompt.svg?branch=master)](https://travis-ci.org/ubnt-intrepid/go-git-prompt)

Informative and fast Git prompt for any shell (Bash, Zsh, and PowerShell).

This project is inspired from oliviervedier's [zsh-git-prompt](https://github.com/olivierverdier/zsh-git-prompt).

## Usage
Bash:
```bash
PS1='\w $(go-git-prompt) % '
```

Zsh:
```zsh
PROMPT='%~ $(go-git-prompt) %% '
```

Fish:
```fish
function fish_prompt
   echo (go-git-prompt)" % "
end
```

PowerShell:
```ps1
function prompt {
  write-host "$(pwd) " -nonewline
  write-host (go-git-prompt) -nonewline
  return "`n> "
}
```

## Install
```shell-session
$ go get -v github.com/ubnt-intrepid/go-git-prompt
```

## License
This software is released under the MIT license.
See [LICENSE](LICENSE) for details.
