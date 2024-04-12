package btmsg

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
)

type MsgHeadWs struct {
	Act  uint16
	Size uint32
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

func (l *MsgHeadWs) Read(r io.Reader) (err error) {
	var hdMap= map[string]uint16{
		"act":0,
	}
	var all []byte

	defer func() {
		if err != nil {
			err = errors.Wrap(err, "MsgHeadWs.Read")
		}
	}()

	all, err = io.ReadAll(r)
	if err != nil {
		err = errors.Wrap(err, "readAll")
		return
	}

	err = json.Unmarshal(all, &hdMap)
	if err != nil {
		err = errors.Wrap(err, "json unmarshall")
		return
	}

	l.Act = hdMap["act"]

	return nil
}

func (l *MsgHeadWs) ReadBody(r io.Reader) (err error, bt []byte) {
	bt, err = io.ReadAll(r)
	if err != nil {
		err = errors.Wrap(err, "msgHeadWs.ReadBody readAll")
		return
	}

	return
}

func (l *MsgHeadWs) ToBytes() []byte {
	return nil
}

func (l *MsgHeadWs) SetSize(size uint32) {
	l.Size = size
}
