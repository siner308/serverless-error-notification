resource "aws_subnet" "public_a1" {
  vpc_id = aws_vpc.main.id
  cidr_block = "10.0.1.0/24"

  availability_zone = "ap-northeast-2a"

  tags = {
    Name = "ttbkk-subnet-public-a1"
  }
}

resource "aws_subnet" "public_b1" {
  vpc_id = aws_vpc.main.id
  cidr_block = "10.0.2.0/24"

  availability_zone = "ap-northeast-2b"

  tags = {
    Name = "ttbkk-subnet-public-b1"
  }
}

resource "aws_subnet" "public_c1" {
  vpc_id = aws_vpc.main.id
  cidr_block = "10.0.3.0/24"

  availability_zone = "ap-northeast-2c"

  tags = {
    Name = "ttbkk-subnet-public-c1"
  }
}

resource "aws_subnet" "public_d1" {
  vpc_id = aws_vpc.main.id
  cidr_block = "10.0.4.0/24"

  availability_zone = "ap-northeast-2d"

  tags = {
    Name = "ttbkk-subnet-public-d1"
  }
}
