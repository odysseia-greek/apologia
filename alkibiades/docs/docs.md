# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [kriton.proto](#kriton-proto)
    - [AggregatedOptions](#apologia_kriton-AggregatedOptions)
    - [AnswerRequest](#apologia_kriton-AnswerRequest)
    - [AnswerResponse](#apologia_kriton-AnswerResponse)
    - [CreationRequest](#apologia_kriton-CreationRequest)
    - [DatabaseHealth](#apologia_kriton-DatabaseHealth)
    - [Dialogue](#apologia_kriton-Dialogue)
    - [DialogueContent](#apologia_kriton-DialogueContent)
    - [DialogueCorrection](#apologia_kriton-DialogueCorrection)
    - [HealthRequest](#apologia_kriton-HealthRequest)
    - [HealthResponse](#apologia_kriton-HealthResponse)
    - [OptionsRequest](#apologia_kriton-OptionsRequest)
    - [QuizMetadata](#apologia_kriton-QuizMetadata)
    - [QuizResponse](#apologia_kriton-QuizResponse)
    - [Segment](#apologia_kriton-Segment)
    - [Speaker](#apologia_kriton-Speaker)
    - [Theme](#apologia_kriton-Theme)
  
    - [Kriton](#apologia_kriton-Kriton)
  
- [Scalar Value Types](#scalar-value-types)



<a name="kriton-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## kriton.proto



<a name="apologia_kriton-AggregatedOptions"></a>

### AggregatedOptions
Response message for quiz options


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| themes | [Theme](#apologia_kriton-Theme) | repeated |  |






<a name="apologia_kriton-AnswerRequest"></a>

### AnswerRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| theme | [string](#string) | optional |  |
| set | [int32](#int32) | optional |  |
| content | [DialogueContent](#apologia_kriton-DialogueContent) | repeated |  |






<a name="apologia_kriton-AnswerResponse"></a>

### AnswerResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| percentage | [double](#double) |  |  |
| input | [DialogueContent](#apologia_kriton-DialogueContent) | repeated |  |
| answer | [DialogueContent](#apologia_kriton-DialogueContent) | repeated |  |
| wronglyPlaced | [DialogueCorrection](#apologia_kriton-DialogueCorrection) | repeated |  |






<a name="apologia_kriton-CreationRequest"></a>

### CreationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| theme | [string](#string) |  |  |
| set | [string](#string) |  |  |






<a name="apologia_kriton-DatabaseHealth"></a>

### DatabaseHealth
Nested message for database health details


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| healthy | [bool](#bool) |  |  |
| cluster_name | [string](#string) |  |  |
| server_name | [string](#string) |  |  |
| server_version | [string](#string) |  |  |






<a name="apologia_kriton-Dialogue"></a>

### Dialogue



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| introduction | [string](#string) |  |  |
| speakers | [Speaker](#apologia_kriton-Speaker) | repeated |  |
| section | [string](#string) |  |  |
| linkToPerseus | [string](#string) |  |  |






<a name="apologia_kriton-DialogueContent"></a>

### DialogueContent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| translation | [string](#string) |  |  |
| greek | [string](#string) | optional |  |
| place | [int32](#int32) | optional |  |
| speaker | [string](#string) | optional |  |






<a name="apologia_kriton-DialogueCorrection"></a>

### DialogueCorrection



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| translation | [string](#string) |  |  |
| greek | [string](#string) | optional |  |
| place | [int32](#int32) | optional |  |
| speaker | [string](#string) | optional |  |
| correctPlace | [int32](#int32) | optional |  |






<a name="apologia_kriton-HealthRequest"></a>

### HealthRequest
Empty request messages since these endpoints require no body






<a name="apologia_kriton-HealthResponse"></a>

### HealthResponse
Response message for health check


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| healthy | [bool](#bool) |  |  |
| time | [string](#string) |  |  |
| version | [string](#string) |  |  |
| database_health | [DatabaseHealth](#apologia_kriton-DatabaseHealth) |  |  |






<a name="apologia_kriton-OptionsRequest"></a>

### OptionsRequest







<a name="apologia_kriton-QuizMetadata"></a>

### QuizMetadata



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| language | [string](#string) |  |  |






<a name="apologia_kriton-QuizResponse"></a>

### QuizResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| quizMetadata | [QuizMetadata](#apologia_kriton-QuizMetadata) |  |  |
| theme | [string](#string) | optional |  |
| set | [int32](#int32) | optional |  |
| segment | [string](#string) | optional |  |
| reference | [string](#string) | optional |  |
| dialogue | [Dialogue](#apologia_kriton-Dialogue) | optional |  |
| content | [DialogueContent](#apologia_kriton-DialogueContent) | repeated |  |






<a name="apologia_kriton-Segment"></a>

### Segment
Structure for segments within a theme


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| max_set | [float](#float) |  |  |






<a name="apologia_kriton-Speaker"></a>

### Speaker



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| shorthand | [string](#string) |  |  |
| translation | [string](#string) |  |  |






<a name="apologia_kriton-Theme"></a>

### Theme
Structure for quiz themes


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| segments | [Segment](#apologia_kriton-Segment) | repeated |  |





 

 

 


<a name="apologia_kriton-Kriton"></a>

### Kriton


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Health | [HealthRequest](#apologia_kriton-HealthRequest) | [HealthResponse](#apologia_kriton-HealthResponse) |  |
| Options | [OptionsRequest](#apologia_kriton-OptionsRequest) | [AggregatedOptions](#apologia_kriton-AggregatedOptions) |  |
| Question | [CreationRequest](#apologia_kriton-CreationRequest) | [QuizResponse](#apologia_kriton-QuizResponse) |  |
| Answer | [AnswerRequest](#apologia_kriton-AnswerRequest) | [AnswerResponse](#apologia_kriton-AnswerResponse) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

