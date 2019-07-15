package main

import (
	"fmt"

	"github.com/smeshkov/lsh"
)

func main() {

	// 1. get shingles from given texts
	aShingles := lsh.Shingle(textA)
	bShingles := lsh.Shingle(textB)

	// 2. minhash given shingles to get signature matrix
	signatureMatrix := lsh.Minhash([][]string{aShingles, bShingles}, 20)

	// 3. perform LSH on the built signature matrix
	bandBuckets := lsh.LSH(signatureMatrix, 4)

	// 4. find candidates in the return buckets
	candidates := bandBuckets.FindCandidates()

	// 5. print results
	fmt.Printf("found %d candidate pair(s)\n", len(candidates.Index))
	if len(candidates.Index) > 0 {
		fmt.Printf("%v\n", candidates.Keys())
	}
}

var (
	// https://www.fbtb.net/video-games/nintendo-switch/2019/07/10/nintendo-switch-lite-announced/
	textA = []string{
		"Nintendo Switch Lite Announced",
		"This morning, Nintendo announced the next iteration of the Switch console: the Nintendo Switch Lite. A dedicated handheld console, the Lite won’t connect to a TV or do any of the Joy-Con switching of the original Switch. However, for someone who wants to use a Switch exclusively in handheld mode, the console is a steal.",
		"Three colors – turqouise, gray, and yellow. Not to mention a clean light gray Pokemon Sword and Shield variant.",
		"At a starting price of $199, $100 less than the original, the Lite still manages to incorporate a lot of the functionality. A 5.5-inch touch screen, 720p resolution, for example. But, there are also some trade offs. For one, motion control and IR sensors are gone, which could lead to some incompatibility in certain games.",
		"The Lite also boasts a “slightly” better battery life, a more power-efficient chip layout, and no additional batteries in the built-in controllers. As well as actually adding a D-pad, so it’s possible to use that for games.",
		"All in all, the Lite looks like a solid addition, though despite what the trailer shows, I can’t help but feel like it does away with the social aspect of the Switch. Oh well, it’s not like I was playing it with anyone to begin with.",
		"The Lite launches on September 20th. Check out more detail on Nintendo’s website.",
	}

	// source: https://gematsu.com/2019/07/nintendo-switch-lite-announced-launches-september-20-for-199
	textB = []string{
		"Nintendo Switch Lite announced, launches September 20 for $199",
		"Nintendo has announced Nintendo Switch Lite, a new version of Switch designed to play Switch games in handheld mode. It will launch worldwide on September 20 for $199.99 in yellow, gray, and turquoise color options. A carrying case and screen protector set will also be available.",
		"“Adding Nintendo Switch Lite to the lineup gives gamers more color and price point options,” said Nintendo of America president Doug Bowser in a press release. “Now consumers can choose the system that best suits how they like to play their favorite Nintendo Switch games.”",
		"Compared to the current Switch model, which includes detachable Joy-Cons, the Switch Lite has integrated controls and a smaller size. It has no kickstand and does not support video output to a TV. It can play all games in the Switch library that support handheld mode, though some games will have restrictions.",
		"Nintendo will also release a Switch Lite Zacian and Zamazenta Edition alongside the release of Pokemon Sword and Pokemon Shield on November 8 for $199.99. The gray model features cyan and magenta buttons, and illustrations of the two legendary Pokemon from the game.",
		"An official Switch system comparison website is now live.",
		"Watch a trailer below.",
	}
)
