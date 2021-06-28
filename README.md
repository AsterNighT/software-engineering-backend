# Backend for ZJU 2021 SE Project

> Should be renamed after

[![CI](https://github.com/AsterNighT/software-engineering-backend/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/AsterNighT/software-engineering-backend/actions/workflows/ci.yml)

## How to run it?
- run `make` on linux (if you work on windows and you have got no `make`, you have to read the Makefile and do it manually)
- open in your browser http://localhost:12448/swagger to check the swagger doc.

## Technology stack

Developing

- Language: https://golang.org/
- Orm: https://gorm.io/index.html
- Validator: https://github.com/go-playground/validator
- Unit test: https://github.com/stretchr/testify
- Config: https://github.com/spf13/viper
- Framework: https://github.com/labstack/echo
- API Definition: https://github.com/swaggo/swag and https://github.com/swaggo/echo-swagger
- Auth: https://github.com/dgrijalva/jwt-go

CI

- Github actions

- Linting: https://github.com/golangci/golangci-lint

CD

- Github actions
- Direct deployment

## Where can I learn relevant things?

- Golang: https://tour.golang.org/welcome/1
- Git: https://learngitbranching.js.org/?locale=zh_CN
- Basic http: https://developer.mozilla.org/zh-CN/docs/Web/HTTP
- Packages: Read their github repo and go document. Search it on google for detailed examples.

## Recommended development environment

- vscode (or goland)
- golang 1.16.1

## How to contribute?

If you are a group member:

1. `Fork` this repo to your account.
2. `git clone` the repo you cloned to your computer.
3. `git remote add` The original repo (https://github.com/2021-ZJU-Software-engineering/software-engineering-backend) as `upstream`.
4. Always `git pull upstream master` on `master` branch before you start developing.
5. `git checkout -b` to a new branch.
6. Write your code and tests.
7. Test it.
8. Commit your code with `-s` for signed-off
9. Submit a pull request to the [original repo](https://github.com/2021-ZJU-Software-engineering/software-engineering-backend) on your group branch.
10. Get review from your group leader and partners.
11. If you are not getting a LGTM, goto 6 and rewrite your code.
12. Get a LGTM and merge your code. 

If you are a group leader, in addition to the above, you need to:

1. Submit a pull request from your group branch to master. 
2. Get review from [@AsterNighT](https://github.com/AsterNighT) .
3. If you are not getting a LGTM, go to the code writer and ask him/her to fix it, goto 2.
4. Get a LGTM and merge your code.

## How to do code review?

Go through the code line by line carefully, start a review and ask yourself questions:

1. Is the code clear enough? If not, does it get enough comments?
2. Is the code badly composed? It there an obvious way to make it better?
3. Is the code necessary? 
4. It the code error-prone? Is it fully tested, manually or automatically?
5. Is there any unit test if necessary?

Don't hesitate to raised questions. It is the author's duty to answer them. Even the mildest suggestion matters. If you don't raise your question now, you may be questioned by [@AsterNighT](https://github.com/AsterNighT) when you try to merge your code into master.

If you feel good enough about the code, apply an approve and comment "LGTM" (Look good to me). This is common practice. If not, feel free to comment and request change from the author.

If the author has fixed the code, he/she can @the reviewer with the word "PTAL" (Please take another look, common practice, too) and ask for a new review.



## Precautions and questions you may ask

- **Stick to a good naming style**, refer to https://github.com/kettanaito/naming-cheatsheet
- **Stick to good code style**. I'll enforce some of them here:
  - Use `\n` as line separator
  - Use camelCase for names
- **Do not just throw your code somewhere**, do good module partition and refer to https://github.com/golang-standards/project-layout for layout.
- Always use English. And use good English for everyone's sake.
- The project is meant to be developed and built on linux tool chains. Its Makefile and similar scripts serve this purpose. It may or may not run properly on windows. Handle it yourself if you work on windows.
- Write markdown for documents in this repo. Microsoft Doc is not acceptable.
- If you have trouble cloning from git, refer to https://doc.fastgit.org/zh-cn/#%E5%85%B3%E4%BA%8E-fastgit
- If you have trouble downloading go modules, refer to https://goproxy.io/zh/

