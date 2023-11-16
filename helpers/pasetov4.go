package helpers

import (
	"fmt"

	"github.com/whatsauth/watoken"
)

func pasetov4() {
	privateKey, publicKey := watoken.GenerateKey()

	//generate token for user awangga
	userid := "raul"
	tokenstring, err := watoken.Encode(userid, privateKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tokenstring)
	//decode token to get userid
	useridstring := watoken.DecodeGetId(publicKey, tokenstring)
	if useridstring == "" {
		fmt.Println("expire token")
	}
	fmt.Println(useridstring)
}
