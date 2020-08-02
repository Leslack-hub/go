package Bridge

import (
	"fmt"
)

type Image struct {
}

type MessagerInterface interface {
	Login(username, passWard string)
	SendMessage(message string)
	SendPicture(image Image)
}

type MessagerImp interface {
	PlaySound()
	DrawShape()
	WriteText()
	Connect()
}

type Messager struct {
	MessagerImp MessagerImp
}

/**
 * 普通版本
 */
type MessagerLite struct {
	Messager
}

func (p *MessagerLite) Login(username, passWard string) {
	p.MessagerImp.Connect()
	fmt.Println(username + "登陆了，密码：" + passWard)
}

func (p *MessagerLite) SendMessage(message string) {
	p.MessagerImp.WriteText()
	fmt.Println("发送了信息：" + message)
}

func (p *MessagerLite) SendPicture(image Image) {
	p.MessagerImp.DrawShape()
	fmt.Println("发送了图片。。。")
}

/**
 * 高级版本
 */
type MessagerPrefect struct {
	Messager
}
func (m *MessagerPrefect) Login(username, passWard string) {
	m.MessagerImp.PlaySound()
	m.MessagerImp.Connect()
	fmt.Println("login")
}

func (m *MessagerPrefect) SendMessage(message string) {
	m.MessagerImp.PlaySound()
	m.MessagerImp.WriteText()
	fmt.Println("SendMessage" + message)
}

func (m *MessagerPrefect) SendPicture(image Image) {
	m.MessagerImp.PlaySound()
	m.MessagerImp.DrawShape()
	fmt.Println("SendPicture")
}

// 平台实现
type PCMessageBase struct {
	PcName string
}

func (p *PCMessageBase) PlaySound() {
	fmt.Println(p.PcName + " play sound")
}

func (p *PCMessageBase) DrawShape() {
	fmt.Println(p.PcName + " draw shape")
}

func (p *PCMessageBase) WriteText() {
	fmt.Println(p.PcName + " write text")
}

func (p *PCMessageBase) Connect() {
	fmt.Println(p.PcName + " connect")
}
