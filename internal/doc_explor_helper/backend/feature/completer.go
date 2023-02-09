package feature

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
)

type Completer[I result.Info] interface {
	Main[I]
	Complete(info *I, source string)
}
