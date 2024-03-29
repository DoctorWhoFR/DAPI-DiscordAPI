<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# DiscordInternal

```go
import "azginfr/dapi/DiscordInternal"
```

## Index

- [Constants](<#constants>)
- [Variables](<#variables>)
- [func GetEnvLogLevel\(\)](<#GetEnvLogLevel>)
- [func LogDebug\(args ...any\)](<#LogDebug>)
- [func LogError\(args ...any\)](<#LogError>)
- [func LogInfo\(args ...any\)](<#LogInfo>)
- [func LogLog\(args ...any\)](<#LogLog>)
- [func LogTrace\(args ...any\)](<#LogTrace>)
- [func MediumValueAPI\(\) int64](<#MediumValueAPI>)
- [func MediumValueBOT\(\) int64](<#MediumValueBOT>)
- [func MediumValueWithoutAPI\(\) int64](<#MediumValueWithoutAPI>)
- [func SimpleSleep\(\)](<#SimpleSleep>)


## Constants

<a name="BaseSleepTime"></a>

```go
const BaseSleepTime = time.Millisecond * 10
```

## Variables

<a name="HandlingTimeBot"></a>

```go
var HandlingTimeBot = []time.Duration{0}
```

<a name="HandlingTimeDiscord"></a>

```go
var HandlingTimeDiscord = []time.Duration{0}
```

<a name="GetEnvLogLevel"></a>
## func [GetEnvLogLevel](<https://github.com/DoctorWhoFR/dapi/blob/master/DiscordInternal/logging.go#L27>)

```go
func GetEnvLogLevel()
```



<a name="LogDebug"></a>
## func [LogDebug](<https://github.com/DoctorWhoFR/dapi/blob/master/DiscordInternal/logging.go#L67>)

```go
func LogDebug(args ...any)
```



<a name="LogError"></a>
## func [LogError](<https://github.com/DoctorWhoFR/dapi/blob/master/DiscordInternal/logging.go#L81>)

```go
func LogError(args ...any)
```



<a name="LogInfo"></a>
## func [LogInfo](<https://github.com/DoctorWhoFR/dapi/blob/master/DiscordInternal/logging.go#L74>)

```go
func LogInfo(args ...any)
```



<a name="LogLog"></a>
## func [LogLog](<https://github.com/DoctorWhoFR/dapi/blob/master/DiscordInternal/logging.go#L61>)

```go
func LogLog(args ...any)
```



<a name="LogTrace"></a>
## func [LogTrace](<https://github.com/DoctorWhoFR/dapi/blob/master/DiscordInternal/logging.go#L55>)

```go
func LogTrace(args ...any)
```

LogTrace DONE big refactoring needed for logging need more readable logging, and less \`\`\` test \`\`\` \- $\{line\} \- $\{fullPath\} \- \[x\] test \- \[ \] test \- \[ \] test \- \[x\] test \- \[x\] test \<\!\-\- epic:"debugging" \#tag @1.0.1 order:0 completed:2023\-06\-22T18:46:48.383Z \-\-\>

<a name="MediumValueAPI"></a>
## func [MediumValueAPI](<https://github.com/DoctorWhoFR/dapi/blob/master/DiscordInternal/mediumCalculation.go#L9>)

```go
func MediumValueAPI() int64
```

MediumValueAPI return medium values of all api call

<a name="MediumValueBOT"></a>
## func [MediumValueBOT](<https://github.com/DoctorWhoFR/dapi/blob/master/DiscordInternal/mediumCalculation.go#L19>)

```go
func MediumValueBOT() int64
```

MediumValueBOT return medium values of all bot process time

<a name="MediumValueWithoutAPI"></a>
## func [MediumValueWithoutAPI](<https://github.com/DoctorWhoFR/dapi/blob/master/DiscordInternal/mediumCalculation.go#L30>)

```go
func MediumValueWithoutAPI() int64
```

MediumValueWithoutAPI comparison from discord call time and real bot process time

<a name="SimpleSleep"></a>
## func [SimpleSleep](<https://github.com/DoctorWhoFR/dapi/blob/master/DiscordInternal/sleeper.go#L7>)

```go
func SimpleSleep()
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
