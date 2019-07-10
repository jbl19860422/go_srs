/*
The MIT License (MIT)

Copyright (c) 2013-2015 GOSRS(gosrs)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package app

import (
	"go_srs/srs/protocol/rtmp"
)

type SrsDvrConsumer struct {
	source          *SrsSource
	req 			*SrsRequest
	queue           *SrsMessageQueue
	plan			SrsDvrPlan
}

func NewSrsDvrConsumer(s *SrsSource, req *SrsRequest) *SrsDvrConsumer {
	p := NewSrsDvrPlan(req)
	if p == nil {
		return nil
	}
	return &SrsDvrConsumer{
		source:s,
		plan:p,
		queue:NewSrsMessageQueue(),
	}
}

func (this *SrsDvrConsumer) OnPublish() error {
	if this.plan != nil {
		return this.plan.OnPublish()
	}

	this.StopConsume()
	return nil
}

func (this *SrsDvrConsumer) OnUnpublish() error {
	if this.plan != nil {
		return this.plan.OnUnpublish()
	}
	return nil
}

func (this *SrsDvrConsumer) ConsumeCycle() error {
	for {
		msg, err := this.queue.Wait()
		if err != nil {
			return err
		}

		if msg != nil {
			if msg.GetHeader().IsVideo() {
				if err := this.plan.OnVideo(msg); err != nil {
					return err
				}
			} else if msg.GetHeader().IsAudio() {
				if err := this.plan.OnAudio(msg); err != nil {
					return err
				}
			} else {
				if err := this.plan.OnMetaData(msg); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (this *SrsDvrConsumer) StopConsume() error {
	this.source.RemoveConsumer(this)
	//send connection close to response writer
	this.queue.Break()
	return nil
}

func (this *SrsDvrConsumer) OnRecvError(err error) {
	this.StopConsume()
}

func (this *SrsDvrConsumer) Enqueue(msg *rtmp.SrsRtmpMessage, atc bool, jitterAlgorithm *SrsRtmpJitterAlgorithm) {
	this.queue.Enqueue(msg)
}