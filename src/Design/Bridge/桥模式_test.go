package Bridge

import "testing"

func Test(t *testing.T) {
	messageBase := &PCMessageBase{"pc"}
	m := Messager{MessagerImp: messageBase}
	messager := MessagerLite{m}
	messager.Login("nihao", "nihao")
	messager2 := MessagerPrefect{m}
	messager2.SendMessage("prefect 2222")
}
