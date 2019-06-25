package app

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

}

/**
* when get sequence header, 
* must write a #EXT-X-DISCONTINUITY to m3u8.
* @see: hls-m3u8-draft-pantos-http-live-streaming-12.txt
* @see: 3.4.11.  EXT-X-DISCONTINUITY
*/
func (this *SrsHlsCache) on_sequence_header(muxer *SrsHlsMuxer) error {

}

 /**
* write audio to cache, if need to flush, flush to muxer.
*/
func (this *SrsHlsCache) write_audio(codec *SrsAvcAacCodec, muxer *SrsHlsMuxer, pts int64, sample *SrsCodecSample) error {

}