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

//type SrsDvr struct {
//	source 	*SrsSource
//	plan	*SrsDvrPlan
//}
//
//func NewSrsDvr() *SrsDvr {
//	return &SrsDvr{}
//}
//
//func (this *SrsDvr) Initialize(s *SrsSource, r *SrsRequest) error {
//	this.source = s
//	this.plan = NewSrsDvrPlan("./record.flv")
//	//todo fix
//	this.plan.Initialize()
//	return nil
//}
//
//func (this *SrsDvr) OnMetaData(metaData *rtmp.SrsRtmpMessage) error {
//	return this.plan.OnMetaData(metaData)
//}
//
//func (this *SrsDvr) on_video(video *rtmp.SrsRtmpMessage) error {
//	return this.plan.On_video(video)
//}
//
//func (this *SrsDvr) on_audio(audio *rtmp.SrsRtmpMessage) error {
//	return this.plan.On_audio(audio)
//}
//
//func (this *SrsDvr) Close() {
//	this.plan.Close()
//}