package models

import (
	"errors"
	"fmt"
)

// SQLInsertStatement
// Struct to hold the insert statement string and insert values for the insertion
// of a Go struct into the corresponding SQL schema
type SQLInsertStatement struct {
	Statement string
	Values    []interface{}
}

// ///////////////////////////////////////////////////////////////

type CodeSource int

const (
	CodeSourcePost CodeSource = iota
	CodeSourceAttempt
	CodeSourceByte
)

func (s CodeSource) String() string {
	switch s {
	case CodeSourcePost:
		return "Challenge"
	case CodeSourceAttempt:
		return "Attempt"
	case CodeSourceByte:
		return "Byte"
	}
	return "Unknown"
}

// ///////////////////////////////////////////////////////////////

type ChallengeType int

const (
	InteractiveChallenge ChallengeType = iota
	PlaygroundChallenge
	CasualChallenge
	CompetitiveChallenge
	DebugChallenge
	BytesChallenge
)

func (c ChallengeType) String() string {
	switch c {
	case InteractiveChallenge:
		return "Interactive"
	case PlaygroundChallenge:
		return "Playground"
	case CasualChallenge:
		return "Casual"
	case CompetitiveChallenge:
		return "Competitive"
	case DebugChallenge:
		return "Debug"
	case BytesChallenge:
		return "Bytes"
	}
	return ""
}

// ///////////////////////////////////////////////////////////////

type AuthenticationRole int64

const (
	BaseUser AuthenticationRole = iota
	Admin
)

func (r AuthenticationRole) String() string {
	switch r {
	case BaseUser:
		return "BaseUser"
	case Admin:
		return "Admin"
	}
	return ""
}

// ///////////////////////////////////////////////////////////////

type TierType int

const (
	Tier1 TierType = iota
	Tier2
	Tier3
	Tier4
	Tier5
	Tier6
	Tier7
	Tier8
	Tier9
	Tier10
)

func (t TierType) String() string {
	switch t {
	case Tier1:
		return "Tier1"
	case Tier2:
		return "Tier2"
	case Tier3:
		return "Tier3"
	case Tier4:
		return "Tier4"
	case Tier5:
		return "Tier5"
	case Tier6:
		return "Tier6"
	case Tier7:
		return "Tier7"
	case Tier8:
		return "Tier8"
	case Tier9:
		return "Tier9"
	case Tier10:
		return "Tier10"
	default:
		return ""
	}
}

func TierTypeFromString(tierType string) (TierType, error) {
	switch tierType {
	case "Tier1":
		return Tier1, nil
	case "Tier2":
		return Tier2, nil
	case "Tier3":
		return Tier3, nil
	case "Tier4":
		return Tier4, nil
	case "Tier5":
		return Tier5, nil
	case "Tier6":
		return Tier6, nil
	case "Tier7":
		return Tier7, nil
	case "Tier8":
		return Tier8, nil
	case "Tier9":
		return Tier9, nil
	case "Tier10":
		return Tier10, nil
	default:
		return Tier1, errors.New(fmt.Sprintf("Invalid TierType input: %v", tierType))
	}
}

// ///////////////////////////////////////////////////////////////

type LevelType int

const (
	Level1 LevelType = iota
	Level2
	Level3
	Level4
	Level5
	Level6
	Level7
	Level8
	Level9
	Level10
)

// ///////////////////////////////////////////////////////////////

type RankType int

const (
	NoobRank RankType = iota
	BronzeRank
	SilverRank
	GoldRank
	PlatinumRank
	NeckBeardRank
	DeveloperRank
	FounderRank
)

func (r RankType) String() string {
	switch r {
	case NoobRank:
		return "Noob"
	case BronzeRank:
		return "Bronze"
	case SilverRank:
		return "Silver"
	case GoldRank:
		return "Gold"
	case PlatinumRank:
		return "Platinum"
	case NeckBeardRank:
		return "NeckBeard"
	}
	return ""
}

// ///////////////////////////////////////////////////////////////

type PostVisibility int

const (
	PublicVisibility PostVisibility = iota
	PrivateVisibility
	FriendsVisibility
	FollowerVisibility
	PremiumVisibility
	ExclusiveVisibility
)

func (w PostVisibility) String() string {
	switch w {
	case PublicVisibility:
		return "Public"
	case PrivateVisibility:
		return "Private"
	case FriendsVisibility:
		return "Friends"
	case FollowerVisibility:
		return "Follower"
	case PremiumVisibility:
		return "Premium"
	case ExclusiveVisibility:
		return "Exclusive"
	}
	return ""
}

// ///////////////////////////////////////////////////////////////

type ProgrammingLanguage int

const (
	AnyProgrammingLanguage ProgrammingLanguage = iota
	CustomProgrammingLanguage
	Java
	JavaScript
	TypeScript
	Python
	Go
	Ruby
	Cpp
	C
	Csharp
	ObjectiveC
	Swift
	PHP
	Rust
	Kotlin
	Dart
	Scala
	CoffeeScript
	Haskell
	Lua
	Clojure
	Perl
	Shell
	Elixir
	Assembly
	Groovy
	Html
	Julia
	OCaml
	R
	Ada
	Erlang
	Matlab
	SQL
	Cobol
	Lisp
	HCL
)

func (l ProgrammingLanguage) String() string {
	switch l {
	case AnyProgrammingLanguage:
		return "Any"
	case CustomProgrammingLanguage:
		return "Java"
	case Java:
		return "Java"
	case JavaScript:
		return "JavaScript"
	case TypeScript:
		return "TypeScript"
	case Python:
		return "Python"
	case Go:
		return "Go"
	case Ruby:
		return "Ruby"
	case Cpp:
		return "Cpp"
	case C:
		return "C"
	case Csharp:
		return "Csharp"
	case ObjectiveC:
		return "ObjectiveC"
	case Swift:
		return "Swift"
	case PHP:
		return "PHP"
	case Rust:
		return "Rust"
	case Kotlin:
		return "Kotlin"
	case Dart:
		return "Dart"
	case Scala:
		return "Scala"
	case CoffeeScript:
		return "CoffeeScript"
	case Haskell:
		return "Haskell"
	case Lua:
		return "Lua"
	case Clojure:
		return "Clojure"
	case Perl:
		return "Perl"
	case Shell:
		return "Shell"
	case Elixir:
		return "Elixir"
	case Assembly:
		return "Assembly"
	case Groovy:
		return "Groovy"
	case Html:
		return "Html"
	case Julia:
		return "Julia"
	case OCaml:
		return "OCaml"
	case R:
		return "R"
	case Ada:
		return "Ada"
	case Erlang:
		return "Erlang"
	case Matlab:
		return "Matlab"
	case SQL:
		return "SQL"
	case Cobol:
		return "Cobol"
	case Lisp:
		return "Lisp"
	case HCL:
		return "HCL"
	default:
		return "Unknown"
	}
}

//////////////////////////////////////////////

type CommunicationType int

const (
	DiscussionLevel CommunicationType = iota
	CommentLevel
	ThreadLevel
	ThreadReplyLevel
)

func (c CommunicationType) String() string {
	switch c {
	case DiscussionLevel:
		return "Discussion"
	case CommentLevel:
		return "Comment"
	case ThreadLevel:
		return "Thread"
	case ThreadReplyLevel:
		return "ThreadReply"
	}
	return ""
}
