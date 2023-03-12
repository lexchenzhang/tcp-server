package main

import (
	"fmt"
	__ "tcp-server/src/demo/protobufDemo/pb"

	"github.com/golang/protobuf/proto"
)

func main() {
	person := &__.Person{
		Name:   "hi",
		Age:    15,
		Emails: []string{"a.b@gmail.com", "a.b@hotmail.com"},
		Phones: []*__.PhoneNumber{
			&__.PhoneNumber{
				Number: "123",
				Type:   __.PhoneType_HOME,
			},
			&__.PhoneNumber{
				Number: "321",
				Type:   __.PhoneType_MOBILE,
			},
		},
	}

	// serialize msg
	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	newData := &__.Person{}
	proto.Unmarshal(data, newData)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	fmt.Println(person)
	fmt.Println(newData)
}
