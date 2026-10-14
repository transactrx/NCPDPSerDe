package responsesegment

import (
	"strings"
	"testing"

	reflectionutils "github.com/transactrx/NCPDPSerDe/pkg/reflectionUtils"
)

func messageSnapshot(t *testing.T, segment *Status) []string {
	t.Helper()

	got := make([]string, 0, len(segment.AdditionalMessages))
	for _, msg := range segment.AdditionalMessages {
		qfr, text, cont := "", "", ""
		if msg.Qualifier != nil {
			qfr = *msg.Qualifier
		}
		if msg.Message != nil {
			text = *msg.Message
		}
		if msg.Continuation != nil {
			cont = *msg.Continuation
		}
		got = append(got, qfr+"|"+text+"|"+cont)
	}

	return got
}

func assertMessages(t *testing.T, segment *Status, want []string) {
	t.Helper()

	got := messageSnapshot(t, segment)
	if len(got) != len(want) {
		t.Fatalf("message count mismatch. Wanted: %v   Got: %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("message[%d] mismatch. Wanted: %q   Got: %q", i, want[i], got[i])
		}
	}
	if segment.AdditionalMessageCount == nil || *segment.AdditionalMessageCount != len(want) {
		t.Errorf("AdditionalMessageCount mismatch. Wanted: %d   Got: %v", len(want), segment.AdditionalMessageCount)
	}
}

func TestAppendMessageToEmptySegment(t *testing.T) {
	segment := &Status{}

	err := segment.AppendMessage("short message")
	if err != nil {
		t.Fatalf("AppendMessage returned error: %v", err)
	}

	assertMessages(t, segment, []string{"01|SHORT MESSAGE|"})
}

func TestAppendMessageChunksWithContinuation(t *testing.T) {
	segment := &Status{}

	err := segment.AppendMessage(strings.Repeat("A", MaxMessageLength) + "TAIL")
	if err != nil {
		t.Fatalf("AppendMessage returned error: %v", err)
	}

	assertMessages(t, segment, []string{
		"01|" + strings.Repeat("A", MaxMessageLength) + "|+",
		"02|TAIL|",
	})
}

// Coded qualifier values (e.g. "12", "38" = Max Fair Price) must be retained
// untouched; free-text chunks get the next unused "01"-"09" values and are
// inserted ahead of the higher coded values to keep the transaction ordered.
func TestAppendMessagePreservesCodedQualifiers(t *testing.T) {
	segment := &Status{
		AdditionalMessages: []AdditionalMessage{
			{Qualifier: reflectionutils.ToPointer("12"), Message: reflectionutils.ToPointer("MAXIMUM DAILY DOSE OF 1.000000")},
			{Qualifier: reflectionutils.ToPointer("38"), Message: reflectionutils.ToPointer("405.66")},
		},
		AdditionalMessageCount: reflectionutils.ToPointer(2),
	}

	err := segment.AppendMessage("new gateway message")
	if err != nil {
		t.Fatalf("AppendMessage returned error: %v", err)
	}

	assertMessages(t, segment, []string{
		"01|NEW GATEWAY MESSAGE|",
		"12|MAXIMUM DAILY DOSE OF 1.000000|",
		"38|405.66|",
	})
}

func TestAppendMessageContinuesAfterExistingFreeText(t *testing.T) {
	segment := &Status{
		AdditionalMessages: []AdditionalMessage{
			{Qualifier: reflectionutils.ToPointer("01"), Message: reflectionutils.ToPointer("FIRST")},
			{Qualifier: reflectionutils.ToPointer("02"), Message: reflectionutils.ToPointer("SECOND")},
			{Qualifier: reflectionutils.ToPointer("38"), Message: reflectionutils.ToPointer("405.66")},
		},
		AdditionalMessageCount: reflectionutils.ToPointer(3),
	}

	err := segment.AppendMessage("third")
	if err != nil {
		t.Fatalf("AppendMessage returned error: %v", err)
	}

	assertMessages(t, segment, []string{
		"01|FIRST|",
		"02|SECOND|",
		"03|THIRD|",
		"38|405.66|",
	})
}

// Future ECL alpha values ("A0"-"ZZ") sort after all numerics; free text must
// be inserted ahead of them and alpha values must not disturb the numbering.
func TestAppendMessageInsertsBeforeAlphaQualifiers(t *testing.T) {
	segment := &Status{
		AdditionalMessages: []AdditionalMessage{
			{Qualifier: reflectionutils.ToPointer("AB"), Message: reflectionutils.ToPointer("FUTURE CODED VALUE")},
		},
		AdditionalMessageCount: reflectionutils.ToPointer(1),
	}

	err := segment.AppendMessage("free text")
	if err != nil {
		t.Fatalf("AppendMessage returned error: %v", err)
	}

	assertMessages(t, segment, []string{
		"01|FREE TEXT|",
		"AB|FUTURE CODED VALUE|",
	})
}

// Some payers send unpadded numeric qualifiers; "1" must count as free-text 01
// for both numbering and ordering.
func TestAppendMessageHandlesUnpaddedQualifiers(t *testing.T) {
	segment := &Status{
		AdditionalMessages: []AdditionalMessage{
			{Qualifier: reflectionutils.ToPointer("1"), Message: reflectionutils.ToPointer("UNPADDED FIRST")},
			{Qualifier: reflectionutils.ToPointer("38"), Message: reflectionutils.ToPointer("405.66")},
		},
		AdditionalMessageCount: reflectionutils.ToPointer(2),
	}

	err := segment.AppendMessage("second")
	if err != nil {
		t.Fatalf("AppendMessage returned error: %v", err)
	}

	assertMessages(t, segment, []string{
		"1|UNPADDED FIRST|",
		"02|SECOND|",
		"38|405.66|",
	})
}

func TestAppendMessageSuccessiveAppends(t *testing.T) {
	segment := &Status{
		AdditionalMessages: []AdditionalMessage{
			{Qualifier: reflectionutils.ToPointer("12"), Message: reflectionutils.ToPointer("MAXIMUM DAILY DOSE OF 1.000000")},
		},
		AdditionalMessageCount: reflectionutils.ToPointer(1),
	}

	if err := segment.AppendMessage("first note"); err != nil {
		t.Fatalf("AppendMessage returned error: %v", err)
	}
	if err := segment.AppendMessage("second note"); err != nil {
		t.Fatalf("AppendMessage returned error: %v", err)
	}

	assertMessages(t, segment, []string{
		"01|FIRST NOTE|",
		"02|SECOND NOTE|",
		"12|MAXIMUM DAILY DOSE OF 1.000000|",
	})
}

func TestAppendMessageFreeTextQualifiersExhausted(t *testing.T) {
	segment := &Status{
		AdditionalMessages: []AdditionalMessage{
			{Qualifier: reflectionutils.ToPointer("09"), Message: reflectionutils.ToPointer("LAST FREE TEXT SLOT")},
		},
		AdditionalMessageCount: reflectionutils.ToPointer(1),
	}

	err := segment.AppendMessage("no room left")
	if err == nil {
		t.Fatal("AppendMessage should have returned an error when free-text qualifiers are exhausted")
	}

	assertMessages(t, segment, []string{"09|LAST FREE TEXT SLOT|"})
}

func TestAppendMessageTooLongForRemainingQualifiers(t *testing.T) {
	segment := &Status{}

	// 10 chunks needed, only 9 free-text qualifiers exist.
	err := segment.AppendMessage(strings.Repeat("A", MaxMessageLength*9+1))
	if err == nil {
		t.Fatal("AppendMessage should have returned an error when the message needs more than 9 chunks")
	}

	if len(segment.AdditionalMessages) != 0 {
		t.Errorf("segment should be unchanged on error. Got: %v", messageSnapshot(t, segment))
	}
}
