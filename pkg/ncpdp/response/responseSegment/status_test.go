package responsesegment

import (
	"fmt"
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

// NCPDP allows gaps in the qualifier sequence. The UF count must always equal
// the number of UH/FQ occurrences present, never the highest qualifier value,
// and appending continues after the highest free-text value without filling
// gaps (which could collide with a receiver's interpretation of order).
func TestAppendMessageCountWithGappedQualifiers(t *testing.T) {
	segment := &Status{
		AdditionalMessages: []AdditionalMessage{
			{Qualifier: reflectionutils.ToPointer("01"), Message: reflectionutils.ToPointer("FIRST")},
			{Qualifier: reflectionutils.ToPointer("03"), Message: reflectionutils.ToPointer("THIRD, GAP AT 02")},
			{Qualifier: reflectionutils.ToPointer("38"), Message: reflectionutils.ToPointer("405.66")},
		},
		AdditionalMessageCount: reflectionutils.ToPointer(3),
	}

	err := segment.AppendMessage("new note")
	if err != nil {
		t.Fatalf("AppendMessage returned error: %v", err)
	}

	// Next free text is 04 (after highest free-text 03), UF is 4 (entry count).
	assertMessages(t, segment, []string{
		"01|FIRST|",
		"03|THIRD, GAP AT 02|",
		"04|NEW NOTE|",
		"38|405.66|",
	})
}

func TestAppendMessageCountWithGappedQualifiersMultiChunk(t *testing.T) {
	segment := &Status{
		AdditionalMessages: []AdditionalMessage{
			{Qualifier: reflectionutils.ToPointer("02"), Message: reflectionutils.ToPointer("ONLY MESSAGE, GAP AT 01")},
		},
		AdditionalMessageCount: reflectionutils.ToPointer(1),
	}

	err := segment.AppendMessage(strings.Repeat("A", MaxMessageLength) + "TAIL")
	if err != nil {
		t.Fatalf("AppendMessage returned error: %v", err)
	}

	assertMessages(t, segment, []string{
		"02|ONLY MESSAGE, GAP AT 01|",
		"03|" + strings.Repeat("A", MaxMessageLength) + "|+",
		"04|TAIL|",
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

func segmentWithStructuredMessages(count int) (*Status, []string) {
	segment := &Status{}
	snapshot := make([]string, 0, count)
	for i := 0; i < count; i++ {
		qualifier := fmt.Sprintf("%02d", maxFreeTextQualifier+1+i)
		segment.AdditionalMessages = append(segment.AdditionalMessages, AdditionalMessage{
			Qualifier: reflectionutils.ToPointer(qualifier),
			Message:   reflectionutils.ToPointer("VALUE"),
		})
		snapshot = append(snapshot, qualifier+"|VALUE|")
	}
	segment.AdditionalMessageCount = reflectionutils.ToPointer(count)

	return segment, snapshot
}

func TestAppendMessageAllowsExactlyMaxMessageCount(t *testing.T) {
	segment, existing := segmentWithStructuredMessages(MaxMessageCount - 1)

	if err := segment.AppendMessage("last"); err != nil {
		t.Fatalf("AppendMessage returned error: %v", err)
	}

	assertMessages(t, segment, append([]string{"01|LAST|"}, existing...))
}

func TestAppendMessageExceedingMaxMessageCountLeavesSegmentUnchanged(t *testing.T) {
	segment, existing := segmentWithStructuredMessages(MaxMessageCount - 1)

	err := segment.AppendMessage(strings.Repeat("A", MaxMessageLength) + "TAIL")
	if err == nil {
		t.Fatal("AppendMessage should have returned an error when the message count is exceeded")
	}

	assertMessages(t, segment, existing)
}

func TestSetMessageReplacesFreeTextAndPreservesStructuredQualifiers(t *testing.T) {
	segment := &Status{
		AdditionalMessages: []AdditionalMessage{
			{Qualifier: reflectionutils.ToPointer("01"), Message: reflectionutils.ToPointer("FIRST")},
			{Qualifier: reflectionutils.ToPointer("2"), Message: reflectionutils.ToPointer("UNPADDED SECOND")},
			{Qualifier: reflectionutils.ToPointer("12"), Message: reflectionutils.ToPointer("MAXIMUM DAILY DOSE OF 1.000000")},
			{Qualifier: reflectionutils.ToPointer("38"), Message: reflectionutils.ToPointer("405.66")},
		},
		AdditionalMessageCount: reflectionutils.ToPointer(4),
	}

	err := segment.SetMessage(strings.Repeat("A", MaxMessageLength) + "TAIL")
	if err != nil {
		t.Fatalf("SetMessage returned error: %v", err)
	}

	assertMessages(t, segment, []string{
		"01|" + strings.Repeat("A", MaxMessageLength) + "|+",
		"02|TAIL|",
		"12|MAXIMUM DAILY DOSE OF 1.000000|",
		"38|405.66|",
	})
}

func TestSetMessageToEmptySegment(t *testing.T) {
	segment := &Status{}

	err := segment.SetMessage("short message")
	if err != nil {
		t.Fatalf("SetMessage returned error: %v", err)
	}

	assertMessages(t, segment, []string{"01|SHORT MESSAGE|"})
}

func TestSetMessageEmptyRemovesFreeText(t *testing.T) {
	segment := &Status{
		AdditionalMessages: []AdditionalMessage{
			{Qualifier: reflectionutils.ToPointer("01"), Message: reflectionutils.ToPointer("FIRST")},
			{Qualifier: reflectionutils.ToPointer("38"), Message: reflectionutils.ToPointer("405.66")},
		},
		AdditionalMessageCount: reflectionutils.ToPointer(2),
	}

	if err := segment.SetMessage(""); err != nil {
		t.Fatalf("SetMessage returned error: %v", err)
	}

	assertMessages(t, segment, []string{"38|405.66|"})
}

func TestSetMessageEmptyClearsCount(t *testing.T) {
	segment := &Status{
		AdditionalMessages: []AdditionalMessage{
			{Qualifier: reflectionutils.ToPointer("01"), Message: reflectionutils.ToPointer("FIRST")},
		},
		AdditionalMessageCount: reflectionutils.ToPointer(1),
	}

	if err := segment.SetMessage(""); err != nil {
		t.Fatalf("SetMessage returned error: %v", err)
	}

	if len(segment.AdditionalMessages) != 0 || segment.AdditionalMessageCount != nil {
		t.Errorf("expected no messages and nil count. Got: %v count %v", messageSnapshot(t, segment), segment.AdditionalMessageCount)
	}
}

func TestSetMessageNeedingMoreThanNineChunksLeavesSegmentUnchanged(t *testing.T) {
	segment := &Status{
		AdditionalMessages: []AdditionalMessage{
			{Qualifier: reflectionutils.ToPointer("01"), Message: reflectionutils.ToPointer("FIRST")},
			{Qualifier: reflectionutils.ToPointer("38"), Message: reflectionutils.ToPointer("405.66")},
		},
		AdditionalMessageCount: reflectionutils.ToPointer(2),
	}

	err := segment.SetMessage(strings.Repeat("A", MaxMessageLength*9+1))
	if err == nil {
		t.Fatal("SetMessage should have returned an error when the message needs more than 9 chunks")
	}

	assertMessages(t, segment, []string{"01|FIRST|", "38|405.66|"})
}

func TestSetMessageRestartsAtFirstQualifierWhenFreeTextExhausted(t *testing.T) {
	segment := &Status{}
	for i := 1; i <= maxFreeTextQualifier; i++ {
		segment.AdditionalMessages = append(segment.AdditionalMessages, AdditionalMessage{
			Qualifier: reflectionutils.ToPointer(fmt.Sprintf("%02d", i)),
			Message:   reflectionutils.ToPointer(fmt.Sprintf("LINE %d", i)),
		})
	}
	segment.AdditionalMessages = append(segment.AdditionalMessages, AdditionalMessage{
		Qualifier: reflectionutils.ToPointer("38"), Message: reflectionutils.ToPointer("405.66"),
	})
	segment.AdditionalMessageCount = reflectionutils.ToPointer(len(segment.AdditionalMessages))

	if err := segment.SetMessage("replacement"); err != nil {
		t.Fatalf("SetMessage returned error: %v", err)
	}

	assertMessages(t, segment, []string{"01|REPLACEMENT|", "38|405.66|"})
}

func TestSetMessageTwiceKeepsOnlyLast(t *testing.T) {
	segment := &Status{
		AdditionalMessages: []AdditionalMessage{
			{Qualifier: reflectionutils.ToPointer("12"), Message: reflectionutils.ToPointer("MAXIMUM DAILY DOSE OF 1.000000")},
		},
		AdditionalMessageCount: reflectionutils.ToPointer(1),
	}

	if err := segment.SetMessage("first"); err != nil {
		t.Fatalf("SetMessage returned error: %v", err)
	}
	if err := segment.SetMessage("second"); err != nil {
		t.Fatalf("SetMessage returned error: %v", err)
	}

	assertMessages(t, segment, []string{"01|SECOND|", "12|MAXIMUM DAILY DOSE OF 1.000000|"})
}

func TestSetMessageDoesNotModifyCallerHeldSlice(t *testing.T) {
	segment := &Status{
		AdditionalMessages: []AdditionalMessage{
			{Qualifier: reflectionutils.ToPointer("01"), Message: reflectionutils.ToPointer("FIRST")},
			{Qualifier: reflectionutils.ToPointer("38"), Message: reflectionutils.ToPointer("405.66")},
		},
		AdditionalMessageCount: reflectionutils.ToPointer(2),
	}
	held := segment.AdditionalMessages

	if err := segment.SetMessage("replacement"); err != nil {
		t.Fatalf("SetMessage returned error: %v", err)
	}

	got := messageSnapshot(t, &Status{AdditionalMessages: held})
	want := []string{"01|FIRST|", "38|405.66|"}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("held message[%d] modified. Wanted: %q   Got: %q", i, want[i], got[i])
		}
	}
}

func TestSetMessageRetainsMessageWithoutQualifier(t *testing.T) {
	segment := &Status{
		AdditionalMessages: []AdditionalMessage{
			{Message: reflectionutils.ToPointer("NO QUALIFIER")},
			{Qualifier: reflectionutils.ToPointer("01"), Message: reflectionutils.ToPointer("FIRST")},
		},
		AdditionalMessageCount: reflectionutils.ToPointer(2),
	}

	if err := segment.SetMessage("replacement"); err != nil {
		t.Fatalf("SetMessage returned error: %v", err)
	}

	assertMessages(t, segment, []string{"|NO QUALIFIER|", "01|REPLACEMENT|"})
}

func TestSetMessageInsertsFreeTextBeforeHigherQualifiersListedFirst(t *testing.T) {
	segment := &Status{
		AdditionalMessages: []AdditionalMessage{
			{Qualifier: reflectionutils.ToPointer("38"), Message: reflectionutils.ToPointer("405.66")},
			{Qualifier: reflectionutils.ToPointer("01"), Message: reflectionutils.ToPointer("FIRST")},
		},
		AdditionalMessageCount: reflectionutils.ToPointer(2),
	}

	if err := segment.SetMessage("replacement"); err != nil {
		t.Fatalf("SetMessage returned error: %v", err)
	}

	assertMessages(t, segment, []string{"01|REPLACEMENT|", "38|405.66|"})
}

func TestSetMessageNilSegment(t *testing.T) {
	var segment *Status

	if err := segment.SetMessage("message"); err != nil {
		t.Fatalf("SetMessage on nil segment returned error: %v", err)
	}
}

func TestSetMessageMaxMessageCountLeavesSegmentUnchanged(t *testing.T) {
	segment, existing := segmentWithStructuredMessages(MaxMessageCount)

	if err := segment.SetMessage("no room"); err == nil {
		t.Fatal("SetMessage should have returned an error when the message count is exceeded")
	}

	assertMessages(t, segment, existing)
}
