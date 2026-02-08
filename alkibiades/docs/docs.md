# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [v1/alkibiades.proto](#v1_alkibiades-proto)
    - [AggregatedOptions](#alkibiades-v1-AggregatedOptions)
    - [Coordinates](#alkibiades-v1-Coordinates)
    - [CreationRequest](#alkibiades-v1-CreationRequest)
    - [Intro](#alkibiades-v1-Intro)
    - [MatchPair](#alkibiades-v1-MatchPair)
    - [MatchQuiz](#alkibiades-v1-MatchQuiz)
    - [MediaDropQuiz](#alkibiades-v1-MediaDropQuiz)
    - [MediaEntry](#alkibiades-v1-MediaEntry)
    - [QuizResponse](#alkibiades-v1-QuizResponse)
    - [QuizStep](#alkibiades-v1-QuizStep)
    - [Segments](#alkibiades-v1-Segments)
    - [StructureQuiz](#alkibiades-v1-StructureQuiz)
    - [Theme](#alkibiades-v1-Theme)
    - [TranslationStep](#alkibiades-v1-TranslationStep)
    - [TriviaQuiz](#alkibiades-v1-TriviaQuiz)
  
    - [Alkibiades](#alkibiades-v1-Alkibiades)
  
- [Scalar Value Types](#scalar-value-types)



<a name="v1_alkibiades-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## v1/alkibiades.proto



<a name="alkibiades-v1-AggregatedOptions"></a>

### AggregatedOptions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| themes | [Theme](#alkibiades-v1-Theme) | repeated |  |






<a name="alkibiades-v1-Coordinates"></a>

### Coordinates



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| x | [float](#float) |  |  |
| y | [float](#float) |  |  |






<a name="alkibiades-v1-CreationRequest"></a>

### CreationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| theme | [string](#string) |  |  |
| segment | [string](#string) |  | Optional |






<a name="alkibiades-v1-Intro"></a>

### Intro



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| author | [string](#string) |  |  |
| work | [string](#string) |  |  |
| background | [string](#string) |  |  |






<a name="alkibiades-v1-MatchPair"></a>

### MatchPair



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| greek | [string](#string) |  |  |
| answer | [string](#string) |  |  |






<a name="alkibiades-v1-MatchQuiz"></a>

### MatchQuiz



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| instruction | [string](#string) |  |  |
| pairs | [MatchPair](#alkibiades-v1-MatchPair) | repeated |  |






<a name="alkibiades-v1-MediaDropQuiz"></a>

### MediaDropQuiz



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| instruction | [string](#string) |  |  |
| mediaFiles | [MediaEntry](#alkibiades-v1-MediaEntry) | repeated |  |






<a name="alkibiades-v1-MediaEntry"></a>

### MediaEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| word | [string](#string) |  |  |
| answer | [string](#string) |  | image filename or URL |






<a name="alkibiades-v1-QuizResponse"></a>

### QuizResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| theme | [string](#string) |  |  |
| segment | [string](#string) |  |  |
| number | [int32](#int32) |  |  |
| sentence | [string](#string) |  | Full Greek sentence |
| translation | [string](#string) |  | English translation |
| contextNote | [string](#string) |  | Informational text about the passage |
| intro | [Intro](#alkibiades-v1-Intro) |  |  |
| quiz | [QuizStep](#alkibiades-v1-QuizStep) | repeated | Quiz is polymorphic — see below |






<a name="alkibiades-v1-QuizStep"></a>

### QuizStep



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| match | [MatchQuiz](#alkibiades-v1-MatchQuiz) |  |  |
| trivia | [TriviaQuiz](#alkibiades-v1-TriviaQuiz) |  |  |
| structure | [StructureQuiz](#alkibiades-v1-StructureQuiz) |  |  |
| media | [MediaDropQuiz](#alkibiades-v1-MediaDropQuiz) |  |  |
| final_translation | [TranslationStep](#alkibiades-v1-TranslationStep) |  |  |






<a name="alkibiades-v1-Segments"></a>

### Segments



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| number | [int32](#int32) |  |  |
| location | [string](#string) |  |  |
| coordinates | [Coordinates](#alkibiades-v1-Coordinates) |  |  |






<a name="alkibiades-v1-StructureQuiz"></a>

### StructureQuiz



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| title | [string](#string) |  |  |
| text | [string](#string) |  |  |
| question | [string](#string) |  |  |
| options | [string](#string) | repeated |  |
| answer | [string](#string) |  |  |
| note | [string](#string) |  | optional |






<a name="alkibiades-v1-Theme"></a>

### Theme



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| segments | [Segments](#alkibiades-v1-Segments) | repeated |  |






<a name="alkibiades-v1-TranslationStep"></a>

### TranslationStep



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| instruction | [string](#string) |  |  |
| options | [string](#string) | repeated |  |
| answer | [string](#string) |  |  |






<a name="alkibiades-v1-TriviaQuiz"></a>

### TriviaQuiz



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| question | [string](#string) |  |  |
| options | [string](#string) | repeated |  |
| answer | [string](#string) |  |  |
| note | [string](#string) |  | optional explanation |





 

 

 


<a name="alkibiades-v1-Alkibiades"></a>

### Alkibiades


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Health | [.koinos.v1.HealthRequest](#koinos-v1-HealthRequest) | [.koinos.v1.HealthResponse](#koinos-v1-HealthResponse) |  |
| Options | [.koinos.v1.OptionsRequest](#koinos-v1-OptionsRequest) | [AggregatedOptions](#alkibiades-v1-AggregatedOptions) |  |
| Question | [CreationRequest](#alkibiades-v1-CreationRequest) | [QuizResponse](#alkibiades-v1-QuizResponse) |  |

 



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

