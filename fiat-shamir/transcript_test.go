package fiatshamir

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func TestConstants(t *testing.T) {
	if min253.BitLen() != bitLimit {
		t.Errorf("Min253.BitLen != bitLimit")
	}
}

func TestTranscript(t *testing.T) {
	testStrings := []string{"111", "aaa", "333"}
	trans1 := InitTranscript(testStrings, Default)

	var trans2 Transcript
	trans2.AppendSlice(testStrings)

	var trans3 Transcript
	trans3.AppendSlice(testStrings[:2])
	trans3.Append(testStrings[2])

	challenge1 := trans1.GetPrimeChallengeUsingTranscript()
	challenge2 := trans2.GetPrimeChallengeUsingTranscript()
	challenge3 := trans3.GetPrimeChallengeUsingTranscript()

	if challenge1.Cmp(challenge2) != 0 {
		t.Errorf("Different ways to init transcript leads to different results")
		trans1.Print()
		trans2.Print()
	}

	if challenge1.Cmp(challenge3) != 0 {
		t.Errorf("Different ways to init transcript leads to different results")
		trans1.Print()
		trans3.Print()
	}
}

func TestChallenge(t *testing.T) {
	testStrings := []string{"111", "aaa", "333"}
	trans1 := InitTranscript(testStrings, Default)

	challenge1 := trans1.GetPrimeChallengeUsingTranscript()
	challenge2 := trans1.GetPrimeChallengeUsingTranscript()
	if challenge1.Cmp(challenge2) == 0 {
		t.Errorf("Updated transcript has old results")
		trans1.Print()
	}

	trans2 := InitTranscript(testStrings, Default)
	trans2.Append(challenge1.String())
	challenge3 := trans2.GetPrimeChallengeUsingTranscript()
	if challenge3.Cmp(challenge2) != 0 {
		t.Errorf("Different ways to update transcript leads to different results")
		trans1.Print()
		trans2.Print()
	}
}

func TestLargeChallenge(t *testing.T) {
	testStrings := []string{"111", "aaa", "333"}
	trans1 := InitTranscript(testStrings, Default)

	lenght := 2048
	challenge1 := trans1.GetLargeChallengeUsingTranscript(lenght)
	if challenge1.BitLen() != lenght && challenge1.BitLen() != lenght-1 {
		t.Error("Wrong bit length for large challenge, length = ", challenge1.BitLen())
	}
	lenght = 2048 + 233
	challenge1 = trans1.GetLargeChallengeUsingTranscript(lenght)
	if challenge1.BitLen() != lenght && challenge1.BitLen() != lenght-1 {
		t.Error("Wrong bit length for large challenge, length = ", challenge1.BitLen())
	}
	lenght = 2048 * 256
	challenge1 = trans1.GetLargeChallengeUsingTranscript(lenght)
	// we may loose several bits because the leading bits are 0
	if challenge1.BitLen() != lenght && challenge1.BitLen() != lenght-1 && challenge1.BitLen() != lenght-2 {
		//fmt.Println("Challenge = ", challenge1.String())
		t.Error("Wrong bit length for large challenge, length = ", challenge1.BitLen())
	}
}

func TestPrimeChallengeLength(t *testing.T) {
	testStrings := []string{"111", "aaa", "333"}
	trans1 := InitTranscript(testStrings, Max252)

	challenge1 := trans1.GetPrimeChallengeUsingTranscript()
	challenge2 := trans1.GetPrimeChallengeUsingTranscript()
	if !challenge1.ProbablyPrime(securityParameter) {
		t.Errorf("Challenge not prime")
	}
	if !challenge2.ProbablyPrime(securityParameter) {
		t.Errorf("Challenge not prime")
	}

	if challenge1.Cmp(fr.Modulus()) != -1 {
		t.Errorf("Challenge larger than fr.Modulus()")
	}
	if challenge2.Cmp(fr.Modulus()) != -1 {
		t.Errorf("Challenge larger than fr.Modulus()")
	}
	if challenge1.Cmp(&min253) != -1 {
		t.Errorf("Challenge larger than min253")
	}
	if challenge2.Cmp(&min253) != -1 {
		t.Errorf("Challenge larger than min253")
	}
}

func FuzzChallenge(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345", "123123123", "0.7"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, testInput string) {
		trans1 := InitTranscript([]string{testInput}, Default)
		challenge1 := trans1.GetPrimeChallengeUsingTranscript()
		if !challenge1.ProbablyPrime(securityParameter) {
			t.Errorf("Challenge not prime")
		}
		trans1 = InitTranscript([]string{testInput}, Max252)
		challenge1 = trans1.GetPrimeChallengeUsingTranscript()
		if !challenge1.ProbablyPrime(securityParameter) {
			t.Errorf("Challenge not prime")
		}
		if challenge1.Cmp(fr.Modulus()) != -1 {
			t.Errorf("Challenge larger than fr.Modulus()")
		}
		if challenge1.Cmp(&min253) != -1 {
			t.Errorf("Challenge larger than min253")
		}
	})
}
