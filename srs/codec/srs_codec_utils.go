package codec

func SrsCodecAacRtmp2Ts(objectType SrsAacObjectType) SrsAacProfile  {
	switch objectType {
		case SrsAacObjectTypeAacMain:
			return SrsAacProfileMain
		case SrsAacObjectTypeAacHE, SrsAacObjectTypeAacHEV2, SrsAacObjectTypeAacLC:
			return SrsAacProfileLC
		case SrsAacObjectTypeAacSSR:
			return SrsAacProfileSSR
		default:
			return SrsAacProfileReserved
	}
}