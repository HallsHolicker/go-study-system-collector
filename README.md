<h1>Collect infromation Servers, Network Devices</h1>

고 스터디를 진행하면서 서버 정보와 네트워크 정보를 수집할 수 있는 프로젝트를 진행하였습니다.

프로그램은 2가지 방식으로 동작합니다.
* Collector mode
  * ssh mode ( Network, Servers)
  * api mode ( Only Servers - 서버에 server mode로 실행 중일 경우 가능 )
* Server mode

<H3>Requirement Target Server package</h3>
* lshw
* nvme

<h3>Collector mode</h3>
 - Server, Network 장비의 정보를 수집

<h3>Server mode</h3>
 - Server 에서 동작하며, 서버의 정보를 수집하여 API로 응답해 줍니다.

<h2>실행 방법</h2>
```<build name> <options>```

|옵션|설명|
|---|---|
|-server| 서버 모드 실행 |

참고 자료:
1. https://gist.github.com/FZambia/b5f5dde1bebb70a3f790