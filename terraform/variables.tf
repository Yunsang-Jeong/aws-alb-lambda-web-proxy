variable "context" {
  description = "The name of service."
  type        = map(string)
  default = {
    delimiter = "-"
  }
}

variable "service_name" {
  description = "The name of service."
  type        = string
  default     = "lambda-web-proxy"
}
