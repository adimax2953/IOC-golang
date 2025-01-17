/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package transaction

import (
	"reflect"

	"github.com/alibaba/ioc-golang/aop"

	"github.com/alibaba/ioc-golang/autowire"
)

type context struct {
	rollbackAbleInvocationContexts []rollbackAbleInvocationCtx
	entranceMethod                 string
}

func newContext(entranceMethod string) *context {
	return &context{
		rollbackAbleInvocationContexts: make([]rollbackAbleInvocationCtx, 0),
		entranceMethod:                 entranceMethod,
	}
}

func (c *context) finish() {

}

func (c *context) failed(err error) {
	for i := len(c.rollbackAbleInvocationContexts) - 1; i >= 0; i-- {
		snapshot := c.rollbackAbleInvocationContexts[i]
		snapshot.rollback(err)
	}
}

func (c *context) addSuccessfullyCalledInvocationCtx(ctx *aop.InvocationContext) {
	sd := autowire.GetStructDescriptor(ctx.SDID)
	if sd == nil {
		// todo: print logs
		return
	}
	if rollbackMethodName, ok := sd.TransactionMethodsMap[ctx.MethodName]; ok && rollbackMethodName != "" {
		c.rollbackAbleInvocationContexts = append(c.rollbackAbleInvocationContexts, rollbackAbleInvocationCtx{
			invocationCtx:      ctx,
			rollbackMethodName: rollbackMethodName,
		})
	}
}

type rollbackAbleInvocationCtx struct {
	invocationCtx      *aop.InvocationContext
	rollbackMethodName string
}

func (c *rollbackAbleInvocationCtx) rollback(err error) {
	valueOf := reflect.ValueOf(c.invocationCtx.ProxyServicePtr)
	valueOfElem := valueOf.Elem()
	// todo what if rollback function annotation is incorrect? it would cause reflect.Value.Call on zero Value
	funcRaw := valueOfElem.FieldByName(c.rollbackMethodName + "_")
	rollbackParam := c.invocationCtx.Params
	rollbackParam = append(rollbackParam, reflect.ValueOf(err.Error()))
	funcRaw.Call(rollbackParam)
}
