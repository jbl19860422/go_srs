/*
The MIT License (MIT)

Copyright (c) 2019 GOSRS(gosrs)

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
	"go_srs/srs/app/config"
	"go_srs/srs/codec"
)

/**
* hls stream cache,
* use to cache hls stream and flush to hls muxer.
*
* when write stream to ts file:
* video frame will directly flush to M3u8Muxer,
* audio frame need to cache, because it's small and flv tbn problem.
*
* whatever, the Hls cache used to cache video/audio,
* and flush video/audio to m3u8 muxer if needed.
*
* about the flv tbn problem:
*   flv tbn is 1/1000, ts tbn is 1/90000,
*   when timestamp convert to flv tbn, it will loose precise,
*   so we must gather audio frame together, and recalc the timestamp @see SrsTsAacJitter,
*   we use a aac jitter to correct the audio pts.
 */

type SrsHlsCache struct {
	cache *SrsTsCache
}

func NewSrsHlsCache() *SrsHlsCache {
	return &SrsHlsCache{
		cache: NewSrsTsCache(),
	}
}

func (this *SrsHlsCache) onPublish(muxer *SrsHlsMuxer, req *SrsRequest, segment_start_dts int64) error {
	//todo vhost
	vhostName := req.vhost
	hlsFragment := config.GetHlsFragment(vhostName)
	hlsWindow := config.GetHlsWindow(vhostName)
	entryPrefix := config.GetHlsEntryPrefix(vhostName)
	m3u8File := config.GetHlsM3u8File(vhostName)
	hlsPath := config.GetHlsPath(vhostName)
	tsFile := config.GetHlsTsFile(vhostName)
	cleanUp := config.GetHlsCleanup(vhostName)
	hlsWaitKeyframe := config.GetHlsWaitKeyframe(vhostName)
	muxer.UpdateConfig(req, entryPrefix, hlsPath, m3u8File, tsFile, float64(hlsFragment), float64(hlsWindow), false, 0.0, cleanUp, hlsWaitKeyframe)

	muxer.segmentOpen(segment_start_dts)
	return nil
}

func (this *SrsHlsCache) onUnpublish(muxer *SrsHlsMuxer) error {
	if err := muxer.flushVideo(this.cache); err != nil {
		return err
	}

	if err := muxer.flushAudio(this.cache); err != nil {
		return err
	}

	if err := muxer.segmentClose(); err != nil {
		return err
	}

	return nil
}

/**
* when get sequence header,
* must write a #EXT-X-DISCONTINUITY to m3u8.
* @see: hls-m3u8-draft-pantos-http-live-streaming-12.txt
* @see: 3.4.11.  EXT-X-DISCONTINUITY
 */
func (this *SrsHlsCache) onSequenceHeader(muxer *SrsHlsMuxer) error {

	return muxer.on_sequence_header()
}

/**
* write audio to cache, if need to flush, flush to muxer.
 */
func (this *SrsHlsCache) writeAudio(c *SrsAvcAacCodec, muxer *SrsHlsMuxer, dts int64, sampler *SrsCodecSampler) error {
	if err := this.cache.cache_audio(c, dts, sampler); err != nil {
		return err
	}

	if err := muxer.flushAudio(this.cache); err != nil {
		return err
	}
	return nil
}

func (this *SrsHlsCache) WriteVideo(c *SrsAvcAacCodec, muxer *SrsHlsMuxer, dts int64, sampler *SrsCodecSampler) error {
	if err := this.cache.cache_video(c, dts, sampler); err != nil {
		return err
	}

	if muxer.is_segment_overflow() {
		if !muxer.hls_wait_keyframe || sampler.FrameType == codec.SrsCodecVideoAVCFrameKeyFrame {
			if err := this.reapSegment("video", muxer, this.cache.video.dts); err != nil {
				return err
			}
		}
	}

	if err := muxer.flushVideo(this.cache); err != nil {
		return err
	}
	return nil
}

func (this *SrsHlsCache) reapSegment(log_desc string, muxer *SrsHlsMuxer, segment_start_dts int64) error {
	if err := muxer.segmentClose(); err != nil {
		return err
	}

	if err := muxer.segmentOpen(segment_start_dts); err != nil {
		return err
	}

	if err := muxer.flushVideo(this.cache); err != nil {
		return err
	}

	if err := muxer.flushAudio(this.cache); err != nil {
		return err
	}
	return nil
}
