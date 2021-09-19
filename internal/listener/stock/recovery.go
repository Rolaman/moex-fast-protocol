package stocklistener

import (
	"fmt"
	decoder "github.com/kdt-wolf/moex-fast/internal/decoder/stock"
	"github.com/kdt-wolf/moex-fast/internal/protocol"
	"log"
	"strings"
)

type Recoverer struct {
	options *protocol.TcpOptions
}

func (r *Recoverer) AddMissed(market string, msgNum uint32, count uint32) []*decoder.XOLRFOND {
	client := protocol.NewTcpClient(r.options)
	logonResp := client.Send(buildLogonMsg())
	logonMsg, success := decoder.Decode(logonResp)
	if !success || logonMsg == nil || logonMsg.Logon == nil {
		log.Printf("Bad logon response: %+v", logonMsg)
		return nil
	}
	missedMsgBytes := client.SendAndReadUntilMessage(receiveUntilLogoutMsg,
		buildMarketDataRequest(market, msgNum, count))
	result := make([]*decoder.XOLRFOND, 0)
	for _, msgByte := range missedMsgBytes {
		message, _ := decoder.Decode(msgByte)
		if message.XOLRFOND != nil {
			result = append(result, message.XOLRFOND)
		} else {
			log.Printf("Unexpected message for stock recovery: %+v", message)
		}
	}
	return result
}

func buildLogonMsg() []byte {
	builder := strings.Builder{}
	builder.WriteString("8=FIXT.1.1")
	builder.WriteByte(1)
	builder.WriteString("35=A")
	builder.WriteByte(1)
	builder.WriteString("56=MOEX")
	builder.WriteByte(1)
	builder.WriteString("553=user1")
	builder.WriteByte(1)
	builder.WriteString("554=user1")
	builder.WriteByte(1)
	builder.WriteString("1137=9")
	builder.WriteByte(1)
	return []byte(builder.String())
}

func buildMarketDataRequest(market string, msgNum uint32, count uint32) []byte {
	builder := strings.Builder{}
	builder.WriteString("8=FIXT.1.1")
	builder.WriteByte(1)
	builder.WriteString("35=A")
	builder.WriteByte(1)
	builder.WriteString("1128=9")
	builder.WriteByte(1)
	builder.WriteString("56=MOEX")
	builder.WriteByte(1)

	builder.WriteString("1180=" + market)
	builder.WriteByte(1)
	builder.WriteString(fmt.Sprintf("1182=%d", msgNum))
	builder.WriteByte(1)
	builder.WriteString(fmt.Sprintf("1183=%d", msgNum+count-1))
	builder.WriteByte(1)

	return []byte(builder.String())
}

func receiveUntilLogoutMsg(b []byte) bool {
	message, _ := decoder.Decode(b)
	if message == nil {
		panic(fmt.Sprintf("Empty message for bytes: %+v", b))
	}
	if message.Logout != nil {
		if message.XOLRFOND != nil {
			return true
		} else {
			log.Printf("Unexpected message for stock recovery: %+v", message)
			return true
		}
	}
	return false
}
