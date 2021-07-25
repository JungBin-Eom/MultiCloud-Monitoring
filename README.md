# MultiCloud Monitoring
본 Repository는 OpenStack(오픈스택)과 CloudStack(클라우드스택)으로 이루어진 멀티 클라우드 환경의 로깅 시스템 코드를 저장하는 저장소입니다. 오픈스택 컴포넌트와 클라우드스택 서버의 로그를 분석하여 빠르게 오류를 파악하고 원활한 클라우드 환경을 운영하기 위함을 목표로 합니다.
## Logging
ELK Stack을 통해 오픈스택 컴포넌트와 클라우드스택 서버의 로그를 수집하고 분석합니다.   
<br>
<img src="./Logging Architecture.png" title="오픈스택 로깅 구조" alt="logging architecture"></img>
<br>
<br>
<img src="./Logging Architecture2.png" title="클라우드스택 로깅 구조" alt="logging architecture2"></img>
<br>

### APIs
#### [GET] /sync
- filebeat에서 설정한 각 로그 경로로부터 로그를 수집하고 logstash로 전송합니다.   
- logstash는 정의된 파이프라인을 통해 전달받은 로그를 정제하고 elasticsearch로 전송합니다.   
- elasticsearch에 저장된 로그들은 index pattern을 가지고 있으며 이를 기반으로 동기화를 수행합니다.   
- 데이터베이스에 수집한 로그를 중복되지 않도록 저장합니다.
   
#### [GET] /{component:[a-z]+}/log
- 경로에서 입력한 component에 해당하는 로그를 데이터베이스로부터 읽어옵니다.
- 로그 목록을 JSON 형태로 파싱하여 보여줍니다.
   
#### [DELETE] /{component:[a-z]+}/log
- 경로에서 입력한 component에 해당하는 로그를 데이터베이스에서 삭제합니다.
- elasticsearch에 저장한 index도 함께 삭제합니다.
- 에러 발생을 인지하고 문제되는 부분을 고친 후 리로드하는 역할을 수행합니다.
   
#### [GET] /check
- 현재 데이터베이스에 저장된 로그 중 ERROR 타입의 에러가 있는지 검사합니다.
- ERROR 타입의 로그를 가지고 있는 component 목록을 보여줍니다.
   
<br> 

## Monitoring
오픈스택 API와 클라우드스택 API를 통해 리소스 사용량을 확인하고 보여줍니다.
<br>
<img src="./Monitoring Architecture.png" title="오픈스택 모니터링 구조" alt="monitoring architecture"></img>
<br>
<br>
<img src="./Monitoring Architecture2.png" title="클라우드스택 모니터링 구조" alt="monitoring architecture2"></img>
<br>

### APIs
#### [POST] /token
- 오픈스택 사용자의 아이디와 비밀번호, Project ID를 body로 보내어 토큰을 요청합니다.
- 먼저 Unscoped 토큰 요청을 통해 올바른 사용자인지 확인합니다.
- 그후 Scoped 토큰 요청을 통해 오픈스택 서비스를 사용할 수 있는 토큰을 발급받습니다.
   
#### [GET] /instances
- 발급받은 토큰과 Project ID를 헤더에 포함하여 인스턴스 목록을 요청합니다.
- 해당 프로젝트에 생성된 인스턴스 목록을 보여줍니다.
- 현재는 오픈스택만 구현하였습니다.
   
#### [GET] /statistics
- 발급받은 토큰과 Project ID를 헤더에 입력합니다.
- body에는 클라우드스택의 apiURL을 구하기 위해 필요한 정보를 입력합니다.
- 그 후 각 플랫폼의 하이퍼바이저 리소스 사용량을 요청합니다.
- 해당 프로젝트에 사용되는 리소스 크기와 여유 리소스 크기를 보여줍니다.
