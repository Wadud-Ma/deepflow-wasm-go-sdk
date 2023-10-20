package main

import (
	"github.com/deepflowio/deepflow-wasm-go-sdk/sdk"
)

func main() {
	sdk.Info("on flow log demo wasm plugin init")
	sdk.SetParser(SomeParser{})
}

type SomeParser struct {
}

func (p SomeParser) HookIn() []sdk.HookBitmap {
	return []sdk.HookBitmap{
		// 一般只需要 hook 协议解析
		sdk.HOOK_POINT_HTTP_REQ,
		sdk.HOOK_POINT_HTTP_RESP,
		sdk.HOOK_POINT_PAYLOAD_PARSE,
	}
}

func (p SomeParser) OnHttpReq(ctx *sdk.HttpReqCtx) sdk.Action {
	sdk.Info("flow log demo wasm plugin handle OnHttpReq")
	return sdk.ActionNext()
}

func (p SomeParser) OnHttpResp(ctx *sdk.HttpRespCtx) sdk.Action {
	sdk.Info("flow log demo wasm plugin handle OnHttpResp")
	return sdk.ActionNext()
}

func (p SomeParser) OnCheckPayload(ctx *sdk.ParseCtx) (uint8, string) {
	// 这里是协议判断的逻辑， 返回 0 表示失败
	// return 0, ""
	sdk.Info("flow log demo wasm plugin handle OnCheckPayload")
	return 0, "custom wasm protocol"
}

func (p SomeParser) OnParsePayload(ctx *sdk.ParseCtx) sdk.Action {
	sdk.Info("flow log demo wasm plugin handle OnParsePayload")
	// 这里是解析协议的逻辑
	if ctx.L4 != sdk.TCP || ctx.L7 != 1 {
		return sdk.ActionNext()
	}
	return sdk.ActionNext()
}
