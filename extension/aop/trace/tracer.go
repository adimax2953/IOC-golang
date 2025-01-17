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

package trace

import (
	"bytes"
	"io"

	"github.com/jaegertracing/jaeger/model"

	"github.com/alibaba/ioc-golang/extension/aop/trace/transport"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

var tracer *wrapperTracer
var collectorAddress = ""
var appName = "ioc-golang-application"

// FIXME: invocation longer than 5s would not be collected
var collectTraceInterval = 5

type wrapperTracer struct {
	rawTracer opentracing.Tracer

	out                  chan []*model.Trace
	subscribingTraceChan chan []*model.Trace

	batchBufferOut             chan *bytes.Buffer
	subscribingBatchBufferChan chan *bytes.Buffer
}

func (w *wrapperTracer) getRawTracer() opentracing.Tracer {
	return w.rawTracer
}

func (w *wrapperTracer) subscribeTrace(subscribingTraceChan chan []*model.Trace) {
	w.subscribingTraceChan = subscribingTraceChan
}

func (w *wrapperTracer) removeSubscribingTrace() {
	w.subscribingTraceChan = nil
}

func (w *wrapperTracer) subscribeBatchBuffer(subscribingBatchBufferChan chan *bytes.Buffer) {
	w.subscribingBatchBufferChan = subscribingBatchBufferChan
}

func (w *wrapperTracer) removeSubscribingBatchBuffer() {
	w.subscribingBatchBufferChan = nil
}

func (w *wrapperTracer) runCollectingTrace() {
	for {
		select {
		case traces := <-w.out:
			if ch := w.subscribingTraceChan; ch != nil {
				select {
				case ch <- traces:
				default:
				}
			}
		case bathBuffer := <-w.batchBufferOut:
			if ch := w.subscribingBatchBufferChan; ch != nil {
				select {
				case ch <- bathBuffer:
				default:
				}
			}
		}
	}
}

func getGlobalTracer() *wrapperTracer {
	if tracer == nil {
		outCh := make(chan []*model.Trace)
		batchBufferOut := make(chan *bytes.Buffer)
		rawJaegerTracer, _ := newJaegerTracer(appName, collectorAddress, outCh, batchBufferOut)
		tracer = &wrapperTracer{
			rawTracer:      rawJaegerTracer,
			batchBufferOut: batchBufferOut,
			out:            outCh,
		}
		go tracer.runCollectingTrace()
	}
	return tracer
}

func getCollectorAddress() string {
	return collectorAddress
}

func setCollectorAddress(addr string) {
	collectorAddress = addr
}

func setAppName(name string) {
	appName = name
}

func newJaegerTracer(service string, collectorAddress string, out chan []*model.Trace, batchBufferOut chan *bytes.Buffer) (opentracing.Tracer, io.Closer) {
	return jaeger.NewTracer(
		service,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(transport.GetLocalWrappedHTTPTransportSingleton(appName, collectorAddress, out, batchBufferOut, collectTraceInterval)),
	)
}
