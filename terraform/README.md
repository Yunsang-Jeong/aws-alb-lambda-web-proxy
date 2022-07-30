# aws-alb-lambda-web-proxy

<!-- BEGIN_TF_DOCS -->
## Providers

| Name | Version |
|------|---------|
| <a name="provider_archive"></a> [archive](#provider\_archive) | 2.2.0 |
| <a name="provider_aws"></a> [aws](#provider\_aws) | 4.22.0 |
| <a name="provider_http"></a> [http](#provider\_http) | 2.2.0 |
| <a name="provider_null"></a> [null](#provider\_null) | 3.1.1 |
## Requirements

No requirements.
## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_context"></a> [context](#input\_context) | The name of service. | `map(string)` | <pre>{<br>  "delimiter": "-"<br>}</pre> | no |
| <a name="input_proxy_url"></a> [proxy\_url](#input\_proxy\_url) | The proxy url. | `string` | n/a | yes |
| <a name="input_service_name"></a> [service\_name](#input\_service\_name) | The name of service. | `string` | `"lambda-web-proxy"` | no |
## Outputs

| Name | Description |
|------|-------------|
| <a name="output_alb_dns_name"></a> [alb\_dns\_name](#output\_alb\_dns\_name) | n/a |
<!-- END_TF_DOCS -->