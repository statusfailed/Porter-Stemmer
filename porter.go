package golearn

import(
  "strings"
)

// Stems a word using Martin Porter's word stemming algorithm described here:
// http://tartarus.org/~martin/PorterStemmer/
// Only takes words in lowercase.
func Porter(word string) string {
  word = strings.TrimSpace(word)
  word = strings.ToLower(word)

  // Web page says to ignore strings length lteq 2, or things like "as"
  // tend to end up like "a".
  if(len(word) <= 2) { return word; }

  //fmt.Printf("step 0[%d] %s\n", len(word), word)
  word = step1(word)
  //fmt.Printf("step 1 [%d] %s\n", len(word), word)
  word = step2(word)
  //fmt.Printf("step 2 [%d] %s\n", len(word), word)
  word = step3(word)
  //fmt.Printf("step 3 [%d] %s\n", len(word), word)
  word = step4(word)
  //fmt.Printf("step 4 [%d] %s\n", len(word), word)
  return step5(word)
  //fmt.Printf("step 5 [%d] %s\n", len(word), word)
  return word
}

// Utility functions

// True if letter 'i' is a vowel.
func cons(word string, i int) bool {
  switch(word[i]) {
  case 'a', 'e', 'i', 'o', 'u': return false
  case 'y':
    if (i <= 0) {
      return true
    } else {
      return !cons(word, i-1)
    }
  default: return true
  }

  return true
}

// Returns the 'measure' of the word by Porter's definition.
// That is, counts the number of vowel sequences, subtracting 1 if it ends with a vowel.
// It does not include the string "suffix" if present.
func m(word, suffix string) int32 {
  // Bug check! make sure length of stem (i.e. word - suffix) is not
  // less than zero
  if len(word) - len(suffix) <= 0 { return 0 }

  word = word[:len(word)-len(suffix)]

  // current state (true if last was vowel)
  var lastVowel bool

  // Number of Vowel lists
  var n int32

  n = 0

  if cons(word, 0) {
    lastVowel = false
  } else {
    lastVowel = true
  }

  // For each letter...
  for i := range(word) {
    //fmt.Printf("%d %c %t %d\n", i, word[i], lastVowel, n)
    if cons(word, i) && lastVowel {
      // If consonant following vowel
      lastVowel = false
      n++
    }

    if !cons(word, i) {
      lastVowel = true
    }

  }

  return n
}

// true if word has a vowel before suffix.
func containsvowel(word string, suffix string) bool {
  for i := 0; i < len(word)-len(suffix); i++ {
    if !cons(word, i) {
      return true
    }
  }
  return false
}

// True if word ends in a double consonant, e.g. TT
func doublec(word string) bool {
  i := len(word) - 1 // Index of last char
  if len(word) < 2 { return false }
  if word[i] != word[i-1] { return false }
  return cons(word, i)
}

func cvc(word string) bool {
  i := len(word)-1
  if(len(word) < 3) { return false }
  if cons(word, i) && !cons(word, i-1) && cons(word, i-2) {
    switch(word[i]) {
    case 'w', 'x', 'y': return false;
    default: return true
    }
  }
  return false
}

func step1(word string) string {
  word = step1a(word)
  word = step1b(word)
  return step1c(word)
}

func step1a(word string) string {
  if strings.HasSuffix(word, "sses") {
    return word[:len(word)-2]
  }

  if strings.HasSuffix(word, "ies") {
    return word[:len(word)-2]
  }

  if strings.HasSuffix(word, "ss") {
    return word
  }

  if strings.HasSuffix(word, "s") {
    return word[:len(word)-1]
  }

  // Otherwise return unchanged.
  return word
}

func step1b(word string) string {
  if strings.HasSuffix(word, "eed") {
    if m(word, "eed") > 0 {
      return word[:len(word)-1]
    } else {
      return word
    }
  }

  if strings.HasSuffix(word, "ed") {
    if containsvowel(word, "ed") {
      return step1b2(word[:len(word)-2])
    } else {
      return word
    }
  }

  if strings.HasSuffix(word, "ing") {
    if containsvowel(word, "ing") {
      return step1b2(word[:len(word)-3])
    } else {
      return word
    }
  }

  return word
}

func step1b2(word string) string {
  if strings.HasSuffix(word, "at") { return word+"e" }
  if strings.HasSuffix(word, "bl") { return word+"e" }
  if strings.HasSuffix(word, "iz") { return word+"e" }
  if doublec(word) &&
    !(strings.HasSuffix(word, "l") ||
      strings.HasSuffix(word, "s") ||
      strings.HasSuffix(word, "z")) {
    return word[:len(word)-1]
  }
  if (m(word, "") == 1) && cvc(word) {
    return word+"e"
  }

  return word
}

func step1c(word string) string {
  if containsvowel(word, "y") && strings.HasSuffix(word, "y") {
    return word[:len(word)-1]+"i"
  }

  return word
}

// If word ends in suffix and m(word) > gtm, replace suffix with end.
// Returns stemmed word and true if rule was applied.
// If word is unchanged, returns word and false.
func replace(word, suffix, end string, gtm int32) (string, bool) {
  //fmt.Printf("checking rule: %s %s %s %d\n", word, suffix, end, gtm)
  if(len(suffix) >= len(word)) { return word, false } // If suffix is too big.
  if (m(word, suffix) > gtm) && strings.HasSuffix(word, suffix) {
    return word[:len(word)-len(suffix)] + end, true
  }
  return word, false
}

type rule struct {
  suffix string // Suffix to search for
  end string // Value to replace with
  gtm int32 // Min value of m minus 1
}

// Check all rules
func replaceRules(word string, rules []rule) (string, bool) {
  applied := false
  for i := 0; i < len(rules); i++ {
    word, applied = replace(word, rules[i].suffix, rules[i].end, rules[i].gtm)
    if applied {
      break
    } else if strings.HasSuffix(word, rules[i].suffix) {
      // If matches but not applied, return word and false.
      return word, false
    }
  }

  return word, applied
}

func step2(word string) string {
  rules := []rule{
    {"ational", "ate", 0},
    {"tional", "tion", 0},
    {"enci", "ence", 0},
    {"anci", "ance", 0},
    {"izer", "ize", 0},
    // This is changed from the published algorithm
    {"bli", "ble", 0},
    {"alli", "al", 0},
    {"entli", "ent", 0},
    {"eli", "e", 0},
    {"ousli", "ous", 0},
    {"ization", "ize", 0},
    {"ation", "ate", 0},
    {"ator", "ate", 0},
    {"alism", "al", 0},
    {"iveness", "ive", 0},
    {"fulness", "ful", 0},
    {"ousness", "ous", 0},
    {"aliti", "al", 0},
    {"iviti", "ive", 0},
    {"biliti", "ble", 0},
    // This was added and is not in the paper. See Porter web page.
    {"logi", "log", 0},
  }
  ret, _ := replaceRules(word, rules)
  return ret
}

func step3(word string) string {
  rules := []rule{
    {"icate", "ic", 0},
    {"ative", "", 0},
    {"alize", "al", 0},
    {"iciti", "ic", 0},
    {"ical", "ic", 0},
    {"ful", "", 0},
    {"ness", "", 0},
  }
  ret, _ := replaceRules(word, rules)
  return ret
}

func step4(word string) string {
  applied := false

  rules1 := []rule{
    {"al", "", 1},
    {"ance", "", 1},
    {"ence", "", 1},
    {"er", "", 1},
    {"ic", "", 1},
    {"able", "", 1},
    {"ible", "", 1},
    {"ant", "", 1},
    {"ement", "", 1},
    {"ment", "", 1},
    {"ent", "", 1},
  }

  ret, applied := replaceRules(word, rules1)
  if applied { return ret }

  // Calling m(word, "sion") first will result in an error.
  if (strings.HasSuffix(word, "sion") || strings.HasSuffix(word, "tion")) &&
     m(word, "ion") > 1 {
    return word[:len(word)-3]
  }

  rules2 := []rule{
    {"ou", "", 1},
    {"ism", "", 1},
    {"ate", "", 1},
    {"iti", "", 1},
    {"ous", "", 1},
    {"ive", "", 1},
    {"ize", "", 1},
  }
  ret, applied = replaceRules(word, rules2)
  if applied { return ret }

  return word
}

func step5(word string) string {
  word = step5a(word)
  return step5b(word)
}

func step5a(word string) string {
  if strings.HasSuffix(word, "e") && (m(word, "e") > 1) {
    return word[:len(word)-1]
  }

  if strings.HasSuffix(word, "e") && (m(word, "e") == 1) && !cvc(word[:len(word)-1]) {
    return word[:len(word)-1]
  }

  return word
}

func step5b(word string) string {
  if m(word, "l") > 1 && doublec(word) && strings.HasSuffix(word, "l") {
    return word[:len(word)-1]
  }
  return word
}
