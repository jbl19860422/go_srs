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

package config

type SrsAppSubscriber struct {}
func(this *SrsAppSubscriber) OnReloadUtcTime() {}
func(this *SrsAppSubscriber) OnReloadMaxConns() {}
func(this *SrsAppSubscriber) OnReloadListen() {}
func(this *SrsAppSubscriber) OnReloadPid() {}
func(this *SrsAppSubscriber) OnReloadLogTank() {}
func(this *SrsAppSubscriber) OnReloadLogLevel() {}
func(this *SrsAppSubscriber) OnReloadLogFile() {}
func(this *SrsAppSubscriber) OnReloadPithyPrint() {}
func(this *SrsAppSubscriber) OnReloadHttpApiEnabled() {}
func(this *SrsAppSubscriber) OnReloadHttpApiDisabled() {}
func(this *SrsAppSubscriber) OnReloadHttpStreamEnabled() {}
func(this *SrsAppSubscriber) OnReloadHttpStreamDisabled() {}
func(this *SrsAppSubscriber) OnReloadHttpStreamUpdated() {}

func(this *SrsAppSubscriber) OnReloadVHostHttpUpdated() {}
func(this *SrsAppSubscriber) OnReloadVHostHttpRemuxUpdated(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostAdded(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostRemoved(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostAtc(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostGopCache(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostQueueLength(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostTimeJitter(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostMixCorrect(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostForward(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostHls(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostHds(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostDvr(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostMr(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostMw(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostSmi(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostTcpNodelay(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostRealtime(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostP1stpt(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostPnt(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostChunkSize(vhost string) {}
func(this *SrsAppSubscriber) OnReloadVHostTranscode(vhost string) {}
func(this *SrsAppSubscriber) OnReloadIngestRemoved(vhost string, ingest_id string) {}
func(this *SrsAppSubscriber) OnReloadIngestAdded(vhost string, ingest_id string) {}
func(this *SrsAppSubscriber) OnReloadIngestUpdated(vhost string, ingest_id string) {}
func(this *SrsAppSubscriber) OnReloadUserInfo() {}
