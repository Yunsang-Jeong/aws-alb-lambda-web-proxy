# aws-alb-lambda-web-proxy

아주 간단한 Web proxy를 구현하기 위해서 AWS ALB와 Lambda를 이용해 구현했습니다.

`terraform destroy` 과정에서 `Still destroying...`로 장기간 대기중인 상태로 빠지는 이슈가 있습니다. 이는 AWS Lambda에 의해 생성된 ENI에 대한 Lifecycle 제어가 불가능하기 때문입니다.
