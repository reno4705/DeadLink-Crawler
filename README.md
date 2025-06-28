# DeadLink-Crawler 🕷️

**DeadLink-Crawler** is a fast, concurrent CLI tool built in Go to **recursively crawl a website** and **detect broken (dead) links** — specifically links that return HTTP status codes in the `4xx` or `5xx` range.

---

## 📌 Features

- ✅ Detects **dead** and **working** links on any given website
- 🔄 **Recursively** crawls internal links
- 🔗 Categorizes internal and external links
- 🔁 Uses **concurrency** to speed up crawling
- 🔒 Thread-safe link tracking to prevent duplicate crawls
- 📋 Clean CLI output showing status of each link

---

## 🛠️ Requirements

- Go 1.20 or higher installed  
- Internet connection (for crawling)

---

## 📦 Installation

1. **Clone the repository**

```bash
git clone https://github.com/your-username/deadlink-crawler.git
cd deadlink-crawler
```

2. **Install dependencies**

This project uses a single external package:

```bash
go get golang.org/x/net/html
```

---

## 🚀 Running the Project

Simply run:

```bash
go run main.go
```

> The tool will crawl the default URL hardcoded in `main.go`, which is:
>
> `https://scrape-me.dreamsofcode.io`

You will see a list of links printed as either:
- `Ok LINK: <url> -> <status code>`
- `DEAD LINK: <url> -> <status code>`
- `Error: <url> -> <error message>`

---

## 🧠 How It Works

1. Starts with a base URL.
2. Downloads and parses HTML to extract all `<a href="...">` links.
3. Converts all links to absolute URLs and classifies them as:
   - **Internal**: Same domain as the base
   - **External**: Different domain
4. Internal links are recursively crawled.
5. Each link is fetched with an HTTP GET:
   - `200–399`: Considered OK
   - `400+`: Marked as dead
6. All crawled URLs are tracked to avoid re-checking or infinite loops.
7. Concurrency (via goroutines + waitgroups) is used to parallelize the crawling.

---

## 📁 File Structure

```bash
.
├── main.go         # Main Go source file
├── go.mod          # Go module file
└── README.md       # Project documentation
```

---

## 🧪 Sample Output

```text
Ok LINK: https://scrape-me.dreamsofcode.io -> 200
Ok LINK: https://youtube.com/@dreamsofcode -> 200
Ok LINK: https://scrape-me.dreamsofcode.io/nirvana -> 200
Ok LINK: https://scrape-me.dreamsofcode.io/about -> 200
DEAD LINK: https://scrape-me.dreamsofcode.io/nevermind -> 404
DEAD LINK: https://scrape-me.dreamsofcode.io/in-utero -> 404
Ok LINK: https://scrape-me.dreamsofcode.io/anime?name=bleach -> 200
Ok LINK: https://scrape-me.dreamsofcode.io/anime?name=Jujutsu%20kaizen -> 200
Ok LINK: https://scrape-me.dreamsofcode.io/naruto -> 200
DEAD LINK: https://scrape-me.dreamsofcode.io/teapot -> 418
DEAD LINK: https://scrape-me.dreamsofcode.io/busted -> 401
DEAD LINK: https://scrape-me.dreamsofcode.io/mars -> 404
DEAD LINK: https://scrape-me.dreamsofcode.io/venus -> 404
Error: http://10.255.255.1 -> Get "http://10.255.255.1": dial tcp 10.255.255.1:80: connect: no route to host
```

---

## ✏️ Customizing the Start URL

To crawl a different website, just replace the value of `startURL` in `main()`:

```go
func main() {
	startURL := "https://yourwebsite.com"
	...
}
```

⚠️ Note:
- Works best with static (non-JavaScript rendered) websites.
- Respects the current domain only — external links are **checked**, not crawled.

---

## 📖 Technologies Used

- [Go (Golang)](https://golang.org)
- [`golang.org/x/net/html`](https://pkg.go.dev/golang.org/x/net/html) for HTML parsing
- `net/http`, `url`, `sync` — Go standard libraries

---

## 💡 Ideas for Future Improvements

- Export dead links to a `.txt` or `.csv` file
- Add depth limiting
- Add support for robots.txt
- Rate limiting to avoid overwhelming servers
- CLI flag support (`flag` or `cobra`)

---

## 📄 License

This project is licensed under the MIT License.  
Feel free to use, modify, and distribute!

---

## 🙌 Acknowledgements

- Inspired by real-world site health checks
- Test site: [scrape-me.dreamsofcode.io](https://scrape-me.dreamsofcode.io)

---

## 🤝 Contributing

Pull requests and suggestions are welcome.  
Please open an issue first to discuss what you would like to change.

---