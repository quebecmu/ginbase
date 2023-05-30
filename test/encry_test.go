package test

import (
	"fmt"
	"github.com/xxandjg/ginbase/pkg"
)

func (s *TestSuite) TestEncrypt() {
	pwd := "muqingcloud.space"
	fmt.Println(pwd)
	hash, g := pkg.PasswordHash(pwd)
	if g.GetCode() != 10000 {
		fmt.Println(g.Error())
	} else {
		//$2a$10$ZJHCn61zJLYeVVkwFl/aRuvfLrQVCCvb6Mrc.9zP.gnzq6xKvEIhq
		fmt.Println(hash)
	}
}

func (s *TestSuite) TestVerify() {
	pwd := "muqingcloud.space"
	hash := "$2a$10$ZJHCn61zJLYeVVkwFl/aRuvfLrQVCCvb6Mrc.9zP.gnzq6xKvEIhq"
	fmt.Println(pwd)
	b, _ := pkg.PasswordVerify(pwd, hash)
	// $2a$10$SOT5p7xSKIgeUP4BpUyvH.U3QSj3B8r5M0XG33J.QsFgWWurerVbe
	fmt.Println(b)
}
