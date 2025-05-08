# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [alkibiades.proto](#alkibiades-proto)
    - [AggregatedOptions](#apologia_alkibiades-AggregatedOptions)
    - [Coordinates](#apologia_alkibiades-Coordinates)
    - [CreationRequest](#apologia_alkibiades-CreationRequest)
    - [DatabaseHealth](#apologia_alkibiades-DatabaseHealth)
    - [HealthRequest](#apologia_alkibiades-HealthRequest)
    - [HealthResponse](#apologia_alkibiades-HealthResponse)
    - [Intro](#apologia_alkibiades-Intro)
    - [MatchPair](#apologia_alkibiades-MatchPair)
    - [MatchQuiz](#apologia_alkibiades-MatchQuiz)
    - [MediaDropQuiz](#apologia_alkibiades-MediaDropQuiz)
    - [MediaEntry](#apologia_alkibiades-MediaEntry)
    - [OptionsRequest](#apologia_alkibiades-OptionsRequest)
    - [QuizResponse](#apologia_alkibiades-QuizResponse)
    - [QuizStep](#apologia_alkibiades-QuizStep)
    - [Segments](#apologia_alkibiades-Segments)
    - [StructureQuiz](#apologia_alkibiades-StructureQuiz)
    - [Theme](#apologia_alkibiades-Theme)
    - [TranslationStep](#apologia_alkibiades-TranslationStep)
    - [TriviaQuiz](#apologia_alkibiades-TriviaQuiz)
  
    - [Alkibiades](#apologia_alkibiades-Alkibiades)
  
- [Scalar Value Types](#scalar-value-types)



<a name="alkibiades-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## alkibiades.proto



<a name="apologia_alkibiades-AggregatedOptions"></a>

### AggregatedOptions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| themes | [Theme](#apologia_alkibiades-Theme) | repeated |  |






<a name="apologia_alkibiades-Coordinates"></a>

### Coordinates



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| x | [float](#float) |  |  |
| y | [float](#float) |  |  |






<a name="apologia_alkibiades-CreationRequest"></a>

### CreationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| theme | [string](#string) |  |  |
| segment | [string](#string) |  | Optional |






<a name="apologia_alkibiades-DatabaseHealth"></a>

### DatabaseHealth
Nested message for database health details


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| healthy | [bool](#bool) |  |  |
| cluster_name | [string](#string) |  |  |
| server_name | [string](#string) |  |  |
| server_version | [string](#string) |  |  |






<a name="apologia_alkibiades-HealthRequest"></a>

### HealthRequest
Empty request messages since these endpoints require no body






<a name="apologia_alkibiades-HealthResponse"></a>

### HealthResponse
Response message for health check


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| healthy | [bool](#bool) |  |  |
| time | [string](#string) |  |  |
| version | [string](#string) |  |  |
| database_health | [DatabaseHealth](#apologia_alkibiades-DatabaseHealth) |  |  |






<a name="apologia_alkibiades-Intro"></a>

### Intro



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| author | [string](#string) |  |  |
| work | [string](#string) |  |  |
| background | [string](#string) |  |  |






<a name="apologia_alkibiades-MatchPair"></a>

### MatchPair



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| greek | [string](#string) |  |  |
| answer | [string](#string) |  |  |






<a name="apologia_alkibiades-MatchQuiz"></a>

### MatchQuiz



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| instruction | [string](#string) |  |  |
| pairs | [MatchPair](#apologia_alkibiades-MatchPair) | repeated |  |






<a name="apologia_alkibiades-MediaDropQuiz"></a>

### MediaDropQuiz



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| instruction | [string](#string) |  |  |
| mediaFiles | [MediaEntry](#apologia_alkibiades-MediaEntry) | repeated |  |






<a name="apologia_alkibiades-MediaEntry"></a>

### MediaEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| word | [string](#string) |  |  |
| answer | [string](#string) |  | image filename or URL |






<a name="apologia_alkibiades-OptionsRequest"></a>

### OptionsRequest







<a name="apologia_alkibiades-QuizResponse"></a>

### QuizResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| theme | [string](#string) |  |  |
| segment | [string](#string) |  |  |
| number | [int32](#int32) |  |  |
| sentence | [string](#string) |  | Full Greek sentence |
| translation | [string](#string) |  | English translation |
| contextNote | [string](#string) |  | Informational text about the passage |
| intro | [Intro](#apologia_alkibiades-Intro) |  |  |
| quiz | [QuizStep](#apologia_alkibiades-QuizStep) | repeated | Quiz is polymorphic — see below |






<a name="apologia_alkibiades-QuizStep"></a>

### QuizStep



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| match | [MatchQuiz](#apologia_alkibiades-MatchQuiz) |  |  |
| trivia | [TriviaQuiz](#apologia_alkibiades-TriviaQuiz) |  |  |
| structure | [StructureQuiz](#apologia_alkibiades-StructureQuiz) |  |  |
| media | [MediaDropQuiz](#apologia_alkibiades-MediaDropQuiz) |  |  |
| final_translation | [TranslationStep](#apologia_alkibiades-TranslationStep) |  |  |






<a name="apologia_alkibiades-Segments"></a>

### Segments



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| number | [int32](#int32) |  |  |
| location | [string](#string) |  |  |
| coordinates | [Coordinates](#apologia_alkibiades-Coordinates) |  |  |






<a name="apologia_alkibiades-StructureQuiz"></a>

### StructureQuiz



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| title | [string](#string) |  |  |
| text | [string](#string) |  |  |
| question | [string](#string) |  |  |
| options | [string](#string) | repeated |  |
| answer | [string](#string) |  |  |
| note | [string](#string) |  | optional |






<a name="apologia_alkibiades-Theme"></a>

### Theme



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| segments | [Segments](#apologia_alkibiades-Segments) | repeated |  |






<a name="apologia_alkibiades-TranslationStep"></a>

### TranslationStep



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| instruction | [string](#string) |  |  |
| options | [string](#string) | repeated |  |
| answer | [string](#string) |  |  |






<a name="apologia_alkibiades-TriviaQuiz"></a>

### TriviaQuiz



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| question | [string](#string) |  |  |
| options | [string](#string) | repeated |  |
| answer | [string](#string) |  |  |
| note | [string](#string) |  | optional explanation |





 

 

 


<a name="apologia_alkibiades-Alkibiades"></a>

### Alkibiades


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Health | [HealthRequest](#apologia_alkibiades-HealthRequest) | [HealthResponse](#apologia_alkibiades-HealthResponse) |  |
| Options | [OptionsRequest](#apologia_alkibiades-OptionsRequest) | [AggregatedOptions](#apologia_alkibiades-AggregatedOptions) |  |
| Question | [CreationRequest](#apologia_alkibiades-CreationRequest) | [QuizResponse](#apologia_alkibiades-QuizResponse) |  |

 



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

