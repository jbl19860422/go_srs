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
package amf0

const (
	RTMP_AMF0_Number      = 0x00
	RTMP_AMF0_Boolean     = 0x01
	RTMP_AMF0_String      = 0x02
	RTMP_AMF0_Object      = 0x03
	RTMP_AMF0_MovieClip   = 0x04 // reserved, not supported
	RTMP_AMF0_Null        = 0x05
	RTMP_AMF0_Undefined   = 0x06
	RTMP_AMF0_Reference   = 0x07
	RTMP_AMF0_EcmaArray   = 0x08
	RTMP_AMF0_ObjectEnd   = 0x09
	RTMP_AMF0_StrictArray = 0x0A
	RTMP_AMF0_Date        = 0x0B
	RTMP_AMF0_LongString  = 0x0C
	RTMP_AMF0_UnSupported = 0x0D
	RTMP_AMF0_RecordSet   = 0x0E
	RTMP_AMF0_XmlDocument = 0x0F
	RTMP_AMF0_TypedObject = 0x10
	// AVM+ object is the AMF3 object.
	RTMP_AMF0_AVMplusObject = 0x11
	// origin array whos data takes the same form as LengthValueBytes
	RTMP_AMF0_OriginStrictArray = 0x20
	// User defined
	RTMP_AMF0_Invalid = 0x3F
)

/**
 * amf0 command message, command name macros
 */
const (
	RTMP_AMF0_COMMAND_CONNECT        = "connect"
	RTMP_AMF0_COMMAND_CREATE_STREAM  = "createStream"
	RTMP_AMF0_COMMAND_CLOSE_STREAM   = "closeStream"
	RTMP_AMF0_COMMAND_PLAY           = "play"
	RTMP_AMF0_COMMAND_PAUSE          = "pause"
	RTMP_AMF0_COMMAND_ON_BW_DONE     = "onBWDone"
	RTMP_AMF0_COMMAND_ON_STATUS      = "onStatus"
	RTMP_AMF0_COMMAND_RESULT         = "_result"
	RTMP_AMF0_COMMAND_ERROR          = "_error"
	RTMP_AMF0_COMMAND_RELEASE_STREAM = "releaseStream"
	RTMP_AMF0_COMMAND_FC_PUBLISH     = "FCPublish"
	RTMP_AMF0_COMMAND_UNPUBLISH      = "FCUnpublish"
	RTMP_AMF0_COMMAND_PUBLISH        = "publish"
	RTMP_AMF0_DATA_SAMPLE_ACCESS     = "|RtmpSampleAccess"
)

type SrsValuePair struct {
	Name  SrsAmf0Utf8
	Value SrsAmf0Any
}

const SRS_CONSTS_RTMP_SET_DATAFRAME  = "@setDataFrame"
const SRS_CONSTS_RTMP_ON_METADATA = "onMetaData"