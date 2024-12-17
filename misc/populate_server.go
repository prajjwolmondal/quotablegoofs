package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"quotablegooofs.prajjmon.net/internal/models"
)

var (
	base_url = "http://localhost:8000"
)

func main() {
	populateOneLiners()
	populateMultiLineJokes()
	populateKnockKnockJokes()
	populateQuotes()
}

type JokePostBody struct {
	JokeType models.JokeType `json:"joke_type"`
	Content  []string        `json:"content"`
	Source   string          `json:"source"`
}

func populateOneLiners() {
	oneLineJokes := []string{"I told them I wanted to be a comedian, and they laughed; then I became a comedian, no one's laughing now", "Procrastination: working tomorrow for a better today.", "No matter how good I get at tennis, I will never be better than a wall.", "I love the word frequently, and I try to use it as much as possible.", "You can't trust atoms, they make up everything.", "I shouldn't have driven home last night... Especially since i walked to the bar..", "I want to die peacefully in my sleep like my grandfather, not screaming in terror like his passengers.", " I used to be indecisive but now I'm not sure.", "I hate Russian dolls, they're so full of themselves.", "A recent survey showed that 6 out of 7 dwarfs are not happy.", "The easiest time to add insult to injury is when you're signing somebody's cast.", "Two antennas meet on a roof and fall in love; the wedding wasn't much but the reception was excellent.", "There are only two hard things in Computer Science: cache invalidation, naming things, and off-by-one errors.", "A man walked into his house and was delighted when he discovered that someone had stolen all of his lamps.", "It's hard to explain puns to kleptomaniacs because they always take things literally.", "I asked my North Korean friend how it was there, he said he couldn't complain.", "I discovered a substance that had no mass, and I was like '0mg!'", "I was so surprised when the stationary store moved", "There are 3 kinds of people in this world, those who can count and those who can't", "A Freudian slip is when you say one thing but mean your mother.", "I went shopping for a pair of camouflage pants. But I couldn't find any.", "Bacon and eggs walk into a diner, and the host says - sorry, we don't serve breakfast here.", "I used to steal soap, but I'm clean now.", "Together, I can beat dissociative identity.", "There's no 'I' in Denial.", "I'm against picketing, but I don't know how to show it.", "I often say to myself - I can't believe that cloning machine worked!", "Exaggerations went up by a million percent last year.", "If life hands you melons, you might be dyslexic.", "Came across a mass grave of snowmen, turns out to be a field of carrots.", "The advantage of easy origami is twofold.", "I've always wanted to be a comedian, but that'll never happen because I always punch up the fuck line.", "Velcro - what a rip-off.", "I went to a seafood disco last week...and pulled a mussel.", "I was reading a book titled - The History of Glue - I couldn't put it down.", "As an agnostic, dyslexic insomniac I lie awake all night wondering if there really is a dog.", "It was an emotional wedding. Even the cake was in tiers."}

	for r := range oneLineJokes {
		joke := JokePostBody{
			JokeType: models.OneLiner,
			Content:  []string{oneLineJokes[r]},
			Source:   "Unknown",
		}

		err := convertToJsonAndPostObject(joke, "joke")
		if err != nil {
			fmt.Printf("ERROR! %s", err.Error())
			return
		}
	}

}

// {"lines": ["Why did the chicken go to the seance?", "To talk to the other side."]}
func populateMultiLineJokes() {
	multiLineJokes := [][]string{
		{"Why did the chicken go to the seance?", "To talk to the other side."},
		{"Why don't some couples go to the gym?", "Because some relationships don't work out."},
		{"What did the shark say when he ate the clownfish?", "This tastes a little funny."},
		{"What do you call a woman with one leg?", "Eileen."},
		{"What did the buffalo say when his son left for college?", "Bison."},
		{"What do you call an apology written in dots and dashes?", "Re-Morse code."},
		{"Did you hear about the two people who stole a calendar?", "They each got six months."},
		{"Why do French people eat snails?", "They don't like fast food."},
		{"What did 0 say to 8?", "Nice belt."},
		{"What did the football coach say to the broken vending machine?", "Give me my quarterback."},
		{"What sits at the bottom of the sea and twitches?", "A nervous wreck."},
		{"Why did the Oreo go to the dentist?", "Because he lost his filling."},
		{"What does a house wear?", "Address."},
		{"Why was the broom late for school?", "It over-swept."},
		{"What did one DNA say to the other DNA?", "Do these genes make me look fat?"},
		{"What happens to an illegally parked frog?", "It gets toad away."},
		{"Why aren't dogs good dancers?", "Because they have two left feet."},
		{"What's brown and sticky?", "A stick."},
		{"My twin sister called me from prison.", "She said: 'You know how we finish each other's sentences?'"},
	}

	for i := range multiLineJokes {
		joke := JokePostBody{
			JokeType: models.MultiLine,
			Content:  multiLineJokes[i],
			Source:   "Unknown",
		}

		err := convertToJsonAndPostObject(joke, "joke")
		if err != nil {
			fmt.Printf("ERROR! %s", err.Error())
			return
		}
	}
}

// {"lines": ["Knock, knock.", "Who's there?", "Atch.", "Atch who?", "Bless you!"]}
func populateKnockKnockJokes() {
	knockKnockJokes := [][]string{
		{"Knock, knock", "Lettuce", "Lettuce who?", "Lettuce in, it's freezing out here!"},
		{"Knock, knock", "Atch", "Atch who?", "Bless you!"},
		{"Knock, knock", "Cows go", "Cows go who?", "No silly, cow says moooo!"},
		{"Knock, knock", "A pile-up", "A pile-up who?", "Oh no, yuck!"},
		{"Knock, knock", "Says", "Says who?", "Says me!"},
		{"Knock, knock", "Nobel", "Nobel who?", "Nobel, that's why I knocked!"},
		{"Knock, knock", "Luke", "Luke who?", "Luke through the peephole and find out."},
		{"Knock, knock", "Candice", "Candice who ?", "Candice joke get any worse?"},
	}

	for i := range knockKnockJokes {
		joke := JokePostBody{
			JokeType: models.KnockKnock,
			Content:  knockKnockJokes[i],
			Source:   "Unknown",
		}

		err := convertToJsonAndPostObject(joke, "joke")
		if err != nil {
			fmt.Printf("ERROR! %s", err.Error())
			return
		}
	}
}

func convertToJsonAndPostObject(object any, path string) error {
	jsonObject, err := json.Marshal(object)
	if err != nil {
		return err
	}

	err = sendPostRequest(jsonObject, path)
	if err != nil {
		return err
	}

	return nil
}

type QuotePostBody struct {
	Content []string `json:"content"`
	Source  string   `json:"source"`
}

func populateQuotes() {
	quotes := []QuotePostBody{
		{Content: []string{"Winners never quit, and quitters never win."}, Source: "Vince Lombardi"},
		{Content: []string{"Don't let the fear of losing be greater than the excitement of winning."}, Source: "Robert Kiyosaki"},
		{Content: []string{"Start where you are. Use what you have. Do what you can."}, Source: "Arthur Ashe"},
		{Content: []string{"You must expect great things of yourself before you can do them."}, Source: "Michael Jordan"},
		{Content: []string{"Do what you have to do until you can do what you want to do."}, Source: "Oprah Winfrey"},
		{Content: []string{"You don't have to see the whole staircase, just take the first step."}, Source: "Martin Luther King Jr."},
		{Content: []string{"We are what we repeatedly do. Excellence, then, is not an act, but a habit."}, Source: "Aristotle"},
		{Content: []string{"Change your thoughts and you change your world."}, Source: "Norman Vincent Peale"},
		{Content: []string{"It's hard to beat a person who never gives up."}, Source: "Babe Ruth"},
		{Content: []string{"The only person you should try to be better than, is the person you were yesterday."}, Source: "Matty Mullens"},
		{Content: []string{"The difference between a stumbling block and a stepping stone is how high you raise your foot."}, Source: "Benny Lewis"},
		{Content: []string{"You don't drown by falling in water; you drown by staying there."}, Source: "Robert Collier"},
		{Content: []string{"Better to do something imperfectly than to do nothing flawlessly."}, Source: "Robert Schuller"},
		{Content: []string{"Our greatest weakness lies in giving up. The most certain way to succeed is always to try just one more time."}, Source: "Unknown"},
		{Content: []string{"If you want something you've never had, you must be willing to do something youve never done."}, Source: "Thomas Jefferson"},
		{Content: []string{"Nobody can go back and start a new beginning, but anyone can start today and make a new ending."}, Source: "Maria Robinson"},
		{Content: []string{"The beginning is the most important part of the work."}, Source: "Plato"},
		{Content: []string{"I cannot express how important it is to believe that taking one tiny—and possibly very uncomfortable—step at a time can ultimately add up to a great distance."}, Source: "Tig Notaro"},
		{Content: []string{"Do your thing and don't care if they like it."}, Source: "Tina Fey"},
		{Content: []string{"Try to be a rainbow in someone else's cloud."}, Source: "Maya Angelou"},
		{Content: []string{"Choose to be optimistic, it feels better."}, Source: "Dali Lama"},
		{Content: []string{"Life is like riding a bicycle. To keep your balance, you must keep moving."}, Source: "Albert Einstein"},
		{Content: []string{"It is never too late to be what you might have been."}, Source: "George Eliot"},
		{Content: []string{"Some people look for a beautiful place. Others make a place beautiful."}, Source: "Hazrat Inayat Khan"},
		{Content: []string{"We must be willing to let go of the life we planned so as to have the life that is waiting for us."}, Source: "Joseph Campbell"},
		{Content: []string{"If I cannot do great things, I can do small things in a great way."}, Source: "Martin Luther King, Jr."},
		{Content: []string{"The bad news is time flies. The good news is you're the pilot."}, Source: "Michael Altshuler"},
		{Content: []string{"There are years that ask questions and years that answer."}, Source: "Zora Neale Hurston"},
		{Content: []string{"Each of us is more than the worst thing we've ever done."}, Source: "Bryan Stevenson"},
		{Content: []string{"An ounce of action is worth a ton of theory."}, Source: "Friedrich Engels"},
		{Content: []string{"Don't take life too seriously. You'll never get out of it alive."}, Source: "Elbert Hubbard"},
	}

	for r := range quotes {
		err := convertToJsonAndPostObject(quotes[r], "quote")
		if err != nil {
			fmt.Printf("ERROR! %s", err.Error())
			return
		}
	}
}

func sendPostRequest(jsonBody []byte, path string) error {
	requestUrl := fmt.Sprintf("%s/%s", base_url, path)
	fmt.Printf("[%s] POST >>> %s with body: \n %s \n", time.Now(), requestUrl, string(jsonBody))

	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Printf("Error with POST request to %s | response statuscode: %d, status: %s, body: %s \n", requestUrl, resp.StatusCode, resp.Status, bytes.TrimSpace(body))
		return errors.New("received error status code")
	}

	fmt.Printf("[%s] POST %s <<< %d \n", time.Now(), requestUrl, resp.StatusCode)

	return nil
}
