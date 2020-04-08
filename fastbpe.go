package subwordnmt

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

const kEndWord = "</w>"
const kEndWordLength = 4
const kTokenDelim = "@@"
const kTokenDelimLength = 2

type fastBPE struct {
	codes map[pair]int
	vocab map[string]int
	reversedCodes map[string]pair
	codesPath string
	vocabPath string
}

type pair struct {
	x string
	y string
}

func FastBPE(codesPath string, vocabPath string) *fastBPE {
	bpe := &fastBPE{
		codes: map[pair]int{},
		vocab: map[string]int{},
		reversedCodes: map[string]pair{},
		codesPath: codesPath,
		vocabPath: vocabPath,
	}
	bpe.readVocab()
	bpe.readCodes()
	return bpe
}

func (bpe *fastBPE) Apply(sentences [][]string) [][]string {
	var results [][]string
	for _, tokens := range sentences {
		results = append(results, bpe.apply(tokens))
	}
	return results
}

func (bpe *fastBPE) ApplyString(sentences []string) [][]string {
	var results [][]string
	for _, sentence := range sentences {
		results = append(results, bpe.apply(strings.Fields(sentence)))
	}
	return results
}

func (bpe *fastBPE) apply(tokens []string) []string {
	var result []string
	for _, token := range tokens {
		var wordBPEs []string
		realLength := 0
		lastStart := 0
		for pos, char := range token {
			newChar := (char & 0xc0) != 0x80  // not a continuation byte
			if newChar {
				realLength++
				if pos > 0 {
					wordBPEs = append(wordBPEs, token[lastStart:pos])
					lastStart = pos
				}
			}
		}
		wordBPEs = append(wordBPEs, token[lastStart:] + kEndWord)
		result = append(result, bpe.process(wordBPEs))
	}
	return result
}

func (bpe *fastBPE) process(subwords []string) string {
	for len(subwords) > 1 {
		var bestPair *pair
		for i := 0; i < len(subwords) - 1; i++ {
			pair := pair{subwords[i], subwords[i+1]}
			pairRank, found := bpe.codes[pair]
			if found && (bestPair == nil || bpe.codes[*bestPair] > pairRank) {
				bestPair = &pair
			}
		}
		if bestPair == nil {
			break
		}
		justMerged := false
		var tmp []string
		for i := 0; i < len(subwords); i++ {
			if (i + 1 < len(subwords)) && !justMerged && subwords[i] == bestPair.x && subwords[i+1] == bestPair.y {
				tmp = append(tmp, subwords[i] + subwords[i+1])
				justMerged = true
			} else {
				if !justMerged {
					tmp = append(tmp, subwords[i])
				}
				justMerged = false
			}
		}
		subwords = tmp
	}
	subwords = bpe.limitVocab(subwords)
	var result string;
	for _, subword := range subwords {
		result += subword + kTokenDelim + " "
	}
	return result[:len(result)-kEndWordLength-kTokenDelimLength-1]
}

func (bpe *fastBPE) limitVocab(subwords []string) []string {
	var results []string
	for i, subword := range subwords {
		var query string
		isFinal := i == len(subwords) - 1
		if isFinal {
			query = subword[:len(subword) - kEndWordLength]
		} else {
			query = subword + kTokenDelim
		}
		if _, found := bpe.vocab[query]; found {
			results = append(results, subword)
		} else {
			bpe.decompose(subword, &results, isFinal)
		}
	}
	return results
}

func (bpe *fastBPE) decompose(s string, results *[]string, isFinal bool) {
	val, found := bpe.reversedCodes[s]
	if !found {
		*results = append(*results, s)
		return
	}
	if _, found = bpe.vocab[val.x + kTokenDelim]; found {
		*results = append(*results, val.x)
	} else {
		bpe.decompose(val.x, results, false)
	}
	var query string
	if isFinal {
		query = val.y[:len(val.y)-kEndWordLength]
	} else {
		query = val.y + kTokenDelim
	}
	if _, found = bpe.vocab[query]; found {
		*results = append(*results, val.y)
	} else {
		bpe.decompose(val.y, results, isFinal)
	}
}

func (bpe *fastBPE) readVocab() {
	f, err := os.OpenFile(bpe.vocabPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Cannot open vocabulary file %s", err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		splits := strings.Split(sc.Text(), " ")
		bpe.vocab[splits[0]], _ = strconv.Atoi(splits[1])
	}
}

func (bpe *fastBPE) readCodes() {
	f, err := os.OpenFile(bpe.codesPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Cannot open codes file %s", err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		splits := strings.Split(sc.Text(), " ")
		pair := pair{splits[0], splits[1]}
		bpe.codes[pair] = len(bpe.codes)
		bpe.reversedCodes[splits[0] + splits[1]] = pair
	}
}