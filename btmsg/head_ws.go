package btmsg

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type MsgHeadWs struct {
	Act  uint16
	Size uint32
	body []byte
}

var _ IHead = (*MsgHeadWs)(nil)

func NewMsgHeadWs() *MsgHeadWs {
	return &MsgHeadWs{}
}

func (l *MsgHeadWs) GetAct() uint16 {
	return l.Act
}

func (l *MsgHeadWs) HeadSize() uint32 {
	return 0
}

func (l *MsgHeadWs) BodySize() uint32 {
	return l.Size
}

func (l *MsgHeadWs) Read(r IReader) (err error) {
	var hdMap= map[string]any{
		"act":0,
	}

	defer func() {
		if err != nil {
			err = errors.Wrap(err, "MsgHeadWs.Read")
		}
	}()

	_, l.body, err = r.ReadMessage()
	if err != nil {
		err = errors.Wrap(err, "readAll")
		return
	}

	err = json.Unmarshal(l.body, &hdMap)
	if err != nil {
		err = errors.Wrap(err, "json unmarshall")
		return
	}

	i,ok := hdMap["act"].(float64);
	if !ok{
		err = errors.Errorf("act must float64:%v",hdMap["act"])
		return
	}

	l.Act = uint16(i)
	return nil
}

// todo 这里head 和 body不能都read,所以先把body存起来，后面优化一下，防止误用
func (l *MsgHeadWs) ReadBody(r IReader) (err error, bt []byte) {
	return nil, l.body
}

func (l *MsgHeadWs) ToBytes() []byte {
	return nil
}

func (l *MsgHeadWs) SetSize(size uint32) {
	l.Size = size
}
