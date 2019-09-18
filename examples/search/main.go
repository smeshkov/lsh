package main

import (
	"fmt"

	"github.com/smeshkov/lsh"
)

func main() {

	// 0. rember the order of documents
	docsToIndex := map[string]int{
		"textA": 0,
		"textB": 1,
		"textC": 2,
		"textD": 3,
		"textE": 4,
		"textF": 5,
	}
	indexToDocs := []string{
		"textA",
		"textB",
		"textC",
		"textD",
		"textE",
		"textF",
	}

	// 1. get shingles from given texts
	a := lsh.Shingle(textA)
	b := lsh.Shingle(textB)
	c := lsh.Shingle(textC)
	d := lsh.Shingle(textD)
	e := lsh.Shingle(textE)

	// 2. build a sets matrix which will serve as an index
	setsMatrix := lsh.ToSetsMatrix([][]string{a, b, c, d, e})

	// 3. create an instance of search based on "setsMatrix" as an index
	search := lsh.NewSearch(
		lsh.Index(setsMatrix),
		lsh.HashersNum(100),
		lsh.BandsNum(20),
	)

	// 4. find all candidates
	allCandidates := search.Find(textF)

	// 5. get candidates for document "textF" sorted by elections.
	candidates := allCandidates.GetByKeySorted(docsToIndex["textF"])

	// 5. print results
	fmt.Printf("found %d candidate(s)\n", len(candidates))
	if len(candidates) > 0 {
		for k, v := range candidates {
			fmt.Printf("[%d] %v\n", k, indexToDocs[v.Index])
		}
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

	// https://en.wikipedia.org/wiki/Locality-sensitive_hashing
	textB = []string{
		"The task of finding nearest neighbours is very common. You can think of applications like finding duplicate or similar documents, audio/video search. Although using brute force to check for all possible combinations will give you the exact nearest neighbour but it’s not scalable at all. Approximate algorithms to accomplish this task has been an area of active research. Although these algorithms don’t guarantee to give you the exact answer, more often than not they’ll be provide a good approximation. These algorithms are faster and scalable.",
		"Locality sensitive hashing (LSH) is one such algorithm. LSH has many applications, including:",
		"Near-duplicate detection: LSH is commonly used to deduplicate large quantities of documents, webpages, and other files.",
		"Genome-wide association study: Biologists often use LSH to identify similar gene expressions in genome databases.",
		"Large-scale image search: Google used LSH along with PageRank to build their image search technology VisualRank.",
		"Audio/video fingerprinting: In multimedia technologies, LSH is widely used as a fingerprinting technique A/V data.",
		"In this blog, we’ll try to understand the workings of this algorithm.",
		"General Idea",
		"LSH refers to a family of functions (known as LSH families) to hash data points into buckets so that data points near each other are located in the same buckets with high probability, while data points far from each other are likely to be in different buckets. This makes it easier to identify observations with various degrees of similarity.",
		"Finding similar documents",
	}

	// https://en.wikipedia.org/wiki/Go_(programming_language)
	textC = []string{
		"Go, also known as Golang,[14] is a statically typed, compiled programming language designed at Google[15] by Robert Griesemer, Rob Pike, and Ken Thompson.[12] Go is syntactically similar to C, but with memory safety, garbage collection, structural typing,[6] and CSP-style concurrency.",
		"Go was designed at Google in 2007 to improve programming productivity in an era of multicore, networked machines and large codebases.[23] The designers wanted to address criticism of other languages in use at Google, but keep their useful characteristics:",
		"Static typing and run-time efficiency (like C++)",
		"Readability and usability (like Python or JavaScript)",
		"High-performance networking and multiprocessing",
		"The designers were primarily motivated by their shared dislike of C++.",
		"Go was publicly announced in November 2009,[29] and version 1.0 was released in March 2012.[30][31] Go is widely used in production at Google[32] and in many other organizations and open-source projects.",
		"In November 2016, the Go and Go Mono fonts which are sans-serif and monospaced respectively were released by type designers Charles Bigelow and Kris Holmes. Both were designed to be legible with a large x-height and distinct letterforms by conforming to the DIN 1450 standard.",
		"In April 2018, the original logo was replaced with a stylized GO slanting right with trailing streamlines. However, the Gopher mascot remained the same.",
		"In August 2018, the Go principal contributors published two ″draft designs″ for new language features, Generics and Error Handling, and asked Go users to submit feedback on them.[36][37] Lack of support for generic programming and the verbosity of Error Handling in Go 1.x had drawn considerable criticism.",
	}

	// source: https://gematsu.com/2019/07/nintendo-switch-lite-announced-launches-september-20-for-199
	textD = []string{
		"Nintendo Switch Lite announced, launches September 20 for $199",
		"Nintendo has announced Nintendo Switch Lite, a new version of Switch designed to play Switch games in handheld mode. It will launch worldwide on September 20 for $199.99 in yellow, gray, and turquoise color options. A carrying case and screen protector set will also be available.",
		"“Adding Nintendo Switch Lite to the lineup gives gamers more color and price point options,” said Nintendo of America president Doug Bowser in a press release. “Now consumers can choose the system that best suits how they like to play their favorite Nintendo Switch games.”",
		"Compared to the current Switch model, which includes detachable Joy-Cons, the Switch Lite has integrated controls and a smaller size. It has no kickstand and does not support video output to a TV. It can play all games in the Switch library that support handheld mode, though some games will have restrictions.",
		"Nintendo will also release a Switch Lite Zacian and Zamazenta Edition alongside the release of Pokemon Sword and Pokemon Shield on November 8 for $199.99. The gray model features cyan and magenta buttons, and illustrations of the two legendary Pokemon from the game.",
		"An official Switch system comparison website is now live.",
		"Watch a trailer below.",
	}

	// https://variety.com/2019/film/news/box-office-spider-man-far-from-home-wednesday-record-1203259644/
	textE = []string{
		"Sony’s “Spider-Man: Far From Home” is dominating North American moviegoing, soaring to $27 million on Wednesday — a record for a Marvel Cinematic Universe film.",
		"The 23rd Marvel movie, starring Tom Holland, topped the previous MCU Wednesday mark set two months ago by Marvel’s “Avengers: Endgame” with $25.3 million. “Spider-Man: Far From Home” has earned a dazzling $65.5 million in its first two days from 4,634 North American locations.",
		"Sony is forecasting that the superhero tale will take in $125 million during its first six days in theaters, though some estimates show that number could reach $150 million. 2017’s “Spider-Man: Homecoming,” debuted domestically with $117 million over the three-day frame on its way to a $334 million North American total and $880 million worldwide.",
		"“Spider-Man: Far From Home” picks up after the events of “Avengers: Endgame” and follows Peter Parker (Holland) being  recruited to save the world while on a class trip to Europe. Jake Gyllenhaal joins the cast as Mysterio, while Samuel L. Jackson, Zendaya, Cobie Smulders, Jon Favreau and Marisa Tomei reprise their roles for the sequel.",
		"Jon Watts directed “Spider-Man: Far From Home,” which cost $160 million to produce, from a script by Chris McKenna and Erik Sommers. Critics have been mostly impressed with the movie, which carries a 92% fresh rating on Rotten Tomatoes.",
		"Disney’s “Toy Story 4” should remain a potent attraction during the holiday weekend. The animated comedy sequel has taken in $256 million domestically in 12 days through Tuesday.",
		"Also opening this weekend is A24’s horror film “Midsommar,” which film collected $1.1 million from 1,951 screens on Tuesday. “Midsommar” is set to bring in around $8 million to $10 million over the weekend when it opens in 2,707 venues. Directed by Ari Aster, “Midsommar” stars Florence Pugh, Jack Reynor and William Jackson Harper in a nightmarish story of friends who travel to Sweden for a festival.",
		"The holiday weekend comes with 2019’s North American box at $5.7 billion through Tuesday, trailing 2018 by 9.3% at the same point — despite the stellar performance of “Avengers: Endgame,” which has topped $843 million after 68 days. Summer box office is also off 6.3% to $2.21 billion, according to Comscore.",
	}

	// query for search
	textF = "Nintendo"
)
