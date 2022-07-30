module "network" {
  source = "github.com/Yunsang-Jeong/terraform-aws-network"
  
  vpc_cidr_block = "10.0.0.0/16"
  vpc_enable_dns_hostnames	= true
  vpc_enable_dns_support	 = true
  create_igw = true
  subnets = [
  {
      identifier            = "public-a"
      name_tag_postfix      = "pub-a"
      availability_zone     = "ap-northeast-2a"
      cidr_block            = "10.0.104.0/24"
      enable_route_with_igw = true
      create_nat            = true
  },
  {
      identifier            = "public-c"
      name_tag_postfix      = "pub-c"
      availability_zone     = "ap-northeast-2c"
      cidr_block            = "10.0.105.0/24"
      enable_route_with_igw = true
  },
  {
      identifier            = "private-a"
      name_tag_postfix      = "pri-a"
      availability_zone     = "ap-northeast-2a"
      cidr_block            = "10.0.106.0/24"
      enable_route_with_nat = true
  },
  {
      identifier            = "private-c"
      name_tag_postfix      = "pri-c"
      availability_zone     = "ap-northeast-2c"
      cidr_block            = "10.0.107.0/24"
      enable_route_with_nat = true
  },
  ]
}