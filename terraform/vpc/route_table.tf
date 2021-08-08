resource "aws_route_table" "route_table" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "ttbkk-route-table"
  }
}

resource "aws_route_table_association" "route_table_association_1" {
  subnet_id = aws_subnet.public_a1.id
  route_table_id = aws_route_table.route_table.id
}

resource "aws_route_table_association" "route_table_association_2" {
  subnet_id = aws_subnet.public_b1.id
  route_table_id = aws_route_table.route_table.id
}

resource "aws_route_table_association" "route_table_association_3" {
  subnet_id = aws_subnet.public_c1.id
  route_table_id = aws_route_table.route_table.id
}

resource "aws_route_table_association" "route_table_association_4" {
  subnet_id = aws_subnet.public_d1.id
  route_table_id = aws_route_table.route_table.id
}