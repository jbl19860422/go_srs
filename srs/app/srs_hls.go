package app

type SrsHls struct {
	source 	*SrsSource
}

func (this *SrsHls) Initialize(s *SrsSource, r *SrsRequest) error {
	this.source = s
	this.plan = flvcodec.NewSrsDvrPlan("./record.flv")
	//todo fix 
	this.plan.Initialize()
	return nil
}

func (this *SrsHls) on_meta_data(metaData *rtmp.SrsRtmpMessage) error {
	return this.plan.On_meta_data(metaData)
}

func (this *SrsHls) on_video(video *rtmp.SrsRtmpMessage) error {
	return this.plan.On_video(video)
}

func (this *SrsHls) on_audio(audio *rtmp.SrsRtmpMessage) error {
	return this.plan.On_audio(audio)
}

func (this *SrsHls) Close() {
	this.plan.Close()
}