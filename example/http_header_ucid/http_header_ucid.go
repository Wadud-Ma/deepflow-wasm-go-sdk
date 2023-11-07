/*
 * Copyright (c) 2022 Yunshan Networks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bufio"
	"bytes"
	"github.com/deepflowio/deepflow-wasm-go-sdk/sdk"
	"net/http"
	"net/url"
	"strings"
)

const (
	ID_MIX1_LEN = 3
	ID_MIX2_LEN = 5
)

type httpHook struct {
}

func (p httpHook) HookIn() []sdk.HookBitmap {
	return []sdk.HookBitmap{
		sdk.HOOK_POINT_HTTP_REQ,
		sdk.HOOK_POINT_HTTP_RESP,
	}
}

/*
assume the http request as follow:

	GET /user_info?username=test&type=1 HTTP/1.1
	Custom-Trace-Info: trace_id: xxx, span_id: sss
*/
func (p httpHook) OnHttpReq(ctx *sdk.HttpReqCtx) sdk.Action {
	baseCtx := &ctx.BaseCtx
	sdk.Info("========= HttpReqCtx: %+v ", ctx)
	if !strings.HasPrefix(ctx.Path, "/web/fe/helpcenter/moduleConfigServiceImpl/getHomePageConfig?") {
		return sdk.ActionNext()
	}
	payload, err := baseCtx.GetPayload()
	if err != nil {
		return sdk.ActionAbortWithErr(err)
	}
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(payload)))
	if err != nil {
		sdk.Info("========= ReadRequest Error: %+v", payload)
		return sdk.ActionAbortWithErr(err)
	}
	sdk.Info("========= Request: %+v", req)

	attr := []sdk.KeyVal{}
	cookies := req.Cookies()
	sdk.Info("========= Cookies: %+v", cookies)
	for _, cookie := range cookies {
		if cookie.Name == "login_ucid" {
			if cookie.Value != "" {
				sdk.Info("========= login_ucid : %s", cookie.Value)
				attr = append(attr, sdk.KeyVal{
					Key: "login_ucid",
					Val: cookie.Value,
				})
			}
		} else if cookie.Name == "_lianjia_link_snid" {
			// URL 解码
			data, err := url.QueryUnescape(cookie.Value)
			if err != nil {
				return sdk.ActionAbortWithErr(err)
			}
			// 根据 "|" 分割字符串
			keyMixs := strings.Split(data, "\\|")
			for _, keyMix := range keyMixs {
				// 根据 ":" 分割字符串
				ucIdMix := strings.Split(keyMix, ":")
				if len(ucIdMix) != ID_MIX1_LEN && len(ucIdMix) != ID_MIX2_LEN {
					continue
				}
				if ucIdMix[0] != "" {
					sdk.Info("========= login_ucid : %s", ucIdMix[0])
					attr = append(attr, sdk.KeyVal{
						Key: "login_ucid",
						Val: ucIdMix[0],
					})
				}
			}
		}
	}

	return sdk.HttpReqActionAbortWithResult(nil, nil, attr)
}

/*
assume resp as follow:

	HTTP/1.1 200 OK

	{"code": 0, "data": {"user_id": 12345, "register_time": 1682050409}}
*/
func (p httpHook) OnHttpResp(ctx *sdk.HttpRespCtx) sdk.Action {
	return sdk.ActionNext()
}

func (p httpHook) OnCheckPayload(baseCtx *sdk.ParseCtx) (uint8, string) {
	return 0, ""
}

func (p httpHook) OnParsePayload(baseCtx *sdk.ParseCtx) sdk.Action {
	return sdk.ActionNext()
}

func main() {
	sdk.Warn("wasm register http hook")
	sdk.SetParser(httpHook{})

}
